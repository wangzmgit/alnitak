package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/ffmpeg"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

// ==============================================================================
// 第一部分：接口抽象与核心服务结构 (解决强耦合与可测试性问题)
//
// 常量定义已统一移至 internal/ffmpeg/constant.go
// ==============================================================================

// StorageUploader 抽象存储接口，屏蔽具体 OSS 实现，便于 Mock 测试
type StorageUploader interface {
	PutObjectFromFile(objectKey, filePath string) error
	DeleteObject(objectKey string) error
}

type TranscodingProcess struct {
	VideoID    uint
	ResourceID uint
	PID        int
	Cmd        *exec.Cmd
	CancelFunc context.CancelFunc
	OutputDir  string
}

type resourceTranscodingProgress struct {
	VideoID    uint
	ResourceID uint
	Details    map[string]vo.TranscodingProgressItem

	// 上传阶段进度
	UploadOSS      string  // aliyun/minio/cloudflare/local
	UploadProgress float64 // 0-100
	UploadStatus   string  // "" / uploading / success / fail / local
}

// TranscodeService 转码服务核心实例 (Orchestrator)
type TranscodeService struct {
	db       *gorm.DB
	uploader StorageUploader // 依赖注入存储层

	maxConcur int           // 最大并发数
	semaphore chan struct{} // 全局并发信号量

	gpuMu        sync.RWMutex
	gpuAvailable bool
	gpuFailCount int

	processMu sync.RWMutex
	processes map[uint][]TranscodingProcess

	progressMu sync.RWMutex
	progress   map[uint]*resourceTranscodingProgress
}

var (
	defaultTranscoder *TranscodeService
	initOnce          sync.Once

	currentTranscoderInstance Transcoder
	currentTranscoderMode     string
	currentTranscoderMu       sync.RWMutex
)

// GetCurrentTranscoder 按 config 返回当前转码后端，支持运行时切换。
//   - mode=local（默认）: 返回 LocalTranscoder，封装现有 TranscodeService
//   - mode=remote:        返回 RemoteTranscoder，通过 Redis + OSS 与 Worker 通信
func GetCurrentTranscoder() Transcoder {
	mode := global.Config.Transcoding.Mode
	if mode == "" {
		mode = "local"
	}

	currentTranscoderMu.RLock()
	if currentTranscoderInstance != nil && currentTranscoderMode == mode {
		currentTranscoderMu.RUnlock()
		return currentTranscoderInstance
	}
	currentTranscoderMu.RUnlock()

	currentTranscoderMu.Lock()
	defer currentTranscoderMu.Unlock()

	// Double check after acquiring write lock
	if currentTranscoderInstance != nil && currentTranscoderMode == mode {
		return currentTranscoderInstance
	}

	switch mode {
	case "remote":
		utils.InfoLog("【转码模式】remote（远程 Worker 池）", "transcoding")
		currentTranscoderInstance = NewRemoteTranscoder()
	default:
		utils.InfoLog("【转码模式】local（本地进程内转码）", "transcoding")
		currentTranscoderInstance = NewLocalTranscoder()
	}
	currentTranscoderMode = mode
	return currentTranscoderInstance
}

// GetTranscoder 获取单例模式的 TranscodeService
func GetTranscoder() *TranscodeService {
	initOnce.Do(func() {
		// 动态确定并发数，使用全局配置和 transcoding_constants.go 中的常量兜底
		maxConcur := global.Config.Transcoding.MaxCpuConcurrency
		if maxConcur <= 0 {
			maxConcur = maxCPUConcurrentTranscoding
		}
		if global.Config.Transcoding.UseGpu {
			maxConcur = global.Config.Transcoding.MaxGpuConcurrency
			if maxConcur <= 0 {
				maxConcur = maxGPUConcurrentTranscoding
			}
		}

		defaultTranscoder = &TranscodeService{
			db:           global.Mysql,
			uploader:     global.Storage, // 注入真实的存储引擎
			maxConcur:    maxConcur,
			semaphore:    make(chan struct{}, maxConcur),
			gpuAvailable: true,
			processes:    make(map[uint][]TranscodingProcess),
			progress:     make(map[uint]*resourceTranscodingProgress),
		}
		utils.InfoLog(fmt.Sprintf("【转码并发控制初始化】最大并发数=%d (GPU=%v)", maxConcur, global.Config.Transcoding.UseGpu), "transcoding")
	})
	return defaultTranscoder
}

// ==============================================================================
// 第三部分：对外暴露的包级函数 (保持原有 API 兼容)
// ==============================================================================

// VideoTransCoding 同步启动转码（阻塞至完成）。
//   - mode=local:  内部调用 ProcessVideo，行为与重构前一致
//   - mode=remote: 直接调用 GetTranscoder().ProcessVideo（不经过 RemoteTranscoder）
//
// 调用方需要异步时自行加 go 关键字，如 upload.go/resource.go 已改为
// GetCurrentTranscoder().Enqueue()。
func VideoTransCoding(transcodingInfo *dto.TranscodingInfo) {
	GetTranscoder().ProcessVideo(context.Background(), transcodingInfo)
}

// GetCurrentTranscoder 见上方定义。
func ResetGPUState() { GetTranscoder().ResetGPUState() }

func GetVideoTranscodingProgress(videoID uint) (float64, []vo.TranscodingProgressItem, *vo.UploadProgressInfo) {
	// 远程模式：从 Redis 读取转码进度
	if global.Config.Transcoding.Mode == "remote" {
		return getRemoteVideoTranscodingProgress(videoID)
	}
	// 本地模式：从内存读取
	return GetTranscoder().GetVideoTranscodingProgress(videoID)
}

