package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/ffmpeg"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

const (
	// Redis Stream 名称，用于转码任务队列
	transcodingQueueStream = "transcoding:queue"
	// Redis List 名称，等待派发的转码任务（Worker 空闲时从此队列取出放入 Stream）
	transcodingPendingQueue = "transcoding:pending"
	// Redis Hash 状态前缀，后接 resourceID
	transcodingStatusPrefix = "transcoding:status:"
	// Pub/Sub 频道：转码取消
	transcodingCancelChannel = "transcoding:cancel"
	// Pub/Sub 频道：Worker 完成任务通知服务端派发下一个
	transcodingDispatchChannel = "transcoding:dispatch"
	// 进度轮询间隔
	progressPollInterval = 3 * time.Second
)

// RemoteTranscoder 通过 Redis Streams + OSS 与远程 Worker 池通信。
type RemoteTranscoder struct {
	rdb *redis.Client
}

const (
	// 心跳 key 前缀（与 worker 端保持一致）
	workerHeartbeatPrefix = "transcoding:workers:"
)

// WorkerHeartbeat Worker 心跳快照
type WorkerHeartbeat struct {
	Healthy     bool      `json:"healthy"`
	StartedAt   time.Time `json:"startedAt"`
	Uptime      string    `json:"uptime"`
	Concurrency int       `json:"concurrency"`
	JobsActive  int32     `json:"jobsActive"`
	JobsTotal   int64     `json:"jobsTotal"`
	JobsFailed  int64     `json:"jobsFailed"`
	GroupID     string    `json:"groupID"`
	LastSeen    int64     `json:"lastSeen"`
}

// GetAliveWorkers 扫描 Redis 中所有 Worker 的心跳 key，返回在线 Worker 列表。
// 心跳由 Worker 的 heartbeatLoop 每 10s 写入一次，TTL=30s。
func GetAliveWorkers(ctx context.Context) ([]WorkerHeartbeat, error) {
	iter := global.Redis.RawClient().Scan(ctx, 0, workerHeartbeatPrefix+"*", 100).Iterator()
	var workers []WorkerHeartbeat
	for iter.Next(ctx) {
		val, err := global.Redis.RawClient().Get(ctx, iter.Val()).Result()
		if err != nil {
			continue // key 可能在 scan 和 get 之间过期
		}
		var hb WorkerHeartbeat
		if err := json.Unmarshal([]byte(val), &hb); err != nil {
			continue
		}
		workers = append(workers, hb)
	}
	return workers, iter.Err()
}

// HasAliveWorker 快速检查是否有至少一个 Worker 在线。
func HasAliveWorker(ctx context.Context) bool {
	// 用 Exists 检查第一个匹配的 key，比 Scan 快得多
	iter := global.Redis.RawClient().Scan(ctx, 0, workerHeartbeatPrefix+"*", 1).Iterator()
	if iter.Next(ctx) {
		return true
	}
	return false
}

func NewRemoteTranscoder() *RemoteTranscoder {
	return &RemoteTranscoder{
		rdb: global.Redis.RawClient(),
	}
}

// Enqueue 将转码任务入队：有空闲 Worker 时直接投递 Stream，否则存入 pending 队列等待派发。
func (t *RemoteTranscoder) Enqueue(ctx context.Context, info *dto.TranscodingInfo) error {
	// 1. 填充编码参数（与主服务配置一致）
	info.UseGpu = global.Config.Transcoding.UseGpu
	info.UseH265 = global.Config.Transcoding.UseH265
	info.UseAv1 = global.Config.Transcoding.UseAv1
	info.Generate1080p60 = global.Config.Transcoding.Generate1080p60

	job, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("transcoding job marshal: %w", err)
	}

	// 2. 检查 Worker 可用容量
	available, err := getAvailableCapacity(ctx)
	if err != nil {
		zap.L().Warn("Failed to check worker capacity, falling back to pending queue", zap.Error(err))
		available = 0
	}

	// 3. 有空闲 Worker → 直接投 Stream；全忙 → 存 pending 队列等 DispatchPending 派发
	if available > 0 {
		if err := t.rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: transcodingQueueStream,
			Values: map[string]interface{}{
				"job":        string(job),
				"resourceID": info.ResourceID,
				"videoID":    info.VideoID,
			},
		}).Err(); err != nil {
			return fmt.Errorf("xadd transcoding job: %w", err)
		}
		utils.InfoLog(fmt.Sprintf("【转码入队→Stream】VideoID=%d ResourceID=%d",
			info.VideoID, info.ResourceID), "transcoding")
	} else {
		// 存入 pending 队列（JSON 序列化的任务 + resourceID）
		pending := map[string]interface{}{
			"job":        string(job),
			"resourceID": info.ResourceID,
			"videoID":    info.VideoID,
		}
		data, _ := json.Marshal(pending)
		if err := t.rdb.RPush(ctx, transcodingPendingQueue, string(data)).Err(); err != nil {
			return fmt.Errorf("rpush pending queue: %w", err)
		}
		utils.InfoLog(fmt.Sprintf("【转码入队→Pending】VideoID=%d ResourceID=%d PendingLen=%d",
			info.VideoID, info.ResourceID, t.rdb.LLen(ctx, transcodingPendingQueue).Val()), "transcoding")
	}

	// 4. 异步轮询进度
	go t.pollProgress(ctx, info)

	return nil
}

