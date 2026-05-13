package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math"
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
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

// ==============================================================================
// 第一部分：全局常量定义
// ==============================================================================

const (
	AudioMaxBitrateBps = 320000
	DefaultAudioSample = 48000
	DefaultAudioChan   = 2
	FragDurationUs     = "2000000" // 2秒切片
	TimeBase30fps      = "30000/1000"
	TimeBase60fps      = "60000/1001"
	DefaultVideoCodec  = "avc1.640028"
	DefaultAudioCodec  = "mp4a.40.2"
)

// ==============================================================================
// 第二部分：接口抽象与核心服务结构 (解决强耦合与可测试性问题)
// ==============================================================================

// StorageUploader 抽象存储接口，屏蔽具体 OSS 实现，便于 Mock 测试
type StorageUploader interface {
	PutObjectFromFile(objectKey, filePath string) error
}

type TranscodingTarget struct {
	Resolution  string
	BitrateRate string
	FPS         string
	FpsName     string
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
)

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

func ResetGPUState() { GetTranscoder().ResetGPUState() }

func GetVideoTranscodingProgress(videoID uint) (float64, []vo.TranscodingProgressItem) {
	return GetTranscoder().GetVideoTranscodingProgress(videoID)
}

func VideoTransCoding(transcodingInfo *dto.TranscodingInfo) {
	GetTranscoder().ProcessVideo(context.Background(), transcodingInfo)
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
	targets := getTranscodingTarget(info)
	s.initResourceProgress(info.VideoID, info.ResourceID, targets)

	utils.InfoLog(fmt.Sprintf("【转码开始】VideoID=%d, ResourceID=%d, 目标数量=%d", info.VideoID, info.ResourceID, len(targets)), "transcoding")

	var wg sync.WaitGroup
	var successCount atomic.Int32
	audioTask := NewSharedTask()

	// 1. 并发处理所有分片
	for _, target := range targets {
		wg.Add(1)
		go func(t TranscodingTarget) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					qualityName := t.Resolution + "_" + t.BitrateRate + "_" + t.FpsName
					s.markQualityFailed(info.ResourceID, qualityName)
					utils.ErrorLog(fmt.Sprintf("【Goroutine panic】%s: %v", qualityName, r), "transcoding", "")
				}
			}()

			if err := s.processSingleQuality(ctx, info, t, audioTask); err != nil {
				utils.ErrorLog(fmt.Sprintf("【转码处理失败】%s", err.Error()), "transcoding", "")
				return
			}
			successCount.Add(1)
		}(target)
	}

	wg.Wait()
	utils.InfoLog(fmt.Sprintf("【所有分段转码完成】成功=%d, 总数=%d", successCount.Load(), len(targets)), "transcoding")

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
			_ = s.completeTransaction(ctx, info, global.PROCESSING_FAIL)
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
}

func (s *TranscodeService) processSingleQuality(
	ctx context.Context, info *dto.TranscodingInfo, t TranscodingTarget, audioTask *SharedTask,
) error {
	qualityName := t.Resolution + "_" + t.BitrateRate + "_" + t.FpsName
	videoFile := info.OutputDir + qualityName + "_video.m4s"
	audioFile := info.OutputDir + "audio.m4s"

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

	// 音频同步处理
		err := audioTask.DoOrWait(subCtx, func() error {
			utils.InfoLog(fmt.Sprintf("【编码音频】audio.m4s (码率=%dk)", info.AudioBitRate/1000), "transcoding")
			leadMs := bFramePresentationLeadMs(t.FPS)
			return encodeAudioOnly(subCtx, info.InputFile, audioFile, info.AudioBitRate, info.AudioSampleRate, info.AudioChannels, info.Duration, leadMs)
		})
	if err != nil {
		s.markQualityFailed(info.ResourceID, qualityName)
		return fmt.Errorf("【音频等待失败】%s: %w", qualityName, err)
	}

	// 解析存储
	if err := s.saveIndexRecord(subCtx, info, qualityName, videoFile, audioFile); err != nil {
		s.markQualityFailed(info.ResourceID, qualityName)
		return fmt.Errorf("【索引生成失败】%s: %w", qualityName, err)
	}

	s.updateProgress(info.ResourceID, qualityName, 100, "success")
	return nil
}