// getRemoteVideoTranscodingProgress 从 Redis 读取远程转码进度。
func getRemoteVideoTranscodingProgress(videoID uint) (float64, []vo.TranscodingProgressItem, *vo.UploadProgressInfo) {
	var resources []model.Resource
	global.Mysql.Model(&model.Resource{}).Where("vid = ?", videoID).Find(&resources)

	if len(resources) == 0 {
		return 0, nil, nil
	}

	details := make([]vo.TranscodingProgressItem, 0)
	var totalProgress float64
	var totalCount int
	var uploadInfo *vo.UploadProgressInfo
	resourceUploads := make(map[uint]*vo.UploadProgressInfo)

	rdb := global.Redis.RawClient()
	ctx := context.Background()

	// 记录有 Redis 数据的 resource，用于后续补齐排队中的分P
	foundInRedis := make(map[uint]bool, len(resources))
	// 收集已有的画质名称，用于排队中分P的占位
	knownQualities := make([]string, 0)

	for _, res := range resources {
		statusKey := fmt.Sprintf("transcoding:status:%d", res.ID)
		statusMap, err := rdb.HGetAll(ctx, statusKey).Result()
		if err != nil || len(statusMap) == 0 {
			continue
		}
		foundInRedis[res.ID] = true

		status := statusMap["status"]

		// 解析总体进度
		var overall float64
		if p, ok := statusMap["progress"]; ok {
			fmt.Sscanf(p, "%f", &overall)
		}

		// 解析每个画质的进度（progress_{qualityName} 字段）
		var resProgressSum float64
		var resProgressCount int
		for key, val := range statusMap {
			if !strings.HasPrefix(key, "progress_") {
				continue
			}
			qualityName := strings.TrimPrefix(key, "progress_")
			// 收集已知画质名
			seen := false
			for _, q := range knownQualities {
				if q == qualityName {
					seen = true
					break
				}
			}
			if !seen {
				knownQualities = append(knownQualities, qualityName)
			}

			// 读取画质独立状态（success/failed），没有则沿用整体状态
			qStatus := status
			if s, ok := statusMap["status_"+qualityName]; ok {
				qStatus = s
			}
			var pct float64
			fmt.Sscanf(val, "%f", &pct)
			details = append(details, vo.TranscodingProgressItem{
				ResourceID: res.ID,
				Quality:    qualityName,
				Progress:   pct,
				Status:     qStatus,
			})
			resProgressSum += pct
			resProgressCount++
		}

		// 整体进度计算：
		//   编码中 → 画质进度平均值
		//   上传/完成 → 画质已全部 100%，用 progress 字段（Worker 设为 90/100）
		if resProgressCount > 0 {
			if strings.Contains(status, "upload") || status == "completed" || status == "encoding_done" {
				totalProgress += overall * float64(resProgressCount)
			} else {
				totalProgress += resProgressSum
			}
			totalCount += resProgressCount
		}

		// 解析上传进度（按分P存储）
		if uploadStr, ok := statusMap["upload"]; ok && uploadStr != "" {
			var up struct {
				OssType  string  `json:"ossType"`
				Progress float64 `json:"progress"`
				Status   string  `json:"status"`
			}
			if err := json.Unmarshal([]byte(uploadStr), &up); err == nil {
				info := &vo.UploadProgressInfo{
					OssType:  up.OssType,
					Progress: up.Progress,
					Status:   up.Status,
				}
				resourceUploads[res.ID] = info
				uploadInfo = info
			}
		}
	}

	// 补齐没有 Redis 数据的 resource（排队中 / 已完成清理）
	// 让前端能看见全部分P的状态，而不是只显示正在处理的那几个
	for _, res := range resources {
		if foundInRedis[res.ID] {
			continue
		}

		var displayStatus string
		switch res.Status {
		case global.CREATED_VIDEO, global.VIDEO_PROCESSING:
			displayStatus = "waiting"
		case global.WAITING_REVIEW, global.AUDIT_APPROVED:
			displayStatus = "success"
		case global.PROCESSING_FAIL:
			displayStatus = "fail"
		default:
			continue
		}

		// 如果已有其他分P的已知画质，为每个画质生成一条排队/完成占位
		if len(knownQualities) > 0 {
			for _, q := range knownQualities {
				var pct float64
				if displayStatus == "success" {
					pct = 100
				}
				details = append(details, vo.TranscodingProgressItem{
					ResourceID: res.ID,
					Quality:    q,
					Progress:   pct,
					Status:     displayStatus,
				})
				totalProgress += pct
				totalCount++
			}
		} else {
			// 没有任何资源有 Redis 数据（首次启动，全排队中），加一条无画质占位
			details = append(details, vo.TranscodingProgressItem{
				ResourceID: res.ID,
				ResourceTitle: res.Title,
				Status:        displayStatus,
			})
		}
	}

	if totalCount == 0 {
		return 0, details, uploadInfo
	}

	sort.Slice(details, func(i, j int) bool {
		// 无画质的放最后
		if details[i].Quality == "" {
			return false
		}
		if details[j].Quality == "" {
			return true
		}
		iw, ih, ifps, _ := parseProgressQualitySortKey(details[i].Quality)
		jw, jh, jfps, _ := parseProgressQualitySortKey(details[j].Quality)
		if ih != jh {
			return ih > jh
		}
		if iw != jw {
			return iw > jw
		}
		if ifps != jfps {
			return ifps > jfps
		}
		return iw < jw
	})

	// 将上传进度挂到对应分P的第一个画质条目上
	assignedUpload := make(map[uint]bool)
	for i := range details {
		if !assignedUpload[details[i].ResourceID] {
			if up, ok := resourceUploads[details[i].ResourceID]; ok {
				details[i].Upload = up
				assignedUpload[details[i].ResourceID] = true
			}
		}
	}

	return totalProgress / float64(totalCount), details, uploadInfo
}

func StopTranscodingAndCleanup(videoID uint) error {
	return GetTranscoder().StopTranscodingAndCleanup(videoID)
}

func HasTranscodingProcess(videoID uint) bool {
	return GetTranscoder().HasTranscodingProcess(videoID)
}

func GetTranscodingProcessCount(videoID uint) int {
	return GetTranscoder().GetTranscodingProcessCount(videoID)
}

// ==============================================================================
// 第四部分：音频并发同步器 (SharedTask)
// ==============================================================================

type SharedTask struct {
	once sync.Once
	done chan struct{}
	err  error
}

func NewSharedTask() *SharedTask {
	return &SharedTask{done: make(chan struct{})}
}