// DispatchPending 从 pending 队列取出任务投递到 Stream（由 Worker 完成通知或定时任务触发）。
func (t *RemoteTranscoder) DispatchPending(ctx context.Context) {
	available, err := getAvailableCapacity(ctx)
	if err != nil || available <= 0 {
		return
	}

	for i := 0; i < available; i++ {
		// RPOP 从 pending 队列取出最早的任务
		result, err := t.rdb.LPop(ctx, transcodingPendingQueue).Result()
		if err == redis.Nil || result == "" {
			break // pending 队列空了
		}
		if err != nil {
			zap.L().Error("LPop pending queue failed", zap.Error(err))
			break
		}

		// 解析 pending 任务
		var pending struct {
			Job        string `json:"job"`
			ResourceID uint   `json:"resourceID"`
			VideoID    uint   `json:"videoID"`
		}
		if err := json.Unmarshal([]byte(result), &pending); err != nil {
			zap.L().Error("Unmarshal pending job failed", zap.Error(err))
			continue
		}

		// 投递到 Stream
		if err := t.rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: transcodingQueueStream,
			Values: map[string]interface{}{
				"job":        pending.Job,
				"resourceID": pending.ResourceID,
				"videoID":    pending.VideoID,
			},
		}).Err(); err != nil {
			zap.L().Error("XAdd from pending failed, pushing back",
				zap.Uint("resourceID", pending.ResourceID), zap.Error(err))
			// 投递失败，放回 pending 队列头部
			t.rdb.LPush(ctx, transcodingPendingQueue, result)
			break
		}

		utils.InfoLog(fmt.Sprintf("【Pending→Stream】VideoID=%d ResourceID=%d",
			pending.VideoID, pending.ResourceID), "transcoding")
	}
}

// ListenDispatch 监听 Worker 完成通知，触发 DispatchPending 派发下一个任务。
func (t *RemoteTranscoder) ListenDispatch(ctx context.Context) {
	pubsub := t.rdb.Subscribe(ctx, transcodingDispatchChannel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-ch:
			if !ok {
				return
			}
			t.DispatchPending(ctx)
		}
	}
}

// getAvailableCapacity 扫描所有 Worker 心跳，计算总剩余容量（Σ 最大并发 - 当前活跃任务数）。
func getAvailableCapacity(ctx context.Context) (int, error) {
	workers, err := GetAliveWorkers(ctx)
	if err != nil {
		return 0, fmt.Errorf("get alive workers: %w", err)
	}
	if len(workers) == 0 {
		return 0, nil
	}
	var total int
	for _, w := range workers {
		remaining := w.Concurrency - int(w.JobsActive)
		if remaining < 0 {
			remaining = 0
		}
		total += remaining
	}
	return total, nil
}