// encodeVideoWithFallback 动态路由 GPU/CPU 策略
func (s *TranscodeService) encodeVideoWithFallback(
	ctx context.Context, info *dto.TranscodingInfo, t TranscodingTarget,
	videoFile, qualityName string, cancel context.CancelFunc,
) error {
	useGpu := global.Config.Transcoding.UseGpu && s.isGPUAvailable()
	utils.InfoLog(fmt.Sprintf("【开始视频编码】%s 使用GPU=%v", qualityName, useGpu), "transcoding")

	err := s.runVideoEncodeTask(ctx, info.VideoID, info.ResourceID, info.InputFile, videoFile, t.Resolution, t.BitrateRate, t.FPS, qualityName, info.Duration, useGpu, cancel)

	if useGpu && err != nil && (strings.Contains(err.Error(), "GPU error") || strings.Contains(err.Error(), "nvenc")) {
		s.handleGPUFailure()
		if ctx.Err() != nil {
			return ctx.Err()
		}
		utils.InfoLog(fmt.Sprintf("【降级CPU编码】%s", qualityName), "transcoding")
		return s.runVideoEncodeTask(ctx, info.VideoID, info.ResourceID, info.InputFile, videoFile, t.Resolution, t.BitrateRate, t.FPS, qualityName, info.Duration, false, cancel)
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

	videoCodec := DefaultVideoCodec
	if c, err := probeH264Avc1CodecString(videoFilePath); err == nil {
		videoCodec = c
	}

	width, height, bandwidth, frameRate := parseQualityInfo(qualityName)

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
		FrameRate:       frameRate,
		VideoInitRange:  vInit,
		VideoIndexRange: vIndex,
		AudioFile:       "audio.m4s",
		AudioBandwidth:  info.AudioBitRate,
		AudioCodec:      DefaultAudioCodec,
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

		if targetStatus == global.WAITING_REVIEW && info.OriginalVideoStatus == global.AUDIT_APPROVED {
			targetStatus = global.AUDIT_APPROVED
		}

		if err := tx.Model(&model.Resource{}).
			Where("id = ? and status = ?", info.ResourceID, global.VIDEO_PROCESSING).
			Update("status", targetStatus).Error; err != nil {
			return err
		}

		if targetStatus != global.PROCESSING_FAIL {
			updateVideoFileStatus(tx, currentResource, info.ResourceID)
		}

		var pendingCount int64
		tx.Model(&model.Resource{}).Where("vid = ? and status = ? and id != ?", info.VideoID, global.VIDEO_PROCESSING, info.ResourceID).Count(&pendingCount)
		if pendingCount > 0 {
			return nil
		}

		var totalCount, failedCount int64
		tx.Model(&model.Resource{}).Where("vid = ?", info.VideoID).Count(&totalCount)
		tx.Model(&model.Resource{}).Where("vid = ? and status = ?", info.VideoID, global.PROCESSING_FAIL).Count(&failedCount)

		finalVideoStatus := global.WAITING_REVIEW
		if failedCount == totalCount {
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

func (s *TranscodeService) uploadToOSS(ctx context.Context, info *dto.TranscodingInfo) error {
	files, err := os.ReadDir(info.OutputDir)
	if err != nil {
		return err
	}

	tasks := make(chan os.DirEntry, len(files))
	results := make(chan bool, len(files))

	var wg sync.WaitGroup
	for i := 0; i < ossUploadMaxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range tasks {
				fileName := file.Name()
				if fileName == "upload"+info.Suffix && !global.Config.Storage.UploadMp4File {
					results <- false
					continue
				}

				objectKey := "video/" + info.DirName + "/" + fileName
				filePath := info.OutputDir + fileName

				err := s.uploader.PutObjectFromFile(objectKey, filePath)
				if err != nil {
					// 使用 context 感知的非阻塞等待进行重试
					select {
					case <-ctx.Done():
						results <- false
						return
					case <-time.After(ossUploadRetryDelay):
						err = s.uploader.PutObjectFromFile(objectKey, filePath)
					}
				}

				if err != nil {
					utils.ErrorLog(fmt.Sprintf("【OSS重试后仍失败】%s", fileName), "oss", err.Error())
					results <- false
				} else {
					results <- true
				}
			}
		}()
	}

	for _, file := range files {
		if !file.IsDir() {
			tasks <- file
		}
	}
	close(tasks)
	wg.Wait()
	close(results)

	successCount := 0
	for success := range results {
		if success {
			successCount++
		}
	}

	if successCount == 0 {
		return fmt.Errorf("OSS全部上传失败")
	}
	return nil
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

func (s *TranscodeService) initResourceProgress(videoID, resourceID uint, targets []TranscodingTarget) {
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

func (s *TranscodeService) clearProgress(resourceID uint) {
	s.progressMu.Lock()
	defer s.progressMu.Unlock()
	delete(s.progress, resourceID)
}

func (s *TranscodeService) GetVideoTranscodingProgress(videoID uint) (float64, []vo.TranscodingProgressItem) {
	s.progressMu.RLock()
	defer s.progressMu.RUnlock()

	details := make([]vo.TranscodingProgressItem, 0)
	var totalProgress float64
	var totalCount int

	for _, state := range s.progress {
		if state.VideoID != videoID {
			continue
		}
		for _, item := range state.Details {
			details = append(details, item)
			totalProgress += item.Progress
			totalCount++
		}
	}
	if totalCount == 0 {
		return 0, details
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
	return totalProgress / float64(totalCount), details
}

func (s *TranscodeService) registerProcess(videoID, resourceID uint, cmd *exec.Cmd, cancelFunc context.CancelFunc, outputDir string) {
	s.processMu.Lock()
	defer s.processMu.Unlock()
	s.processes[videoID] = append(s.processes[videoID], TranscodingProcess{
		VideoID:    videoID,
		ResourceID: resourceID,
		PID:        cmd.Process.Pid,
		Cmd:        cmd,
		CancelFunc: cancelFunc,
		OutputDir:  outputDir,
	})
}

func (s *TranscodeService) unregisterProcess(videoID uint, pid int) {
	s.processMu.Lock()
	defer s.processMu.Unlock()
	processes, exists := s.processes[videoID]
	if !exists {
		return
	}
	newProcesses := make([]TranscodingProcess, 0)
	for _, p := range processes {
		if p.PID != pid {
			newProcesses = append(newProcesses, p)
		}
	}
	if len(newProcesses) == 0 {
		delete(s.processes, videoID)
	} else {
		s.processes[videoID] = newProcesses
	}
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
			_ = p.Cmd.Process.Kill()
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

// runVideoEncodeTask 动态合并构建 CPU 与 GPU 编码指令，消除 DRY 违规。
// 帧率模式使用 -fps_mode cfr（FFmpeg 文档：-vsync 已弃用，宜用 -fps_mode）。
// 输出时长使用与音轨相同的 -t，见 ffmpegOutputDurationArgs。
func (s *TranscodeService) runVideoEncodeTask(
	ctx context.Context, videoID, resourceID uint, inputFile, outputFile, quality, rate, fps, progressQuality string,
	totalDuration float64, useGpu bool, cancelFunc context.CancelFunc,
) error {
	fpsFloat := parseFPS(fps)
	gopSize := int(math.Round(fpsFloat * 2))
	if gopSize < 1 {
		gopSize = 60
	}
	gopSizeStr := strconv.Itoa(gopSize)
	targetRate, maxrate, bufsize := buildRateControlParams(rate)
	// fps 过滤器强制精确 CFR（修复 59.94 源→30 输出时 avg_fps 漂移到 28.5 的问题）
	scaleFilter := fmt.Sprintf("fps=%s,scale=%s:flags=lanczos", fps, quality)

	// 【核心修复】防止 FFmpeg 默认吃光所有 CPU 引发 OS 上下文切换风暴
	numCPU := runtime.NumCPU()
	threads := numCPU / s.maxConcur
	if threads < 1 {
		threads = 1
	}
	if threads > 8 {
		threads = 8
	} // 限制单实例线程上限避免内部损耗

	args := []string{
		"-i", inputFile,
		"-filter_complex", fmt.Sprintf("[0:v]setpts=PTS-STARTPTS,%s", scaleFilter),
		"-an",
	}

	// 保留 B 帧；首帧展示时刻相对 t=0 常延后约 2 帧（参见 ffprobe stream start_time）。
	// 音轨侧用 adelay 注入等量前置静声与视频首画对齐（ffmpeg 滤镜 adelay）。
	if useGpu {
		args = append(args,
			"-c:v", "h264_nvenc", "-cq", "23", "-preset", "p4", "-rc", "vbr",
			"-profile:v", "high", "-pix_fmt", "yuv420p", "-bf", "2", "-b_ref_mode", "middle",
			"-forced-idr", "1",
		)
	} else {
		args = append(args,
			"-c:v", "libx264", "-preset", "medium", "-crf", "20", "-tune", "film",
			"-profile:v", "high", "-pix_fmt", "yuv420p", "-bf", "3", "-b_strategy", "2",
			"-flags", "+cgop", "-sc_threshold", "0",
			"-refs", "6", "-me_method", "hex", "-subq", "9",
		)
	}

	args = append(args,
		"-b:v", targetRate, "-maxrate", maxrate, "-bufsize", bufsize,
		"-r", fps, "-g", gopSizeStr, "-keyint_min", gopSizeStr,
		"-threads", strconv.Itoa(threads),
	)

	if useGpu {
		args = append(args,
			// NVENC：与 CFR GOP 对齐，减少与 -fps_mode cfr 组合的边界抖动（参见 FFmpeg h264_nvenc 文档 strict_gop）
			"-strict_gop", "1",
			"-delay", "0",
		)
	}

	// cfr：按请求帧率重复/丢桢，与前置 fps 滤镜共同保证恒定帧率（ffmpeg.html「fps_mode / cfr」）
	args = append(args,
		"-fps_mode", "cfr",
		"-progress", "pipe:1", "-nostats",
		"-f", "mp4",
		"-frag_duration", FragDurationUs,
		"-movflags", "+frag_keyframe+empty_moov+default_base_moof+dash+global_sidx+negative_cts_offsets",
		"-avoid_negative_ts", "make_zero",
	)
	if dur := ffmpegOutputDurationArgs(totalDuration); dur != nil {
		args = append(args, dur...)
	}
	args = append(args, "-y", outputFile)

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go watchFFmpegProgress(bufio.NewScanner(stdout), resourceID, progressQuality, totalDuration, s)

	s.registerProcess(videoID, resourceID, cmd, cancelFunc, filepath.Dir(outputFile))
	defer s.unregisterProcess(videoID, cmd.Process.Pid)

	if err := cmd.Wait(); err != nil {
		errOut := stderr.String()
		if useGpu && (strings.Contains(errOut, "No NVENC") || strings.Contains(errOut, "nvenc")) {
			return fmt.Errorf("GPU error: %s", errOut)
		}
		return fmt.Errorf("encode failed: %s\n%s", err.Error(), errOut)
	}
	return nil
}

// encodeAudioOnly 与视频轨共用同一内容基准时长 durationSec（来自 ProcessVideoInfo）。
// presentationLeadMs>0 时：adelay 在头部插入静声（https://ffmpeg.org/ffmpeg-filters.html#adelay）；
// 输出 -t 使用 durationSec + leadSec，与 ffmpeg「-t 作输出选项」语义一致，为前置静声留出时间轴，
// 避免原先固定 durationSec 截断导致挤掉末尾约 leadMs 的有效采样。视频轨仅用 durationSec。
func encodeAudioOnly(ctx context.Context, inputFile, outputFile string, audioBitRate, audioSampleRate, audioChannels int, durationSec float64, presentationLeadMs int) error {
	bitRateStr := fmt.Sprintf("%dk", audioBitRate/1000)
	sampleRateStr := strconv.Itoa(audioSampleRate)
	channelsStr := strconv.Itoa(audioChannels)

	adelayArg := adelayPerChannelArg(presentationLeadMs, audioChannels)
	audioFilter := fmt.Sprintf("[0:a]asetpts=PTS-STARTPTS,aresample=osr=%s", sampleRateStr)
	if adelayArg != "" {
		audioFilter += ",adelay=" + adelayArg
	}
	audioFilter += "[aout]"
	audioOutputDur := durationSec
	if presentationLeadMs > 0 {
		audioOutputDur += float64(presentationLeadMs) / 1000.0
	}
	command := []string{
		"-i", inputFile,
		"-filter_complex", audioFilter,
		"-map", "[aout]", "-vn", "-c:a", "aac", "-b:a", bitRateStr,
		"-ar", sampleRateStr, "-ac", channelsStr,
		"-f", "mp4",
		"-frag_duration", FragDurationUs,
		"-movflags", "+frag_keyframe+empty_moov+default_base_moof+dash+global_sidx+negative_cts_offsets",
		"-avoid_negative_ts", "make_zero",
	}
	if dur := ffmpegOutputDurationArgs(audioOutputDur); dur != nil {
		command = append(command, dur...)
	}
	command = append(command, "-y", outputFile)
	cmd := exec.CommandContext(ctx, "ffmpeg", command...)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("音频编码失败: %s, stderr: %s", err.Error(), stderr.String())
	}
	return nil
}

// ==============================================================================
// 第八部分：辅助工具函数 (纯函数)
// ==============================================================================

// ffmpegOutputDurationArgs 在写出文件前附加「-t duration」作为**输出选项**。
// 官方说明见 https://ffmpeg.org/ffmpeg.html 中「-t duration (input/output)」：
// 作为输出选项时，表示输出时长达到给定值后停止写入。音视频两次编码传入同一 duration
// 可显著收敛分离封装（fMP4 + DASH）下片尾时长漂移。
func ffmpegOutputDurationArgs(seconds float64) []string {
	if seconds <= 0 || math.IsNaN(seconds) || math.IsInf(seconds, 0) {
		return nil
	}
	// 与 ffmpeg 时间语法兼容的十进制秒（去掉无意义尾随 0）
	s := strings.TrimRight(strings.TrimRight(strconv.FormatFloat(seconds, 'f', 6, 64), "0"), ".")
	if s == "" {
		return nil
	}
	return []string{"-t", s}
}

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

// bFramePresentationLeadMs 按目标帧率估计视频首帧「展示」相对 t=0 的毫秒延迟（与 IPB 下约 2 帧同量纲）。
func bFramePresentationLeadMs(fps string) int {
	f := parseFPS(fps)
	if f <= 0 {
		return 0
	}
	ms := int(math.Round(2000.0 / f))
	if ms < 1 {
		return 0
	}
	if ms > 200 {
		return 200
	}
	return ms
}

// adelayPerChannelArg 生成 FFmpeg adelay 的「del0|del1|…」语法（各声道相等）。
func adelayPerChannelArg(delayMs, channels int) string {
	if delayMs <= 0 || channels < 1 {
		return ""
	}
	s := strconv.Itoa(delayMs)
	parts := make([]string, channels)
	for i := range parts {
		parts[i] = s
	}
	return strings.Join(parts, "|")
}

func watchFFmpegProgress(scanner *bufio.Scanner, resourceID uint, quality string, totalDuration float64, svc *TranscodeService) {
	if totalDuration <= 0 {
		return
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "out_time_ms=") {
			raw := strings.TrimPrefix(line, "out_time_ms=")
			if ms, err := strconv.ParseFloat(raw, 64); err == nil {
				svc.updateProgress(resourceID, quality, (ms/1000000.0)/totalDuration*100, "processing")
			}
		} else if strings.HasPrefix(line, "out_time=") {
			raw := strings.TrimPrefix(line, "out_time=")
			seconds := parseFfmpegClockToSeconds(raw)
			svc.updateProgress(resourceID, quality, seconds/totalDuration*100, "processing")
		}
	}
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

	var videoStream, audioStream *global.Streams
	for i := range videoData.Stream {
		if videoData.Stream[i].CodecType == "video" && videoStream == nil {
			videoStream = &videoData.Stream[i]
		} else if videoData.Stream[i].CodecType == "audio" && audioStream == nil {
			audioStream = &videoData.Stream[i]
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

	ti.Duration = minEncodeDurationSeconds(videoData.Format.Duration, videoStream, audioStream)
	if ti.Duration <= 0 {
		durStr := videoStream.Duration
		if durStr == "" {
			durStr = videoData.Format.Duration
		}
		ti.Duration, _ = strconv.ParseFloat(durStr, 64)
	}
	ti.FPS = videoStream.AvgFrameRate
	ti.FPS30, ti.FPS60 = getFpsInfo(ti.FPS)

	if br, err := strconv.Atoi(videoStream.BitRate); err == nil && br > 0 {
		ti.VideoBitRate = br
	}

	if audioStream != nil {
		if sr, err := strconv.Atoi(audioStream.SampleRate); err == nil && sr > 0 {
			ti.AudioSampleRate = sr
		}
		if audioStream.Channels > 0 {
			ti.AudioChannels = audioStream.Channels
		}
		if br, err := strconv.Atoi(audioStream.BitRate); err == nil && br > 0 {
			ti.AudioBitRate = br
		}
	}

	if ti.AudioSampleRate == 0 {
		ti.AudioSampleRate = DefaultAudioSample
	}
	if ti.AudioChannels == 0 {
		ti.AudioChannels = DefaultAudioChan
	}
	if ti.AudioBitRate == 0 {
		ti.AudioBitRate = AudioMaxBitrateBps
	}
	if ti.AudioBitRate > AudioMaxBitrateBps {
		ti.AudioBitRate = AudioMaxBitrateBps
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
		maxLvl := getMaxQualityLevel(ti.Width, ti.Height)
		ti.VideoBitRate = getDefaultVideoBitRateByLevel(maxLvl, ti.Width < ti.Height)
	}
	return &ti, nil
}

func getMP4InitRange(filePath string) (initRange, indexRange string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", "", err
	}
	fileSize := fileInfo.Size()

	var moovEnd, sidxStart, sidxEnd int64
	offset := int64(0)
	for offset < fileSize {
		boxSize, boxType, err := readMP4BoxSizeAndType(file, offset, fileSize)
		if err != nil || boxSize <= 0 {
			break
		}
		switch boxType {
		case "moov":
			moovEnd = offset + boxSize
		case "sidx":
			sidxStart = offset
			sidxEnd = offset + boxSize
		case "moof":
			goto ParseDone
		}
		offset += boxSize
		if sidxEnd > 0 {
			break
		}
	}
ParseDone:
	if moovEnd == 0 {
		return "", "", fmt.Errorf("moov box not found")
	}
	initRange = fmt.Sprintf("0-%d", moovEnd-1)
	if sidxEnd > 0 {
		indexRange = fmt.Sprintf("%d-%d", sidxStart, sidxEnd-1)
	} else {
		indexRange = fmt.Sprintf("%d-%d", moovEnd, fileSize-1)
	}
	return initRange, indexRange, nil
}

type qualityPreset struct {
	LongSide   int
	ShortSide  int
	BitrateH   string
	BitrateV   string
	Bitrate60H string
	Bitrate60V string
}

var qualityPresets = []qualityPreset{
	{1920, 1080, "8000k", "5000k", "12000k", "8000k"},
	{1280, 720, "5000k", "3000k", "7500k", "5000k"},
	{854, 480, "2500k", "1500k", "", ""},
	{640, 360, "1000k", "700k", "", ""},
}

func getMaxQualityLevel(width, height int) int {
	shortSide, longSide := width, height
	if height < width {
		shortSide, longSide = height, width
	}
	if longSide >= 1920 && shortSide >= 1000 {
		return 1080
	}
	if shortSide >= 1080 {
		return 1080
	}
	if shortSide >= 720 {
		return 720
	}
	if shortSide >= 480 {
		return 480
	}
	return 360
}

func getFpsInfo(avgFrameRate string) (string, string) {
	parts := strings.Split(avgFrameRate, "/")
	if len(parts) == 2 {
		num, den := utils.StringToInt(parts[0]), utils.StringToInt(parts[1])
		if den == 0 {
			return TimeBase30fps, ""
		}
		fps := float64(num) / float64(den)
		if fps < 30 {
			return avgFrameRate, ""
		}
		if fps >= 59 {
			if global.Config.Transcoding.Generate1080p60 {
				return TimeBase30fps, TimeBase60fps
			}
			return TimeBase30fps, ""
		}
	}
	return TimeBase30fps, ""
}

func calcResolution(srcWidth, srcHeight, targetShortSide int) (w, h int) {
	isPortrait := srcWidth < srcHeight
	srcAspect := float64(srcWidth) / float64(srcHeight)
	if isPortrait {
		w = targetShortSide
		h = int(math.Round(float64(w) / srcAspect))
	} else {
		h = targetShortSide
		w = int(math.Round(float64(h) * srcAspect))
	}
	return w &^ 1, h &^ 1
}

func parseBitrateKbps(rate string) int {
	rate = strings.TrimSpace(strings.TrimSuffix(rate, "k"))
	if v, err := strconv.Atoi(rate); err == nil && v > 0 {
		return v
	}
	return 0
}

func formatBitrateKbps(rateKbps int) string {
	if rateKbps < 1 {
		rateKbps = 1
	}
	return fmt.Sprintf("%dk", rateKbps)
}

func getDefaultVideoBitRateByLevel(maxLevel int, isPortrait bool) int {
	for _, preset := range qualityPresets {
		if preset.ShortSide == maxLevel {
			br := preset.BitrateH
			if isPortrait {
				br = preset.BitrateV
			}
			if kbps := parseBitrateKbps(br); kbps > 0 {
				return kbps * 1000
			}
		}
	}
	return 1000000
}

func scaleBitrateBySource(sourceMaxKbps, maxPresetKbps, currentPresetKbps int) int {
	if sourceMaxKbps <= 0 {
		return currentPresetKbps
	}
	if maxPresetKbps <= 0 || currentPresetKbps <= 0 {
		return sourceMaxKbps
	}
	if currentPresetKbps >= maxPresetKbps {
		return sourceMaxKbps
	}
	scaled := int(math.Round(float64(sourceMaxKbps) * float64(currentPresetKbps) / float64(maxPresetKbps)))
	if scaled < 200 {
		scaled = 200
	}
	if scaled > sourceMaxKbps {
		scaled = sourceMaxKbps
	}
	return scaled
}

func getPresetBitrateKbps(preset qualityPreset, isPortrait, fps60 bool) int {
	var bitrate string
	if fps60 {
		bitrate = preset.Bitrate60H
		if isPortrait {
			bitrate = preset.Bitrate60V
		}
		if kbps := parseBitrateKbps(bitrate); kbps > 0 {
			return kbps
		}
	}
	bitrate = preset.BitrateH
	if isPortrait {
		bitrate = preset.BitrateV
	}
	return parseBitrateKbps(bitrate)
}

func getTranscodingTarget(videoInfo *dto.TranscodingInfo) []TranscodingTarget {
	targets := make([]TranscodingTarget, 0)
	maxLevel := getMaxQualityLevel(videoInfo.Width, videoInfo.Height)
	isPortrait := videoInfo.Width < videoInfo.Height

	var maxPreset qualityPreset
	for _, p := range qualityPresets {
		if p.ShortSide == maxLevel {
			maxPreset = p
			break
		}
	}
	if maxPreset.ShortSide == 0 {
		return targets
	}

	maxPresetKbps30 := getPresetBitrateKbps(maxPreset, isPortrait, false)
	maxPresetKbps60 := getPresetBitrateKbps(maxPreset, isPortrait, true)
	if maxPresetKbps60 <= 0 {
		maxPresetKbps60 = maxPresetKbps30
	}

	sourceMaxKbps := videoInfo.VideoBitRate / 1000
	if sourceMaxKbps <= 0 {
		sourceMaxKbps = maxPresetKbps30
	}

	enable60 := global.Config.Transcoding.Generate1080p60 && videoInfo.FPS60 != ""

	for _, preset := range qualityPresets {
		if preset.ShortSide > maxLevel {
			continue
		}
		w, h := calcResolution(videoInfo.Width, videoInfo.Height, preset.ShortSide)
		resStr := fmt.Sprintf("%dx%d", w, h)

		currKbps30 := getPresetBitrateKbps(preset, isPortrait, false)
		dynKbps30 := scaleBitrateBySource(sourceMaxKbps, maxPresetKbps30, currKbps30)

		if preset.ShortSide == maxLevel && enable60 {
			currKbps60 := getPresetBitrateKbps(preset, isPortrait, true)
			dynKbps60 := scaleBitrateBySource(sourceMaxKbps, maxPresetKbps60, currKbps60)
			targets = append(targets, TranscodingTarget{Resolution: resStr, BitrateRate: formatBitrateKbps(dynKbps60), FPS: videoInfo.FPS60, FpsName: "60"})
		}
		targets = append(targets, TranscodingTarget{Resolution: resStr, BitrateRate: formatBitrateKbps(dynKbps30), FPS: videoInfo.FPS30, FpsName: "30"})
	}
	return targets
}

func getVideoInfo(input string) (info global.VideoInfo, err error) {
	cmd := exec.Command("ffprobe", "-i", input, "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", "-probesize", "5000000", "-analyzeduration", "5000000")
	out, err := utils.RunCmd(cmd)
	if err != nil {
		return info, err
	}
	err = json.Unmarshal(out.Bytes(), &info)
	return info, err
}

func parseFPS(fps string) float64 {
	if strings.Contains(fps, "/") {
		parts := strings.Split(fps, "/")
		if len(parts) == 2 {
			num, _ := strconv.ParseFloat(parts[0], 64)
			den, _ := strconv.ParseFloat(parts[1], 64)
			if den > 0 {
				return num / den
			}
		}
	}
	f, _ := strconv.ParseFloat(fps, 64)
	return f
}

func parseFfmpegClockToSeconds(clock string) float64 {
	parts := strings.Split(clock, ":")
	if len(parts) != 3 {
		return 0
	}
	h, _ := strconv.ParseFloat(parts[0], 64)
	m, _ := strconv.ParseFloat(parts[1], 64)
	s, _ := strconv.ParseFloat(parts[2], 64)
	return h*3600 + m*60 + s
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
	profile = strings.ToLower(profile)
	var profileHex string
	switch {
	case strings.Contains(profile, "baseline"):
		profileHex = "42"
	case strings.Contains(profile, "main"):
		profileHex = "4D"
	case strings.Contains(profile, "high"):
		profileHex = "64"
	default:
		return "", fmt.Errorf("unsupported h264 profile: %s", profile)
	}

	levelHex := fmt.Sprintf("%02X", level)
	return fmt.Sprintf("avc1.%s00%s", profileHex, levelHex), nil
}

func buildRateControlParams(rate string) (string, string, string) {
	if kbps, err := strconv.Atoi(strings.TrimSuffix(rate, "k")); err == nil && kbps > 0 {
		maxKbps := int(math.Round(float64(kbps) * 1.4))
		if maxKbps < kbps {
			maxKbps = kbps
		}
		return rate, fmt.Sprintf("%dk", maxKbps), fmt.Sprintf("%dk", maxKbps*2)
	}
	return rate, rate, "4000k"
}

func parseQualityInfo(quality string) (int, int, int, float64) {
	w, h, bw, fr := 1920, 1080, 3000000, 30.0
	parts := strings.Split(quality, "_")
	if len(parts) > 0 {
		res := strings.Split(parts[0], "x")
		if len(res) == 2 {
			w, _ = strconv.Atoi(res[0])
			h, _ = strconv.Atoi(res[1])
		}
	}
	if len(parts) > 1 {
		r, _ := strconv.Atoi(strings.TrimSuffix(parts[1], "k"))
		bw = r * 1000
	}
	if len(parts) > 2 {
		fr, _ = strconv.ParseFloat(parts[2], 64)
	}
	return w, h, bw, fr
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
		b = parseBitrateKbps(parts[1])
	}
	if len(parts) > 2 {
		f, _ = strconv.Atoi(parts[2])
	}
	return
}