func (s *SharedTask) DoOrWait(ctx context.Context, taskFn func() error) error {
	s.once.Do(func() {
		defer close(s.done)
		defer func() {
			if r := recover(); r != nil {
				s.err = fmt.Errorf("task panic: %v", r)
			}
		}()
		s.err = taskFn()
	})

	select {
	case <-s.done:
		return s.err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// ==============================================================================
// 第五部分：核心管线编排 (ProcessVideo 流水线)
// ==============================================================================

func (s *TranscodeService) ProcessVideo(ctx context.Context, info *dto.TranscodingInfo) {
	// 前端元数据可能缺少 VideoBitRate/AudioStreams/FPS 等字段，
	// 当关键字段为空且本地源文件可用时，做一次完整 ffprobe 补全。
	needsProbe := info.InputFile != "" && (info.VideoBitRate == 0 || len(info.AudioStreams) == 0 || info.FPS == "")
	if needsProbe {
		if _, err := os.Stat(info.InputFile); err == nil {
			if probed, err := ProcessVideoInfo(info.InputFile); err == nil {
				if info.VideoBitRate == 0 {
					info.VideoBitRate = probed.VideoBitRate
				}
				if len(info.AudioStreams) == 0 {
					info.AudioStreams = probed.AudioStreams
				}
				if info.FPS == "" {
					info.FPS = probed.FPS
					info.FPS30 = probed.FPS30
					info.FPS60 = probed.FPS60
				}
				if info.Width == 0 || info.Height == 0 {
					info.Width = probed.Width
					info.Height = probed.Height
				}
				if info.Duration == 0 {
					info.Duration = probed.Duration
				}
				if info.AudioBitRate == 0 {
					info.AudioBitRate = probed.AudioBitRate
					info.AudioSampleRate = probed.AudioSampleRate
					info.AudioChannels = probed.AudioChannels
				}
				utils.InfoLog(fmt.Sprintf("【自动补全视频信息】VideoID=%d, BitRate=%d, 音轨数=%d, FPS=%s",
					info.VideoID, info.VideoBitRate, len(info.AudioStreams), info.FPS), "transcoding")
			}
		}
	}
	if info.FPS == "" {
		info.FPS = "30/1"
		info.FPS30 = "30/1"
		info.FPS60 = "60000/1001"
		utils.InfoLog(fmt.Sprintf("【FPS探测失败，使用默认值】VideoID=%d, FPS=30/1", info.VideoID), "transcoding")
	}

	targets := ffmpeg.GetTranscodingTargets(info.Width, info.Height, info.VideoBitRate, info.FPS30, info.FPS60, global.Config.Transcoding.Generate1080p60)
	s.initResourceProgress(info.VideoID, info.ResourceID, targets)
	defer s.clearProgress(info.ResourceID)

	utils.InfoLog(fmt.Sprintf("【转码开始】VideoID=%d, ResourceID=%d, 目标数量=%d, 音轨数=%d", info.VideoID, info.ResourceID, len(targets), len(info.AudioStreams)), "transcoding")

	// 【多音轨】每个音轨对应一个 SharedTask（确保每条音轨只编一次，跨 quality goroutine 共享）
	audioTrackTasks := make([]*SharedTask, len(info.AudioStreams))
	for i := range audioTrackTasks {
		audioTrackTasks[i] = NewSharedTask()
	}
	// 无多音轨探测时，回退单音轨模式
	singleAudioTask := NewSharedTask()

	var wg sync.WaitGroup
	var successCount atomic.Int32

	// 1. 并发处理所有画质分片
	for _, target := range targets {
		wg.Add(1)
		go func(t ffmpeg.TranscodingTarget) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					qualityName := t.Resolution + "_" + t.BitrateRate + "_" + t.FpsName
					s.markQualityFailed(info.ResourceID, qualityName)
					utils.ErrorLog(fmt.Sprintf("【Goroutine panic】%s: %v", qualityName, r), "transcoding", "")
				}
			}()

			// 如果 context 已取消，goroutine 直接返回
			if ctx.Err() != nil {
				return
			}

			if err := s.processSingleQuality(ctx, info, t, audioTrackTasks, singleAudioTask); err != nil {
				utils.ErrorLog(fmt.Sprintf("【转码处理失败】%s", err.Error()), "transcoding", "")
				return
			}
			successCount.Add(1)
		}(target)
	}

	// 等待所有分片完成或 context 取消
	waitDone := make(chan struct{}, 1)
	go func() {
		wg.Wait()
		waitDone <- struct{}{}
	}()
	select {
	case <-waitDone:
	case <-ctx.Done():
		utils.InfoLog(fmt.Sprintf("【转码取消】VideoID=%d, ResourceID=%d", info.VideoID, info.ResourceID), "transcoding")
		// 等待 goroutine 自然结束（FFmpeg 被 exec.CommandContext 杀死后很快返回）
		<-waitDone
	}
	utils.InfoLog(fmt.Sprintf("【所有分段转码完成】成功=%d, 总数=%d", successCount.Load(), len(targets)), "transcoding")

	// 1.5 多音轨落库（在所有 quality 完成后执行一次）
	if len(info.AudioStreams) > 0 {
		_ = s.saveAudioTrackRecords(ctx, info)
	}

	// 2. 失败拦截与处理
	if successCount.Load() == 0 {
		utils.InfoLog(fmt.Sprintf("【所有转码失败】VideoID=%d, ResourceID=%d, 标记为失败", info.VideoID, info.ResourceID), "transcoding")
		_ = s.completeTransaction(ctx, info, global.PROCESSING_FAIL)
		return
	}

	// 3. 上传 OSS (安全重试)
	if global.Config.Storage.OssType != "local" {
		if err := s.uploadToOSS(ctx, info); err != nil {
			utils.ErrorLog("【OSS上传流程失败】", "transcoding", err.Error())
			s.markUploadFailed(ctx, info)
			return
		}
	} else {
		utils.InfoLog("【跳过OSS上传】使用本地存储", "transcoding")
	}

	// 4. 收尾落库
	utils.InfoLog(fmt.Sprintf("【准备最终收尾】VideoID=%d", info.VideoID), "transcoding")
	if err := s.completeTransaction(ctx, info, global.WAITING_REVIEW); err != nil {
		utils.ErrorLog("【业务收尾落库失败】", "transcoding", err.Error())
	}

	// 5. 异步上传到备用 OSS（多源容灾，不影响主流程）
	if global.StorageBackup != nil {
		utils.InfoLog(fmt.Sprintf("【备用OSS上传】触发异步上传 VideoID=%d, DirName=%s", info.VideoID, info.DirName), "transcoding")
		go s.uploadToBackup(info)
	}
}

func (s *TranscodeService) processSingleQuality(
	ctx context.Context, info *dto.TranscodingInfo, t ffmpeg.TranscodingTarget,
	audioTrackTasks []*SharedTask, singleAudioTask *SharedTask,
) error {
	qualityName := t.Resolution + "_" + t.BitrateRate + "_" + t.FpsName
	videoFile := info.OutputDir + qualityName + "_video.m4s"

	s.semaphore <- struct{}{}
	defer func() { <-s.semaphore }()

	s.updateProgress(info.ResourceID, qualityName, 0, "processing")
	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	if subCtx.Err() != nil {
		s.markQualityFailed(info.ResourceID, qualityName)
		return subCtx.Err()
	}

	// 视频编码
	if err := s.encodeVideoWithFallback(subCtx, info, t, videoFile, qualityName, cancel); err != nil {
		s.markQualityFailed(info.ResourceID, qualityName)
		return fmt.Errorf("【视频编码失败】%s: %w", qualityName, err)
	}

	if subCtx.Err() != nil {
		return fmt.Errorf("【被取消】%s", qualityName)
	}

	// ===== 音频编码（多音轨支持） =====
	if len(info.AudioStreams) > 0 {
		// 多音轨模式：每个音轨独立编码，默认音轨失败时整体失败
		for i, stream := range info.AudioStreams {
			audioFileName := ffmpeg.AudioFileNameForTrack(stream.Language)
			audioFilePath := info.OutputDir + audioFileName

			// 当前帧率用于计算 B 帧偏移（使用首个 quality 的 fps 即可，所有 quality 共享）
			leadMs := ffmpeg.BFramePresentationLeadMs(t.FPS)
			err := audioTrackTasks[i].DoOrWait(subCtx, func() error {
				utils.InfoLog(fmt.Sprintf("【编码音频】%s 语言=%s 码率=%dk", audioFileName, stream.Language, stream.BitRate/1000), "transcoding")
				return encodeAudioTrack(subCtx, info.InputFile, audioFilePath,
					stream.StreamIndex, stream.BitRate, stream.SampleRate, stream.Channels,
					info.Duration, leadMs)
			})
			if err != nil {
				if i == 0 {
					// 默认音轨失败 → 该 quality 失败
					s.markQualityFailed(info.ResourceID, qualityName)
					return fmt.Errorf("【默认音轨编码失败】%s: %w", qualityName, err)
				}
				// 非默认音轨失败 → 仅记录日志
				utils.ErrorLog(fmt.Sprintf("【附加音轨编码失败】语言=%s: %s", stream.Language, err.Error()), "transcoding", "")
			}
		}
	} else {
		// 单音轨回退模式
		audioFile := info.OutputDir + "audio.m4s"
		err := singleAudioTask.DoOrWait(subCtx, func() error {
			utils.InfoLog(fmt.Sprintf("【编码音频】audio.m4s (码率=%dk)", info.AudioBitRate/1000), "transcoding")
			leadMs := ffmpeg.BFramePresentationLeadMs(t.FPS)
			return encodeAudioOnly(subCtx, info.InputFile, audioFile, info.AudioBitRate, info.AudioSampleRate, info.AudioChannels, info.Duration, leadMs)
		})
		if err != nil {
			s.markQualityFailed(info.ResourceID, qualityName)
			return fmt.Errorf("【音频编码失败】%s: %w", qualityName, err)
		}
	}

	// 解析存储
	audioFile := info.OutputDir + "audio.m4s"
	if len(info.AudioStreams) > 0 {
		audioFile = info.OutputDir + ffmpeg.AudioFileNameForTrack(info.AudioStreams[0].Language)
	}
	if err := s.saveIndexRecord(subCtx, info, qualityName, videoFile, audioFile); err != nil {
		s.markQualityFailed(info.ResourceID, qualityName)
		return fmt.Errorf("【索引生成失败】%s: %w", qualityName, err)
	}

	s.updateProgress(info.ResourceID, qualityName, 100, "success")
	return nil
}