// pollProgress 在后台轮询 Redis Hash 中的转码进度，直至完成或取消。
func (t *RemoteTranscoder) pollProgress(ctx context.Context, info *dto.TranscodingInfo) {
	statusKey := fmt.Sprintf("%s%d", transcodingStatusPrefix, info.ResourceID)
	ticker := time.NewTicker(progressPollInterval)
	defer ticker.Stop()

	// 订阅取消信号
	pubsub := t.rdb.Subscribe(ctx, transcodingCancelChannel)
	defer pubsub.Close()
	cancelCh := pubsub.Channel()

	timeout := time.After(24 * time.Hour) // 最长等待24小时
	var offlineSince *time.Time

	for {
		select {
		case <-ctx.Done():
			return
		case <-timeout:
			utils.ErrorLog(fmt.Sprintf("【远程转码超时】ResourceID=%d，24h 上限到达，标记失败", info.ResourceID), "transcoding", "")
			GetTranscoder().completeTransaction(ctx, info, global.PROCESSING_FAIL)
			t.rdb.Del(ctx, statusKey)
			return
		case msg, ok := <-cancelCh:
			if ok && msg.Payload == fmt.Sprintf("%d", info.VideoID) {
				utils.InfoLog(fmt.Sprintf("【远程转码取消】VideoID=%d", info.VideoID), "transcoding")
				return
			}
		case <-ticker.C:
			// 查 Worker 心跳，离线时不读进度但不标记失败
			alive := HasAliveWorker(ctx)
			if !alive {
				if offlineSince == nil {
					now := time.Now()
					offlineSince = &now
					utils.InfoLog(fmt.Sprintf("【远程转码】Worker 离线，等待恢复... ResourceID=%d", info.ResourceID), "transcoding")
				}
				continue // Worker 离线期间不读进度
			}
			if offlineSince != nil {
				utils.InfoLog(fmt.Sprintf("【远程转码】Worker 已恢复在线（离线 %v），ResourceID=%d",
					time.Since(*offlineSince), info.ResourceID), "transcoding")
				offlineSince = nil
			}

			statusMap, err := t.rdb.HGetAll(ctx, statusKey).Result()
			if err != nil || len(statusMap) == 0 {
				continue
			}
			status := statusMap["status"]
			if status == "failed" {
				utils.InfoLog(fmt.Sprintf("【远程转码失败】ResourceID=%d", info.ResourceID), "transcoding")
				GetTranscoder().completeTransaction(ctx, info, global.PROCESSING_FAIL)
				t.rdb.Del(ctx, statusKey)
				return
			}
			if status == "completed" {
				utils.InfoLog(fmt.Sprintf("【远程转码完成】ResourceID=%d", info.ResourceID), "transcoding")
				t.handleRemoteCompletion(ctx, info, statusMap)
				t.rdb.Del(ctx, statusKey)
				return
			}
		}
	}
}