// encodeVideoWithFallback 动态路由 GPU/CPU 策略
func (s *TranscodeService) encodeVideoWithFallback(
	ctx context.Context, info *dto.TranscodingInfo, t ffmpeg.TranscodingTarget,
	videoFile, qualityName string, cancel context.CancelFunc,
) error {
	useGpu := global.Config.Transcoding.UseGpu && s.isGPUAvailable()
	useHevc := global.Config.Transcoding.UseH265
	useAv1 := global.Config.Transcoding.UseAv1
	utils.InfoLog(fmt.Sprintf("【开始视频编码】%s 使用GPU=%v AV1=%v H.265=%v", qualityName, useGpu, useAv1, useHevc), "transcoding")

	err := s.runVideoEncodeTask(ctx, info.VideoID, info.ResourceID, info.InputFile, videoFile, t.Resolution, t.BitrateRate, t.FPS, qualityName, info.Duration, useGpu, useAv1, useHevc, cancel)

	if useGpu && err != nil && (strings.Contains(err.Error(), "GPU error") || strings.Contains(err.Error(), "nvenc")) {
		s.handleGPUFailure()
		if ctx.Err() != nil {
			return ctx.Err()
		}
		// GPU AV1 → GPU H.265 → GPU H.264 → CPU H.264 逐级降级
		if useAv1 {
			utils.InfoLog(fmt.Sprintf("【GPU降级】%s AV1→H.265 GPU", qualityName), "transcoding")
			err = s.runVideoEncodeTask(ctx, info.VideoID, info.ResourceID, info.InputFile, videoFile, t.Resolution, t.BitrateRate, t.FPS, qualityName, info.Duration, true, false, true, cancel)
			if err == nil {
				s.recordGPUSuccess()
				return nil
			}
			s.handleGPUFailure()
			if ctx.Err() != nil {
				return ctx.Err()
			}
		}
		if useAv1 || useHevc {
			utils.InfoLog(fmt.Sprintf("【GPU降级】%s H.265→H.264 GPU", qualityName), "transcoding")
			err = s.runVideoEncodeTask(ctx, info.VideoID, info.ResourceID, info.InputFile, videoFile, t.Resolution, t.BitrateRate, t.FPS, qualityName, info.Duration, true, false, false, cancel)
			if err == nil {
				s.recordGPUSuccess()
				return nil
			}
			s.handleGPUFailure()
			if ctx.Err() != nil {
				return ctx.Err()
			}
		}
		utils.InfoLog(fmt.Sprintf("【GPU降级CPU编码】%s", qualityName), "transcoding")
		return s.runVideoEncodeTask(ctx, info.VideoID, info.ResourceID, info.InputFile, videoFile, t.Resolution, t.BitrateRate, t.FPS, qualityName, info.Duration, false, false, false, cancel)
	} else if useGpu && err == nil {
		s.recordGPUSuccess()
	}
	return err
}

func (s *TranscodeService) saveIndexRecord(ctx context.Context, info *dto.TranscodingInfo, qualityName, videoFilePath, audioFilePath string) error {
	vInit, vIndex, err := getMP4InitRange(videoFilePath)
	if err != nil {
		return err
	}
	aInit, aIndex, err := getMP4InitRange(audioFilePath)
	if err != nil {
		return err
	}

	// 从已编码文件探测实际 codec（正确处理 GPU 降级场景）
	videoCodec := ffmpeg.ProbeVideoActualCodec(videoFilePath)

	width, height, bandwidth, frameRateStr := ffmpeg.ParseQualityInfo(qualityName)

	indexFile := &model.VideoIndexFile{
		ResourceID:      info.ResourceID,
		Quality:         qualityName,
		DirName:         info.DirName,
		TotalDuration:   info.Duration,
		VideoFile:       qualityName + "_video.m4s",
		VideoBandwidth:  bandwidth,
		VideoCodec:      videoCodec,
		Width:           width,
		Height:          height,
		FrameRate:       ffmpeg.ParseFPS(frameRateStr),
		VideoInitRange:  vInit,
		VideoIndexRange: vIndex,
		AudioFile:       "audio.m4s",
		AudioBandwidth:  info.AudioBitRate,
		AudioCodec:      ffmpeg.DefaultAudioCodec,
		AudioSampleRate: info.AudioSampleRate,
		AudioInitRange:  aInit,
		AudioIndexRange: aIndex,
	}

	return s.db.WithContext(ctx).Create(indexFile).Error
}

func (s *TranscodeService) completeTransaction(ctx context.Context, info *dto.TranscodingInfo, targetStatus int) error {
	defer s.clearProgress(info.ResourceID)

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var fileCount int64
		tx.Model(&model.VideoIndexFile{}).Where("resource_id = ?", info.ResourceID).Count(&fileCount)
		if fileCount == 0 {
			targetStatus = global.PROCESSING_FAIL
		}

		var currentVideo model.Video
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", info.VideoID).First(&currentVideo).Error; err != nil {
			return err
		}

		var currentResource model.Resource
		if err := tx.Where("id = ?", info.ResourceID).First(&currentResource).Error; err != nil {
			return err
		}

		// 所有新分P统一进审核队列（无论是否替换），审核通过后才对外公开
		isReplace := currentResource.ReplaceID > 0

		// 转码成功后不设可见，等审核通过后统一由 ReviewVideoApproved 设 VisibleStatus=1
		if targetStatus != global.PROCESSING_FAIL && targetStatus != global.UPLOAD_FAILED && isReplace {
			utils.InfoLog(fmt.Sprintf("【资源替换】新ResourceID=%d 转码完成，等待审核，旧ResourceID=%d 仍可见",
				info.ResourceID, currentResource.ReplaceID), "transcoding")
		}

		result := tx.Model(&model.Resource{}).
			Where("id = ? and status = ?", info.ResourceID, global.VIDEO_PROCESSING).
			Update("status", targetStatus)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			// 资源状态已被其他流程变更（例如之前的 completeTransaction 已完成），
			// 跳过后续所有 Video 级别状态更新，避免重复写入错误的状态。
			utils.InfoLog(fmt.Sprintf("资源状态已变更，跳过: ResourceID=%d", info.ResourceID), "transcoding")
			return nil
		}

		if targetStatus != global.PROCESSING_FAIL {
			updateVideoFileStatus(tx, currentResource, info.ResourceID)
		}

		var totalCount, failedCount int64
		tx.Model(&model.Resource{}).Where("vid = ?", info.VideoID).Count(&totalCount)
		tx.Model(&model.Resource{}).Where("vid = ? and status = ?", info.VideoID, global.PROCESSING_FAIL).Count(&failedCount)

		var successCount int64
		tx.Model(&model.Resource{}).Where("vid = ? and status in ?", info.VideoID, []int{global.WAITING_REVIEW, global.AUDIT_APPROVED}).Count(&successCount)

		// 成功路径：等全部分P都完成再更新视频状态，避免单P完成就提前暴露稿件
		if targetStatus != global.PROCESSING_FAIL && successCount < totalCount {
			utils.InfoLog(fmt.Sprintf("多分P等待: VideoID=%d, 已处理 %d/%d 个分P, 暂不更新视频状态",
				info.VideoID, successCount, totalCount), "transcoding")
			return nil
		}

		finalVideoStatus := global.WAITING_REVIEW
		if targetStatus == global.PROCESSING_FAIL {
			// 当前分P失败：如果已有其他分P成功，不标记视频为失败（可单独重试失败的分P）
			if successCount > 0 {
				return nil
			}
			if currentVideo.Status == global.AUDIT_APPROVED || info.OriginalVideoStatus == global.AUDIT_APPROVED {
				finalVideoStatus = global.AUDIT_APPROVED
			} else if failedCount == totalCount {
				// 全部分P都失败时才标记视频失败
				finalVideoStatus = global.PROCESSING_FAIL
			} else {
				// 还有分P在处理中，等待
				return nil
			}
		} else if failedCount == totalCount {
			// 当前成功但其他全失败了——不可能，但兜底
			finalVideoStatus = global.PROCESSING_FAIL
		} else if currentVideo.Status == global.AUDIT_APPROVED || info.OriginalVideoStatus == global.AUDIT_APPROVED {
			finalVideoStatus = global.AUDIT_APPROVED
		}

		return tx.Model(&model.Video{}).Where("id = ?", info.VideoID).Update("status", finalVideoStatus).Error
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	cache.DelVideoInfo(info.VideoID)
	utils.InfoLog(fmt.Sprintf("【事务流转成功】已清理VideoID=%d缓存", info.VideoID), "transcoding")
	return nil
}

// markUploadFailed 在转码成功但上传 OSS 失败时调用。
// 与 completeTransaction(PROCESSING_FAIL) 不同，它只标记资源状态为 UPLOAD_FAILED，
// 不改变 video 状态，并保留转码产物在磁盘上以便后续重试上传。
func (s *TranscodeService) markUploadFailed(ctx context.Context, info *dto.TranscodingInfo) {
	defer s.clearProgress(info.ResourceID)
	if err := s.db.WithContext(ctx).
		Model(&model.Resource{}).
		Where("id = ? and status = ?", info.ResourceID, global.VIDEO_PROCESSING).
		Update("status", global.UPLOAD_FAILED).Error; err != nil {
		utils.ErrorLog("【标记上传失败】", "transcoding", err.Error())
	}
	// 不更新 video 表状态，保持 VIDEO_PROCESSING 以便前端在「处理中」列表可见
}

type uploadTaskResult struct {
	fileName  string
	objectKey string
	success   bool
}

func (s *TranscodeService) uploadToOSS(ctx context.Context, info *dto.TranscodingInfo) error {
	files, err := os.ReadDir(info.OutputDir)
	if err != nil {
		return err
	}

	// 过滤出需要上传的文件
	var uploadEntries []os.DirEntry
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name() == "upload"+info.Suffix && !global.Config.Storage.UploadMp4File {
			continue
		}
		uploadEntries = append(uploadEntries, file)
	}

	totalFiles := len(uploadEntries)
	if totalFiles == 0 {
		s.setUploadProgress(info.ResourceID, 100, "local", "local")
		return nil
	}

	ossType := global.Config.Storage.OssType
	if ossType == "" {
		ossType = "local"
	}
	s.setUploadProgress(info.ResourceID, 0, "uploading", ossType)

	tasks := make(chan os.DirEntry, totalFiles)
	results := make(chan uploadTaskResult, totalFiles)

	var completed atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < ossUploadMaxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range tasks {
				fileName := file.Name()
				objectKey := "video/" + info.DirName + "/" + fileName
				filePath := info.OutputDir + fileName

				// 每次上传新文件前检查 context 是否已取消
				if ctx.Err() != nil {
					results <- uploadTaskResult{fileName: fileName, objectKey: objectKey, success: false}
					return
				}

				var lastErr error
				for attempt := 0; attempt <= ossUploadMaxRetries; attempt++ {
					if attempt > 0 {
						select {
						case <-ctx.Done():
							results <- uploadTaskResult{fileName: fileName, objectKey: objectKey, success: false}
							return
						case <-time.After(ossUploadBackoff[attempt-1]):
						}
					}
					lastErr = s.uploader.PutObjectFromFile(objectKey, filePath)
					if lastErr == nil {
						break
					}
				}

				if lastErr != nil {
					utils.ErrorLog(fmt.Sprintf("【OSS上传失败】%s (重试%d次)", fileName, ossUploadMaxRetries), "oss", lastErr.Error())
					results <- uploadTaskResult{fileName: fileName, objectKey: objectKey, success: false}
				} else {
					done := completed.Add(1)
					s.setUploadProgress(info.ResourceID, float64(done)/float64(totalFiles)*100, "uploading", "")
					results <- uploadTaskResult{fileName: fileName, objectKey: objectKey, success: true}
				}
			}
		}()
	}

	for _, file := range uploadEntries {
		tasks <- file
	}
	close(tasks)
	wg.Wait()
	close(results)

	var uploaded []string
	var failedCount int
	for r := range results {
		if r.success {
			uploaded = append(uploaded, r.objectKey)
		} else {
			failedCount++
		}
	}

	if failedCount > 0 {
		s.setUploadProgress(info.ResourceID, 0, "fail", "")
		// 回滚：删除所有已成功上传的文件，防止 OSS 残留孤儿文件
		for _, key := range uploaded {
			if key == "" {
				continue
			}
			if err := s.uploader.DeleteObject(key); err != nil {
				utils.ErrorLog(fmt.Sprintf("【OSS回滚删除失败】%s", key), "oss", err.Error())
			}
		}
		return fmt.Errorf("OSS上传部分失败: %d个文件失败, 已回滚%d个已上传文件", failedCount, len(uploaded))
	}

	s.setUploadProgress(info.ResourceID, 100, "success", "")
	return nil
}

// uploadToBackup 异步上传转码产物到备用 OSS（多源容灾）。
// 不影响主流程，失败仅记录日志。
func (s *TranscodeService) uploadToBackup(info *dto.TranscodingInfo) {
	backup := global.StorageBackup
	if backup == nil {
		return
	}

	files, err := os.ReadDir(info.OutputDir)
	if err != nil {
		utils.ErrorLog("【备用OSS上传】读取目录失败", "transcoding", err.Error())
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name() == "upload"+info.Suffix && !global.Config.Storage.UploadMp4File {
			continue
		}

		fileName := file.Name()
		objectKey := "video/" + info.DirName + "/" + fileName
		filePath := info.OutputDir + fileName

		var lastErr error
		for attempt := 0; attempt <= ossUploadMaxRetries; attempt++ {
			if attempt > 0 {
				time.Sleep(ossUploadBackoff[attempt-1])
			}
			lastErr = backup.PutObjectFromFile(objectKey, filePath)
			if lastErr == nil {
				break
			}
		}

		if lastErr != nil {
			utils.ErrorLog(fmt.Sprintf("【备用OSS上传失败】%s (重试%d次)", fileName, ossUploadMaxRetries), "oss", lastErr.Error())
			RecordBackupUploadFailure(objectKey, filePath, "video", lastErr.Error())
		} else {
			clearBackupFailure(objectKey)
			utils.InfoLog(fmt.Sprintf("【备用OSS上传成功】%s", fileName), "transcoding")
		}
	}
}

// ==============================================================================
// 第六部分：内部状态管理与追踪
// ==============================================================================

func (s *TranscodeService) isGPUAvailable() bool {
	s.gpuMu.RLock()
	defer s.gpuMu.RUnlock()
	return s.gpuAvailable
}

func (s *TranscodeService) handleGPUFailure() {
	s.gpuMu.Lock()
	defer s.gpuMu.Unlock()
	s.gpuFailCount++
	utils.InfoLog(fmt.Sprintf("【GPU失败】失败次数=%d/%d", s.gpuFailCount, maxGpuFailCountThreshold), "transcoding")
	if s.gpuFailCount >= maxGpuFailCountThreshold {
		s.gpuAvailable = false
		utils.InfoLog("【GPU禁用】达到阈值，已切换到CPU模式", "transcoding")
	}
}

func (s *TranscodeService) recordGPUSuccess() {
	s.gpuMu.Lock()
	defer s.gpuMu.Unlock()
	if s.gpuFailCount > 0 {
		s.gpuFailCount = 0
		utils.InfoLog("【GPU恢复】转码成功，重置失败计数", "transcoding")
	}
}