// handleRemoteCompletion Worker 完成后，从 Redis 读取索引数据并创建 DB 记录。
func (t *RemoteTranscoder) handleRemoteCompletion(ctx context.Context, info *dto.TranscodingInfo, statusMap map[string]string) {
	svc := GetTranscoder()

	// 1. 从 statusMap 中找出所有 idx_{qualityName} 键，解析索引数据
	var indexRecords []model.VideoIndexFile
	for key, val := range statusMap {
		if !strings.HasPrefix(key, "idx_") {
			continue
		}
		qualityName := strings.TrimPrefix(key, "idx_")

		var idx struct {
			VideoInitRange  string `json:"vInit"`
			VideoIndexRange string `json:"vIndex"`
			AudioInitRange  string `json:"aInit"`
			AudioIndexRange string `json:"aIndex"`
			VideoCodec      string `json:"codec"`
		}
		if err := json.Unmarshal([]byte(val), &idx); err != nil {
			utils.ErrorLog("【远程转码】解析索引数据失败", "transcoding",
				fmt.Sprintf("quality=%s err=%v", qualityName, err))
			continue
		}

		width, height, bandwidth, _ := ffmpeg.ParseQualityInfo(qualityName)
		frameRate := ffmpeg.ParseFPS(qualityName[strings.LastIndex(qualityName, "_")+1:])

		videoCodec := idx.VideoCodec
		if videoCodec == "" {
			videoCodec = ffmpeg.DefaultVideoCodec
		}

		audioCodec := ffmpeg.DefaultAudioCodec
		audioBitrate := info.AudioBitRate
		audioSampleRate := info.AudioSampleRate

		// 多音轨：取主音轨的语言（第一个 AudioStream 的语言）
		// 单音轨回退：Worker 在无 AudioStreams 时始终上传 audio.m4s
		audioLanguage := "und"
		if len(info.AudioStreams) > 0 {
			audioLanguage = info.AudioStreams[0].Language
		}
		audioFile := ffmpeg.AudioFileNameForTrack(audioLanguage)

		indexRecords = append(indexRecords, model.VideoIndexFile{
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
			VideoInitRange:  idx.VideoInitRange,
			VideoIndexRange: idx.VideoIndexRange,
			AudioFile:       audioFile,
			AudioBandwidth:  audioBitrate,
			AudioCodec:      audioCodec,
			AudioSampleRate: audioSampleRate,
			AudioInitRange:  idx.AudioInitRange,
			AudioIndexRange: idx.AudioIndexRange,
		})
	}

	if len(indexRecords) == 0 {
		utils.ErrorLog("【远程转码】没有有效的索引数据，标记失败", "transcoding",
			fmt.Sprintf("ResourceID=%d", info.ResourceID))
		if err := svc.completeTransaction(ctx, info, global.PROCESSING_FAIL); err != nil {
			utils.ErrorLog("【远程转码】标记失败时事务出错", "transcoding", err.Error())
		}
		return
	}

	// 2. 事务内批量创建 VideoIndexFile 记录（防止部分写入后失败导致脏数据）
	if err := svc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, rec := range indexRecords {
			if err := tx.Create(&rec).Error; err != nil {
				return fmt.Errorf("create index record for %s: %w", rec.Quality, err)
			}
		}
		return nil
	}); err != nil {
		utils.ErrorLog("【远程转码】批量创建索引记录失败", "transcoding", err.Error())
		if err := svc.completeTransaction(ctx, info, global.PROCESSING_FAIL); err != nil {
			utils.ErrorLog("【远程转码】标记失败时事务出错", "transcoding", err.Error())
		}
		return
	}

	// 3. 多音轨记录（远程模式下音频文件在 OSS 上，暂不创建 AudioTrack 记录）
	// TODO: 远程模式需要 Worker 同时上报各音轨的 init/index range

	// 4. 完成事务（更新资源状态）
	utils.InfoLog(fmt.Sprintf("【远程转码】索引创建完成，共 %d 个画质", len(indexRecords)), "transcoding")
	if err := svc.completeTransaction(ctx, info, global.WAITING_REVIEW); err != nil {
		utils.ErrorLog("【远程转码】完成事务失败", "transcoding", err.Error())
	}

	// 5. 清理本地源文件和目录
	// 远程模式下 Worker 从 OSS 拉取源文件进行转码，转码产物也直接写入 OSS，
	// 本地副本（upload.mp4 等）已无用处，删除以释放 VPS 磁盘空间。
	if info.OutputDir != "" {
		if err := os.RemoveAll(info.OutputDir); err != nil {
			utils.ErrorLog("【远程转码】清理本地源文件目录失败", "transcoding",
				fmt.Sprintf("dir=%s err=%v", info.OutputDir, err))
		} else {
			utils.InfoLog(fmt.Sprintf("【远程转码】已清理本地源文件目录: %s", info.OutputDir), "transcoding")
		}
	}
}

// GetProgress 从 Redis Hash 读取转码进度。
func (t *RemoteTranscoder) GetProgress(ctx context.Context, resourceID uint) (*TranscodingProgress, error) {
	statusKey := fmt.Sprintf("%s%d", transcodingStatusPrefix, resourceID)
	statusMap, err := t.rdb.HGetAll(ctx, statusKey).Result()
	if err != nil {
		return nil, fmt.Errorf("hgetall %s: %w", statusKey, err)
	}
	if len(statusMap) == 0 {
		return nil, nil
	}

	// 写入的 key 由 Worker 定义，这里做映射
	tp := &TranscodingProgress{}
	if s, ok := statusMap["status"]; ok {
		tp.Status = s
	}
	if p, ok := statusMap["progress"]; ok {
		fmt.Sscanf(p, "%f", &tp.OverallProgress)
	}

	// 各 quality 进度
	if qMap, ok := statusMap["qualities"]; ok {
		_ = json.Unmarshal([]byte(qMap), &tp.Qualities)
	}

	// 上传进度
	if uMap, ok := statusMap["upload"]; ok {
		_ = json.Unmarshal([]byte(uMap), &tp.Upload)
	}

	return tp, nil
}

// Cancel 发布取消信号到 Redis Pub/Sub，Worker 收到后终止对应任务。
func (t *RemoteTranscoder) Cancel(ctx context.Context, videoID uint) error {
	if err := t.rdb.Publish(ctx, transcodingCancelChannel, fmt.Sprintf("%d", videoID)).Err(); err != nil {
		return fmt.Errorf("publish cancel %d: %w", videoID, err)
	}
	return nil
}