func (s *TranscodeService) ResetGPUState() {
	s.gpuMu.Lock()
	defer s.gpuMu.Unlock()
	if s.gpuFailCount > 0 || !s.gpuAvailable {
		s.gpuFailCount = 0
		s.gpuAvailable = true
		utils.InfoLog("【GPU重置】状态恢复", "transcoding")
	}
}

func (s *TranscodeService) initResourceProgress(videoID, resourceID uint, targets []ffmpeg.TranscodingTarget) {
	s.progressMu.Lock()
	defer s.progressMu.Unlock()
	state := &resourceTranscodingProgress{
		VideoID:    videoID,
		ResourceID: resourceID,
		Details:    make(map[string]vo.TranscodingProgressItem, len(targets)),
	}
	for _, target := range targets {
		quality := target.Resolution + "_" + target.BitrateRate + "_" + target.FpsName
		state.Details[quality] = vo.TranscodingProgressItem{
			ResourceID: resourceID,
			Quality:    quality,
			Progress:   0,
			Status:     "waiting",
		}
	}
	s.progress[resourceID] = state
}

func (s *TranscodeService) updateProgress(resourceID uint, quality string, progress float64, status string) {
	s.progressMu.Lock()
	defer s.progressMu.Unlock()
	state, ok := s.progress[resourceID]
	if !ok {
		return
	}
	item, exists := state.Details[quality]
	if !exists {
		item = vo.TranscodingProgressItem{ResourceID: resourceID, Quality: quality, Status: "processing"}
	}
	if progress < 0 {
		progress = 0
	} else if progress > 100 {
		progress = 100
	}
	item.Progress = progress
	if status != "" {
		item.Status = status
	}
	state.Details[quality] = item
}

func (s *TranscodeService) markQualityFailed(resourceID uint, quality string) {
	s.updateProgress(resourceID, quality, 0, "fail")
}

func (s *TranscodeService) setUploadProgress(resourceID uint, pct float64, status, ossType string) {
	s.progressMu.Lock()
	defer s.progressMu.Unlock()
	state, ok := s.progress[resourceID]
	if !ok {
		return
	}
	if pct < 0 {
		pct = 0
	} else if pct > 100 {
		pct = 100
	}
	state.UploadProgress = pct
	state.UploadStatus = status
	if ossType != "" {
		state.UploadOSS = ossType
	}
}

func (s *TranscodeService) clearProgress(resourceID uint) {
	s.progressMu.Lock()
	defer s.progressMu.Unlock()
	delete(s.progress, resourceID)
}

// getResourceProgress 返回指定 resource 的进度快照（读锁安全）。
// 供 LocalTranscoder 及内部使用；返回 nil 表示未找到。
func (s *TranscodeService) getResourceProgress(resourceID uint) *resourceTranscodingProgress {
	s.progressMu.RLock()
	defer s.progressMu.RUnlock()
	state, ok := s.progress[resourceID]
	if !ok {
		return nil
	}
	// 返回浅拷贝，调用方不应修改字段
	return state
}

func (s *TranscodeService) GetVideoTranscodingProgress(videoID uint) (float64, []vo.TranscodingProgressItem, *vo.UploadProgressInfo) {
	s.progressMu.RLock()
	defer s.progressMu.RUnlock()

	details := make([]vo.TranscodingProgressItem, 0)
	var totalProgress float64
	var totalCount int
	var uploadInfo *vo.UploadProgressInfo

	for _, state := range s.progress {
		if state.VideoID != videoID {
			continue
		}
		for _, item := range state.Details {
			details = append(details, item)
			totalProgress += item.Progress
			totalCount++
		}
		if state.UploadStatus != "" {
			uploadInfo = &vo.UploadProgressInfo{
				OssType:  state.UploadOSS,
				Progress: state.UploadProgress,
				Status:   state.UploadStatus,
			}
		}
	}
	if totalCount == 0 {
		return 0, details, uploadInfo
	}

	sort.Slice(details, func(i, j int) bool {
		iw, ih, ifps, ik := parseProgressQualitySortKey(details[i].Quality)
		jw, jh, jfps, jk := parseProgressQualitySortKey(details[j].Quality)
		if ih != jh {
			return ih > jh
		}
		if iw != jw {
			return iw > jw
		}
		if ifps != jfps {
			return ifps > jfps
		}
		if ik != jk {
			return ik > jk
		}
		return details[i].Quality > details[j].Quality
	})
	return totalProgress / float64(totalCount), details, uploadInfo
}

func (s *TranscodeService) StopTranscodingAndCleanup(videoID uint) error {
	s.processMu.Lock()
	processes, exists := s.processes[videoID]
	if !exists || len(processes) == 0 {
		s.processMu.Unlock()
		return nil
	}
	processesCopy := make([]TranscodingProcess, len(processes))
	copy(processesCopy, processes)
	delete(s.processes, videoID)
	s.processMu.Unlock()

	cleanedDirs := make(map[string]struct{})
	for _, p := range processesCopy {
		if p.CancelFunc != nil {
			p.CancelFunc()
		}
		if p.Cmd != nil && p.Cmd.Process != nil {
			if err := p.Cmd.Process.Kill(); err != nil {
				utils.ErrorLog(fmt.Sprintf("【终止进程失败】PID=%d", p.PID), "transcoding", err.Error())
			}
		}
		if p.OutputDir != "" {
			if _, done := cleanedDirs[p.OutputDir]; !done {
				_ = cleanupTranscodedFilesInOutputDir(p.OutputDir)
				cleanedDirs[p.OutputDir] = struct{}{}
			}
		}
	}
	return nil
}

func (s *TranscodeService) HasTranscodingProcess(videoID uint) bool {
	s.processMu.RLock()
	defer s.processMu.RUnlock()
	return len(s.processes[videoID]) > 0
}

func (s *TranscodeService) GetTranscodingProcessCount(videoID uint) int {
	s.processMu.RLock()
	defer s.processMu.RUnlock()
	return len(s.processes[videoID])
}

// ==============================================================================
// 第七部分：FFmpeg 指令动态拼接 (解决严重代码重复与CPU饥饿问题)
// ==============================================================================

// runVideoEncodeTask 执行视频编码：参数构建委托 ffmpeg.VideoEncodeArgs，
// 进程管理与 GPU 错误检测在本地处理。
func (s *TranscodeService) runVideoEncodeTask(
	ctx context.Context, _ /*videoID*/, resourceID uint, inputFile, outputFile, quality, rate, fps, progressQuality string,
	totalDuration float64, useGpu bool, useAv1 bool, useHevc bool, _ /*cancelFunc*/ context.CancelFunc,
) error {
	args := ffmpeg.VideoEncodeArgs(inputFile, quality, rate, fps, totalDuration, useGpu, useAv1, useHevc)

	// CPU/GPU 线程分配：GPU 固定 4 线程；CPU 预留 2 核给系统，其余按 maxConcur 平分
	numCPU := runtime.NumCPU()
	threads := 4
	if !useGpu {
		threads = (numCPU - 2) / s.maxConcur
		if threads < 1 {
			threads = 1
		}
		if threads > 8 {
			threads = 8
		}
	}
	args = append(args, "-progress", "pipe:1", "-nostats",
		"-threads", strconv.Itoa(threads),
		"-y", outputFile)

	// 执行并报告进度
	stderr, err := ffmpeg.EncodeVideo(ctx, args, totalDuration, func(pct float64) {
		s.updateProgress(resourceID, progressQuality, pct, "processing")
	})
	if err != nil {
		if useGpu && (strings.Contains(stderr, "No NVENC") || strings.Contains(stderr, "nvenc")) {
			return fmt.Errorf("GPU error: %s", stderr)
		}
		return err
	}
	return nil
}

// encodeAudioOnly 与视频轨共用同一内容基准时长 durationSec（来自 ProcessVideoInfo）。
// presentationLeadMs>0 时：adelay 在头部插入静声（https://ffmpeg.org/ffmpeg-filters.html#adelay）；
// 输出 -t 使用 durationSec + leadSec，与 ffmpeg「-t 作输出选项」语义一致，为前置静声留出时间轴，
// 避免原先固定 durationSec 截断导致挤掉末尾约 leadMs 的有效采样。视频轨仅用 durationSec。
func encodeAudioOnly(ctx context.Context, inputFile, outputFile string, audioBitRate, audioSampleRate, audioChannels int, durationSec float64, presentationLeadMs int) error {
	args := ffmpeg.AudioEncodeArgs(inputFile, audioBitRate, audioSampleRate, audioChannels, durationSec, presentationLeadMs)
	args = append(args, "-y", outputFile)
	err := ffmpeg.EncodeAudio(ctx, args)
	if err != nil {
		return fmt.Errorf("音频编码失败: %w", err)
	}
	return nil
}

// ==============================================================================
// 第七部分 B：多音轨辅助函数
// ==============================================================================

// encodeAudioTrack 编码指定音轨（通过 streamIndex 选择流），输出唯一文件名。
func encodeAudioTrack(ctx context.Context, inputFile, outputFile string, streamIndex, audioBitRate, audioSampleRate, audioChannels int, durationSec float64, presentationLeadMs int) error {
	args := ffmpeg.AudioTrackEncodeArgs(inputFile, streamIndex, audioBitRate, audioSampleRate, audioChannels, durationSec, presentationLeadMs)
	args = append(args, "-y", outputFile)
	err := ffmpeg.EncodeAudio(ctx, args)
	if err != nil {
		return fmt.Errorf("音轨%d编码失败: %w", streamIndex, err)
	}
	return nil
}

// saveAudioTrackRecords 保存所有音轨记录到 audio_track 表。
// 默认音轨（第一条）同时写入 VideoIndexFile 的音频字段（向后兼容）。
func (s *TranscodeService) saveAudioTrackRecords(ctx context.Context, info *dto.TranscodingInfo) error {
	if len(info.AudioStreams) == 0 {
		return nil
	}

	for i, stream := range info.AudioStreams {
		audioFile := ffmpeg.AudioFileNameForTrack(stream.Language)
		audioFilePath := info.OutputDir + audioFile

		aInit, aIndex, err := getMP4InitRange(audioFilePath)
		if err != nil {
			utils.ErrorLog(fmt.Sprintf("【音轨索引解析失败】语言=%s: %s", stream.Language, err.Error()), "transcoding", "")
			continue
		}

		title := languageToTitle(stream.Language)
		track := model.AudioTrack{
			ResourceID: info.ResourceID,
			DirName:    info.DirName,
			Language:   stream.Language,
			Title:      title,
			TrackIndex: stream.StreamIndex,
			IsDefault:  i == 0,
			Channels:   stream.Channels,
			AudioFile:  audioFile,
			Codec:      ffmpeg.DefaultAudioCodec,
			Bandwidth:  stream.BitRate,
			SampleRate: stream.SampleRate,
			InitRange:  aInit,
			IndexRange: aIndex,
		}

		if err := s.db.WithContext(ctx).Where("resource_id = ? AND language = ?", info.ResourceID, stream.Language).
			FirstOrCreate(&track).Error; err != nil {
			utils.ErrorLog(fmt.Sprintf("【音轨记录保存失败】语言=%s: %s", stream.Language, err.Error()), "transcoding", "")
		}
	}
	return nil
}

// languageToTitle 将 ISO 639-2 语言代码转为可读标题
func languageToTitle(lang string) string {
	switch lang {
	case "jpn":
		return "日语"
	case "eng":
		return "英语"
	case "kor":
		return "韩语"
	case "chi", "zho":
		return "中文"
	case "fre", "fra":
		return "法语"
	case "ger", "deu":
		return "德语"
	case "spa":
		return "西班牙语"
	case "rus":
		return "俄语"
	case "tha":
		return "泰语"
	case "vie":
		return "越南语"
	case "und":
		return "未知"
	default:
		return lang
	}
}

// ==============================================================================
// 第八部分：辅助工具函数 (纯函数)
// ==============================================================================

func parseProbDuration(s string) float64 {
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v <= 0 {
		return 0
	}
	return v
}

// minEncodeDurationSeconds 取 ffprobe 中 format / 视频轨 / 音频轨 时长的**最小正数**，
// 避免按单一轨时长编码导致另一轨提前结束或片尾 Silent 拉长（分离轨时表现为 A/V 时长差）。
func minEncodeDurationSeconds(formatDur string, videoStream, audioStream *global.Streams) float64 {
	var candidates []float64
	if d := parseProbDuration(formatDur); d > 0 {
		candidates = append(candidates, d)
	}
	if videoStream != nil {
		if d := parseProbDuration(videoStream.Duration); d > 0 {
			candidates = append(candidates, d)
		}
	}
	if audioStream != nil {
		if d := parseProbDuration(audioStream.Duration); d > 0 {
			candidates = append(candidates, d)
		}
	}
	if len(candidates) == 0 {
		return 0
	}
	m := candidates[0]
	for _, d := range candidates[1:] {
		if d < m {
			m = d
		}
	}
	return m
}

// 【修复修复】防止路径穿越引发目录安全隐患
func cleanupTranscodedFilesInOutputDir(outputDir string) error {
	files, err := os.ReadDir(outputDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		if strings.HasPrefix(fileName, "upload.") || fileName == "cover.jpg" {
			continue
		}
		// 使用 filepath.Base 规避可能的路径穿越
		safePath := filepath.Join(outputDir, filepath.Base(fileName))
		if err := os.Remove(safePath); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func updateVideoFileStatus(db *gorm.DB, currentResource model.Resource, resourceId uint) {
	if currentResource.FileID != 0 {
		db.Model(&model.VideoFile{}).Where("id = ? AND status != ?", currentResource.FileID, model.FileStatusReady).Update("status", model.FileStatusReady)
	} else {
		var videoIndex model.VideoIndexFile
		if err := db.Where("resource_id = ?", resourceId).First(&videoIndex).Error; err == nil && videoIndex.DirName != "" {
			res := db.Model(&model.VideoFile{}).Where("dir_name = ? AND status != ?", videoIndex.DirName, model.FileStatusReady).Update("status", model.FileStatusReady)
			if res.RowsAffected > 0 {
				var vf model.VideoFile
				if err := db.Where("dir_name = ?", videoIndex.DirName).First(&vf).Error; err == nil {
					db.Model(&model.Resource{}).Where("id = ? AND file_id = 0", resourceId).Update("file_id", vf.ID)
				}
			}
		}
	}
}

func GenerateCover(inputFile, outputFile string) error {
	_, err := utils.RunCmd(exec.Command("ffmpeg", "-i", inputFile, "-vframes", "1", "-y", outputFile))
	return err
}

func ProcessVideoInfo(input string) (*dto.TranscodingInfo, error) {
	ti := dto.TranscodingInfo{OriginalVideoStatus: -1}
	videoData, err := getVideoInfo(input)
	if err != nil {
		return &ti, err
	}

	var videoStream *global.Streams
	audioStreams := make([]global.Streams, 0)
	for i := range videoData.Stream {
		switch videoData.Stream[i].CodecType {
		case "video":
			if videoStream == nil {
				videoStream = &videoData.Stream[i]
			}
		case "audio":
			audioStreams = append(audioStreams, videoData.Stream[i])
		}
	}
	if videoStream == nil && len(videoData.Stream) > 0 {
		videoStream = &videoData.Stream[0]
	}
	if videoStream == nil {
		return &ti, fmt.Errorf("未找到视频流: %s", input)
	}

	ti.Width = videoStream.Width
	ti.Height = videoStream.Height
	ti.CodecName = videoStream.CodecName

	// 【HDR/10-bit 检测】探测源视频位深，在转码为 8-bit 时告警
	if is10BitPixelFormat(videoStream.PixFmt) {
		isTarget8Bit := !global.Config.Transcoding.UseAv1 && !global.Config.Transcoding.UseH265
		if isTarget8Bit {
			utils.InfoLog(fmt.Sprintf("【HDR告警】源视频为 %s (%s)，但转码目标为 8-bit H.264，HDR/10-bit 色彩信息将丢失",
				videoStream.CodecName, videoStream.PixFmt), "transcoding")
		} else {
			utils.InfoLog(fmt.Sprintf("【HDR检测】源视频为 %s (%s)，转码目标支持高位深，保留 HDR", videoStream.CodecName, videoStream.PixFmt), "transcoding")
		}
	}

	// 【多音轨】取第一个音轨作为默认音频参数（向后兼容）
	var primaryAudio *global.Streams
	if len(audioStreams) > 0 {
		primaryAudio = &audioStreams[0]
	}

	ti.Duration = minEncodeDurationSeconds(videoData.Format.Duration, videoStream, primaryAudio)
	if ti.Duration <= 0 {
		durStr := videoStream.Duration
		if durStr == "" {
			durStr = videoData.Format.Duration
		}
		ti.Duration, _ = strconv.ParseFloat(durStr, 64)
	}
	ti.FPS = videoStream.AvgFrameRate
	ti.FPS30, ti.FPS60 = ffmpeg.GetFpsInfo(ti.FPS, global.Config.Transcoding.Generate1080p60)

	if br, err := strconv.Atoi(videoStream.BitRate); err == nil && br > 0 {
		ti.VideoBitRate = br
	}

	// 【多音轨】填充所有音频流
	ti.AudioStreams = make([]dto.AudioStreamProbe, 0, len(audioStreams))
	for idx, a := range audioStreams {
		lang := a.Tags.Language
		if lang == "" {
			lang = "und"
		}
		sr, _ := strconv.Atoi(a.SampleRate)
		br, _ := strconv.Atoi(a.BitRate)
		ti.AudioStreams = append(ti.AudioStreams, dto.AudioStreamProbe{
			StreamIndex: idx,
			Language:    lang,
			SampleRate:  sr,
			Channels:    a.Channels,
			BitRate:     br,
		})
	}

	// 默认音频参数仍取第一个音轨（向后兼容）
	if primaryAudio != nil {
		if sr, err := strconv.Atoi(primaryAudio.SampleRate); err == nil && sr > 0 {
			ti.AudioSampleRate = sr
		}
		if primaryAudio.Channels > 0 {
			ti.AudioChannels = primaryAudio.Channels
		}
		if br, err := strconv.Atoi(primaryAudio.BitRate); err == nil && br > 0 {
			ti.AudioBitRate = br
		}
	}

	if ti.AudioSampleRate == 0 {
		ti.AudioSampleRate = ffmpeg.DefaultAudioSample
	}
	if ti.AudioChannels == 0 {
		ti.AudioChannels = ffmpeg.DefaultAudioChan
	}
	if ti.AudioBitRate == 0 {
		ti.AudioBitRate = ffmpeg.AudioMaxBitrateBps
	}
	if ti.AudioBitRate > ffmpeg.AudioMaxBitrateBps {
		ti.AudioBitRate = ffmpeg.AudioMaxBitrateBps
	}

	if ti.VideoBitRate == 0 {
		if tbr, err := strconv.Atoi(videoData.Format.BitRate); err == nil && tbr > 0 {
			ti.VideoBitRate = tbr - ti.AudioBitRate
			if ti.VideoBitRate <= 0 {
				ti.VideoBitRate = tbr
			}
		}
	}
	if ti.VideoBitRate == 0 {
		maxLvl := ffmpeg.GetMaxQualityLevel(ti.Width, ti.Height)
		ti.VideoBitRate = ffmpeg.GetDefaultVideoBitRateByLevel(maxLvl, ti.Width < ti.Height)
	}
	return &ti, nil
}

func getMP4InitRange(filePath string) (initRange, indexRange string, err error) {
	return ffmpeg.GetMP4InitRange(filePath)
}
func getVideoInfo(input string) (global.VideoInfo, error) {
	result, err := ffmpeg.RunFFprobe(input)
	if err != nil {
		return global.VideoInfo{}, err
	}
	// 转换 *ffmpeg.ProbeResult → global.VideoInfo
	vi := global.VideoInfo{
		Stream: make([]global.Streams, len(result.Streams)),
		Format: global.Format{
			BitRate:  result.Format.BitRate,
			Duration: result.Format.Duration,
		},
	}
	for i, s := range result.Streams {
		vi.Stream[i] = global.Streams{
			CodecType:    s.CodecType,
			CodecName:    s.CodecName,
			Width:        s.Width,
			Height:       s.Height,
			PixFmt:       s.PixFmt,
			Duration:     s.Duration,
			RFrameRate:   s.RFrameRate,
			AvgFrameRate: s.AvgFrameRate,
			SampleRate:   s.SampleRate,
			Channels:     s.Channels,
			BitRate:      s.BitRate,
			Tags: global.StreamTags{
				Language: s.Tags.Language,
				Title:    s.Tags.Title,
			},
		}
	}
	return vi, nil
}

func ProbeH264Avc1CodecString(filePath string) (string, error) {
	return probeH264Avc1CodecString(filePath)
}

func probeH264Avc1CodecString(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=profile,level", "-of", "default=nw=1:nk=1", filePath)
	out, err := utils.RunCmd(cmd)
	if err != nil {
		return "", err
	}
	raw := strings.TrimSpace(out.String())
	lines := strings.Split(raw, "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("ffprobe output unexpected")
	}
	profile := strings.ToLower(strings.TrimSpace(lines[0]))
	level, _ := strconv.Atoi(strings.TrimSpace(lines[1]))

	return avc1CodecStringFromH264ProfileLevel(profile, level)
}

func avc1CodecStringFromH264ProfileLevel(profile string, level int) (string, error) {
	s := ffmpeg.Avc1CodecString(profile, level)
	if s == "" {
		return "", fmt.Errorf("unsupported h264 profile: %s", profile)
	}
	return s, nil
}

func parseProgressQualitySortKey(quality string) (w, h, f, b int) {
	parts := strings.Split(quality, "_")
	if len(parts) > 0 {
		res := strings.Split(parts[0], "x")
		if len(res) == 2 {
			w, _ = strconv.Atoi(res[0])
			h, _ = strconv.Atoi(res[1])
		}
	}
	if len(parts) > 1 {
		b = ffmpeg.ParseBitrateKbps(parts[1])
	}
	if len(parts) > 2 {
		f, _ = strconv.Atoi(parts[2])
	}
	return
}

// is10BitPixelFormat 判断 ffprobe pix_fmt 是否为 10-bit 及以上位深
func is10BitPixelFormat(pixFmt string) bool {
	switch pixFmt {
	case "yuv420p10le", "yuv422p10le", "yuv444p10le",
		"yuv420p12le", "yuv422p12le", "yuv444p12le",
		"yuva420p10le", "yuva422p10le", "yuva444p10le",
		"gbrp10le", "gbrp12le",
		"rgb10le", "rgb12le":
		return true
	}
	return false
}
