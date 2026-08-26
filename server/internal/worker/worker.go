package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"interastral-peace.com/alnitak/internal/config"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/ffmpeg"
	"interastral-peace.com/alnitak/pkg/oss"
)

// =============================================================================
// 常量
// =============================================================================

// Redis Key 常量（与服务端 remote_transcoder.go 保持一致）
const (
	transcodingQueueStream     = "transcoding:queue"
	transcoderGroup            = "transcoder"
	transcodingStatusPrefix    = "transcoding:status:"
	transcodingCancelChannel   = "transcoding:cancel"
	transcodingCompleteChannel = "transcoding:complete"
	transcodingDispatchChannel = "transcoding:dispatch"
	deadLetterStream           = "transcoding:deadletter"
)

const (
	// Redis Key 前缀：Worker 心跳（主服务通过此 key 监测 Worker 存活）
	heartbeatKeyPrefix = "transcoding:workers:"
	// 心跳写入间隔
	heartbeatInterval = 10 * time.Second
	// 心跳 TTL（比间隔大，允许丢 2 次心跳仍视为存活）
	heartbeatTTL = 30 * time.Second

	// 超过此空闲时间的 pending 消息视为被遗弃（前一个 Worker crash）
	maxPendingIdle = 2 * time.Minute
	// 单个消息最大重试次数（含重新投递），超限进入死信
	maxRetryCount = 3
	// XReadGroup 阻塞时长
	readBlockTime = 5 * time.Second
	// XReadGroup 每次批量拉取上限
	readBatchSize = 10
)

// =============================================================================
// Redis 重试工具
// =============================================================================

// redisMaxRetries 关键 Redis 写操作的最大重试次数
const redisMaxRetries = 3

// retryHSet 带重试的 HSet，防止网络丢包导致进度/状态丢失
func (w *Worker) retryHSet(ctx context.Context, key string, values ...interface{}) error {
	var lastErr error
	for i := 0; i < redisMaxRetries; i++ {
		if err := w.rdb.HSet(ctx, key, values...).Err(); err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
			continue
		}
		return nil
	}
	zap.L().Error("HSet failed after retries",
		zap.String("key", key),
		zap.Int("retries", redisMaxRetries),
		zap.Error(lastErr))
	return lastErr
}

// retryXAdd 带重试的 XAdd
func (w *Worker) retryXAdd(ctx context.Context, args *redis.XAddArgs) error {
	var lastErr error
	for i := 0; i < redisMaxRetries; i++ {
		if err := w.rdb.XAdd(ctx, args).Err(); err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
			continue
		}
		return nil
	}
	zap.L().Error("XAdd failed after retries",
		zap.String("stream", args.Stream),
		zap.Int("retries", redisMaxRetries),
		zap.Error(lastErr))
	return lastErr
}

// retryPublish 带重试的 Publish
func (w *Worker) retryPublish(ctx context.Context, channel string, message interface{}) error {
	var lastErr error
	for i := 0; i < redisMaxRetries; i++ {
		if err := w.rdb.Publish(ctx, channel, message).Err(); err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
			continue
		}
		return nil
	}
	zap.L().Error("Publish failed after retries",
		zap.String("channel", channel),
		zap.Int("retries", redisMaxRetries),
		zap.Error(lastErr))
	return lastErr
}

// =============================================================================
// Worker 主结构
// =============================================================================

// Worker 远程转码 Worker，从 Redis Stream 消费任务，
// 执行 ffmpeg 编码后将结果上传到 OSS。
type Worker struct {
	rdb           *redis.Client
	storage       oss.Storage
	backupStorage oss.Storage // 备用 OSS（多源容灾），可为 nil
	cfg           *config.Config
	groupID       string
	concurrency   int
	cancelSub     *redis.PubSub
	sem           chan struct{}
	encodingSem   chan struct{} // 画质编码并发限制
	nvencSem      chan struct{} // 全局 NVENC 硬件上限 (cap=8, GeForce 驱动限制)
	cpuSem        chan struct{} // CPU fallback 保护 (cap=2, 防线程打满)
	workDir       string       // 临时文件工作目录，默认 os.TempDir()

	mu           sync.RWMutex
	activeVideos map[uint]context.CancelFunc

	// 健康与统计
	startedAt  time.Time
	jobsActive atomic.Int32
	jobsTotal  atomic.Int64
	jobsFailed atomic.Int64
	healthy    atomic.Bool
}

// NewWorker 创建 Worker 实例，建立 Redis/OSS 连接。
func NewWorker(cfg *config.Config, workerID string, concurrency int) (*Worker, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
	})

	var storage oss.Storage
	if cfg.Storage.OssType != "local" {
		storage = oss.InitStorage(cfg.Storage)
	}
	backupStorage := oss.InitBackupStorage(cfg.Storage)

	workDir := cfg.Transcoding.WorkDir
	if workDir == "" {
		workDir = os.TempDir()
	} else if !filepath.IsAbs(workDir) {
		// 相对路径基于工作目录（nssm AppDirectory）解析
		if absDir, err := filepath.Abs(workDir); err == nil {
			workDir = absDir
		}
	}

	encConcurrency := cfg.Transcoding.EncodingConcurrency
	if encConcurrency <= 0 {
		encConcurrency = 0 // 0 = unlimited
	}

	w := &Worker{
		rdb:           rdb,
		storage:       storage,
		backupStorage: backupStorage,
		cfg:           cfg,
		groupID:       fmt.Sprintf("transcoder-worker-%s", workerID),
		concurrency:   concurrency,
		workDir:       workDir,
		sem:           make(chan struct{}, concurrency),
		activeVideos:  make(map[uint]context.CancelFunc),
		startedAt:     time.Now(),
	}

	if encConcurrency > 0 {
		w.encodingSem = make(chan struct{}, encConcurrency)
	}

	// 全局 NVENC 上限（GeForce 驱动限制 8 路）和 CPU fallback 保护
	w.nvencSem = make(chan struct{}, 8)
	w.cpuSem = make(chan struct{}, 2)

	w.healthy.Store(true)
	return w, nil
}

// =============================================================================
// 健康检查（供 HTTP 管理接口调用）
// =============================================================================

// HealthStatus Worker 运行时状态快照
type HealthStatus struct {
	Healthy     bool      `json:"healthy"`
	StartedAt   time.Time `json:"startedAt"`
	Uptime      string    `json:"uptime"`
	Concurrency int       `json:"concurrency"`
	JobsActive  int32     `json:"jobsActive"`
	JobsTotal   int64     `json:"jobsTotal"`
	JobsFailed  int64     `json:"jobsFailed"`
	GroupID     string    `json:"groupID"`
}

// Health 返回当前状态快照。
func (w *Worker) Health() HealthStatus {
	return HealthStatus{
		Healthy:     w.healthy.Load(),
		StartedAt:   w.startedAt,
		Uptime:      time.Since(w.startedAt).Truncate(time.Second).String(),
		Concurrency: w.concurrency,
		JobsActive:  w.jobsActive.Load(),
		JobsTotal:   w.jobsTotal.Load(),
		JobsFailed:  w.jobsFailed.Load(),
		GroupID:     w.groupID,
	}
}

// Ready 检查 Worker 是否就绪（Redis 连通 + 正常运行）。
func (w *Worker) Ready() bool {
	if !w.healthy.Load() {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return w.rdb.Ping(ctx).Err() == nil
}

// =============================================================================
// 心跳 — 向 Redis 写入 Worker 存活状态
// =============================================================================

// heartbeatLoop 每隔 heartbeatInterval 向 Redis 写入心跳 key，TTL = heartbeatTTL。
// 主服务通过检查该 key 是否存在来判断 Worker 是否存活。
func (w *Worker) heartbeatLoop(ctx context.Context) {
	key := heartbeatKeyPrefix + w.groupID
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	// 立即写一次
	w.writeHeartbeat(ctx, key)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.writeHeartbeat(ctx, key)
		}
	}
}

func (w *Worker) writeHeartbeat(ctx context.Context, key string) {
	h := w.Health()
	data, err := json.Marshal(map[string]interface{}{
		"healthy":     h.Healthy,
		"startedAt":   h.StartedAt,
		"uptime":      h.Uptime,
		"concurrency": h.Concurrency,
		"jobsActive":  h.JobsActive,
		"jobsTotal":   h.JobsTotal,
		"jobsFailed":  h.JobsFailed,
		"groupID":     h.GroupID,
		"lastSeen":    time.Now().Unix(),
	})
	if err != nil {
		zap.L().Warn("Heartbeat marshal failed", zap.Error(err))
		return
	}

	if err := w.rdb.Set(ctx, key, string(data), heartbeatTTL).Err(); err != nil {
		zap.L().Warn("Heartbeat write failed", zap.Error(err))
	}
}

// =============================================================================
// Run — 主消费循环
// =============================================================================

// Run 进入消费循环，阻塞至 ctx 取消。
func (w *Worker) Run(ctx context.Context) error {
	// 确保 Redis Stream 消费者组存在（幂等）
	if err := w.rdb.XGroupCreateMkStream(ctx, transcodingQueueStream, transcoderGroup, "0").Err(); err != nil {
		if !strings.Contains(err.Error(), "BUSYGROUP") {
			return fmt.Errorf("create consumer group: %w", err)
		}
	}

	// 启动时恢复 crash 遗留的 pending 消息
	if err := w.recoverPending(ctx); err != nil {
		zap.L().Warn("Pending recovery partial", zap.Error(err))
		// 不阻断启动，部分消息可能被其他消费者持有
	}

	// 订阅取消频道
	w.cancelSub = w.rdb.Subscribe(ctx, transcodingCancelChannel)
	defer w.cancelSub.Close()
	cancelCh := w.cancelSub.Channel()

	var wg sync.WaitGroup
	defer wg.Wait() // 退出时等待所有 goroutine 排空

	zap.L().Info("Worker started",
		zap.String("groupID", w.groupID),
		zap.Int("concurrency", w.concurrency),
	)

	// 启动心跳（goroutine 随 ctx 退出）
	heartbeatCtx, heartbeatCancel := context.WithCancel(ctx)
	defer heartbeatCancel()
	go w.heartbeatLoop(heartbeatCtx)

	for {
		select {
		case <-ctx.Done():
			zap.L().Info("Worker shutting down, waiting for in-flight jobs...")
			wg.Wait()
			return ctx.Err()

		case msg := <-cancelCh:
			w.handleCancel(ctx, msg.Payload)

		default:
			// 周期性 Redis 连通性检测（每轮循环顺便检查）
			if err := w.rdb.Ping(ctx).Err(); err != nil {
				w.healthy.Store(false)
				zap.L().Error("Redis ping failed, retrying in 5s", zap.Error(err))
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(5 * time.Second):
				}
				continue
			}
			w.healthy.Store(true)

			streams, err := w.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    transcoderGroup,
				Consumer: w.groupID,
				Streams:  []string{transcodingQueueStream, ">"},
				Count:    readBatchSize,
				Block:    readBlockTime,
			}).Result()
			if err != nil {
				if err == redis.Nil {
					continue
				}
				zap.L().Error("XReadGroup error", zap.Error(err))
				continue
			}

			for _, stream := range streams {
				for _, msg := range stream.Messages {
					msgID := msg.ID
					jobJSON, ok := msg.Values["job"].(string)
					if !ok {
						w.rdb.XAck(ctx, transcodingQueueStream, transcoderGroup, msgID)
						w.rdb.XDel(ctx, transcodingQueueStream, msgID)
						continue
					}

					w.sem <- struct{}{}
					wg.Add(1)
					w.jobsActive.Add(1)
					w.jobsTotal.Add(1)
					go func(id string, rawJob string) {
						defer func() {
							if r := recover(); r != nil {
								w.jobsFailed.Add(1)
								zap.L().Error("Job goroutine panic", zap.String("msgID", id), zap.Any("panic", r))
							}
						}()
						defer func() { <-w.sem }()
						defer wg.Done()
						defer w.jobsActive.Add(-1)
						if err := w.processJob(ctx, rawJob, id); err != nil {
							w.jobsFailed.Add(1)
							zap.L().Error("Job failed", zap.String("msgID", id), zap.Error(err))
						}
					}(msgID, jobJSON)
				}
			}
		}
	}
}

// =============================================================================
// Pending 消息恢复
// =============================================================================

// recoverPending 启动时处理遗留的 pending 消息：
//  1. 自己（相同 groupID）之前 crash 留下的 → 立即认领，不等 idle
//  2. 其他 Worker 遗弃且 idle 超时的 → 按 maxPendingIdle 认领
//  3. 超过重试上限的 → 移入死信
func (w *Worker) recoverPending(ctx context.Context) error {
	pendingItems, err := w.rdb.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream: transcodingQueueStream,
		Group:  transcoderGroup,
		Start:  "-",
		End:    "+",
		Count:  100,
	}).Result()
	if err != nil {
		return fmt.Errorf("xpendingext: %w", err)
	}
	if len(pendingItems) == 0 {
		return nil
	}

	zap.L().Info("Found pending messages for recovery",
		zap.Int("count", len(pendingItems)))

	var ownIDs []string      // 自己 crash 留下的（立即认领）
	var abandonedIDs []string // 其他 Worker 遗弃的（等 idle 超时）
	for _, p := range pendingItems {
		if p.RetryCount >= maxRetryCount {
			w.sendToDeadLetter(ctx, p.ID, fmt.Sprintf("retry count %d exceeded max (consumer=%s, idle=%s)", p.RetryCount, p.Consumer, p.Idle))
			w.rdb.XAck(ctx, transcodingQueueStream, transcoderGroup, p.ID)
			w.rdb.XDel(ctx, transcodingQueueStream, p.ID)
			zap.L().Warn("Moved to dead letter",
				zap.String("msgID", p.ID),
				zap.Int64("retryCount", p.RetryCount),
				zap.String("consumer", p.Consumer))
			continue
		}
		if p.Consumer == w.groupID {
			ownIDs = append(ownIDs, p.ID)
		} else if p.Idle >= maxPendingIdle {
			abandonedIDs = append(abandonedIDs, p.ID)
		}
	}

	// 认领并处理一批消息，返回 claimed 的 msg 列表
	claimAndRun := func(ids []string, minIdle time.Duration) {
		if len(ids) == 0 {
			return
		}
		claimed, err := w.rdb.XClaim(ctx, &redis.XClaimArgs{
			Stream:   transcodingQueueStream,
			Group:    transcoderGroup,
			Consumer: w.groupID,
			MinIdle:  minIdle,
			Messages: ids,
		}).Result()
		if err != nil {
			zap.L().Error("XClaim failed", zap.Error(err), zap.Int("count", len(ids)))
			return
		}
		for _, msg := range claimed {
			jobJSON, ok := msg.Values["job"].(string)
			if !ok {
				w.rdb.XAck(ctx, transcodingQueueStream, transcoderGroup, msg.ID)
				w.rdb.XDel(ctx, transcodingQueueStream, msg.ID)
				continue
			}
			w.sem <- struct{}{}
			w.jobsActive.Add(1)
			w.jobsTotal.Add(1)
			go func(id, raw string) {
				defer func() {
					if r := recover(); r != nil {
						w.jobsFailed.Add(1)
						zap.L().Error("Recovered job goroutine panic", zap.String("msgID", id), zap.Any("panic", r))
					}
				}()
				defer func() { <-w.sem }()
				defer w.jobsActive.Add(-1)
				if err := w.processJob(ctx, raw, id); err != nil {
					w.jobsFailed.Add(1)
					zap.L().Error("Recovered job failed", zap.String("msgID", id), zap.Error(err))
				}
			}(msg.ID, jobJSON)
		}
	}

	// 先认领自己的（立即，MinIdle=0），再认领遗弃的
	claimAndRun(ownIDs, 0)
	claimAndRun(abandonedIDs, maxPendingIdle)

	msg := ""
	if n := len(ownIDs); n > 0 {
		msg += fmt.Sprintf("own=%d ", n)
	}
	if n := len(abandonedIDs); n > 0 {
		msg += fmt.Sprintf("abandoned=%d", n)
	}
	if msg != "" {
		zap.L().Info("Recovered pending messages", zap.String("summary", msg))
	}

	return nil
}

// sendToDeadLetter 将无法处理的消息写入死信队列以便人工排查。
func (w *Worker) sendToDeadLetter(ctx context.Context, msgID, reason string) {
	w.retryXAdd(ctx, &redis.XAddArgs{
		Stream: deadLetterStream,
		Values: map[string]interface{}{
			"originalMsgID": msgID,
			"reason":        reason,
			"movedAt":       time.Now().Unix(),
		},
	})
}

// =============================================================================
// 任务处理
// =============================================================================

// processJob 处理单个转码任务。
func (w *Worker) processJob(ctx context.Context, rawJob, msgID string) (err error) {
	var info dto.TranscodingInfo
	if err := json.Unmarshal([]byte(rawJob), &info); err != nil {
		return fmt.Errorf("unmarshal job: %w", err)
	}

	jobCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 无论成功还是失败都 ACK + XDEL 消息：
	//   - XAck 防止 Worker 重启后重新消费已处理的任务
	//   - XDEL 防止 Stream 堆积导致 maxDepth 上限被占满
	defer func() {
		w.rdb.XAck(ctx, transcodingQueueStream, transcoderGroup, msgID)
		w.rdb.XDel(ctx, transcodingQueueStream, msgID)
	}()

	// 注册当前任务以便取消
	w.mu.Lock()
	w.activeVideos[info.VideoID] = cancel
	w.mu.Unlock()
	defer func() {
		w.mu.Lock()
		delete(w.activeVideos, info.VideoID)
		w.mu.Unlock()
	}()

	// 清理本 resourceID 的历史残留文件 + 退出时清理本次临时文件
	cleanPattern := filepath.Join(w.workDir, fmt.Sprintf("alnitak_*_%d*", info.ResourceID))
	defer func() {
		if leftovers, err := filepath.Glob(cleanPattern); err == nil {
			for _, f := range leftovers {
				if removeErr := os.Remove(f); removeErr == nil {
					zap.L().Debug("Cleaned temp file", zap.String("path", f))
				}
			}
		}
	}()
	if oldFiles, err := filepath.Glob(cleanPattern); err == nil {
		for _, f := range oldFiles {
			os.Remove(f)
		}
	}

	// 1. 下载输入文件
	statusKey := fmt.Sprintf("%s%d", transcodingStatusPrefix, info.ResourceID)
	w.writeStatus(statusKey, "processing", 0, "")

	inputLocal := filepath.Join(w.workDir, fmt.Sprintf("alnitak_input_%d%s", info.ResourceID, info.Suffix))
	zap.L().Info("Downloading input file", zap.String("objectKey", fmt.Sprintf("video/%s/upload%s", info.DirName, info.Suffix)))
	if err := w.downloadInput(&info, inputLocal); err != nil {
		w.writeStatus(statusKey, "failed", 0, err.Error())
		return fmt.Errorf("download input: %w", err)
	}
	zap.L().Info("Input file downloaded", zap.String("local", inputLocal))
	defer os.Remove(inputLocal)
	info.InputFile = inputLocal

	// 1.5 如果元数据不全（远程重转码场景），ffprobe 探测补充
	if info.Width == 0 || info.Height == 0 {
		if err := w.probeInput(&info); err != nil {
			w.writeStatus(statusKey, "failed", 0, err.Error())
			return fmt.Errorf("probe input: %w", err)
		}
	}

	// 2. 计算转码目标
	targets := w.computeTargets(&info)
	totalQualities := len(targets)
	zap.L().Info("Encoding targets",
		zap.Int("total", totalQualities),
		zap.Any("qualities", func() []string {
			var names []string
			for _, t := range targets {
				names = append(names, t.Resolution+"_"+t.BitrateRate+"_"+t.FpsName)
			}
			return names
		}()),
	)
	var completed atomic.Int32

	// 2a. 清理旧的画质进度字段（避免多次重转码残留不同码率的旧数据）
	progressKey := fmt.Sprintf("%s%d", transcodingStatusPrefix, info.ResourceID)
	currentQualities := make(map[string]bool, len(targets))
	for _, target := range targets {
		qn := target.Resolution + "_" + target.BitrateRate + "_" + target.FpsName
		currentQualities[qn] = true
	}
	if oldFields, err := w.rdb.HKeys(ctx, progressKey).Result(); err == nil {
		for _, field := range oldFields {
			if strings.HasPrefix(field, "progress_") || strings.HasPrefix(field, "status_") {
				qualName := strings.TrimPrefix(field, "progress_")
				qualName = strings.TrimPrefix(qualName, "status_")
				if !currentQualities[qualName] {
					w.rdb.HDel(ctx, progressKey, field)
				}
			}
		}
	}

	// 2b. 初始化所有画质进度为 0%（排队中的也能看见）
	for _, target := range targets {
		qualityName := target.Resolution + "_" + target.BitrateRate + "_" + target.FpsName
		progressField := fmt.Sprintf("progress_%s", qualityName)
		// 只在尚未开始时写入（避免覆盖已完成的进度）
		exists, _ := w.rdb.HExists(ctx, progressKey, progressField).Result()
		if !exists {
			w.rdb.HSet(ctx, progressKey, progressField, "0.00")
			w.rdb.HSet(ctx, progressKey, fmt.Sprintf("status_%s", qualityName), "waiting")
		}
	}

	// 3. 并发编码视频 + 音轨
	var encWg sync.WaitGroup
	errCh := make(chan error, totalQualities)

	// 3a. 视频编码（受 encodingSem 限制并发数）
	for _, target := range targets {
		qualityName := target.Resolution + "_" + target.BitrateRate + "_" + target.FpsName
		videoOutput := filepath.Join(w.workDir, fmt.Sprintf("alnitak_v_%d_%s.m4s", info.ResourceID, qualityName))

		// 检查该画质是否已成功完成（部分重转码场景）
		existingStatus, _ := w.rdb.HGet(ctx, progressKey, fmt.Sprintf("status_%s", qualityName)).Result()
		if existingStatus == "success" {
			completed.Add(1)
			zap.L().Info("Skipping already completed quality", zap.String("quality", qualityName))
			continue
		}

		zap.L().Info("Starting video encode", zap.String("quality", qualityName))
		encWg.Add(1)
		go func(t ffmpeg.TranscodingTarget, outPath, qName string) {
			defer encWg.Done()
			if w.encodingSem != nil {
				w.encodingSem <- struct{}{}
				defer func() { <-w.encodingSem }()
			}
			// 标记开始编码（从 waiting → processing）
			w.rdb.HSet(ctx, statusKey, fmt.Sprintf("status_%s", qName), "processing")
			if err := w.encodeQuality(jobCtx, &info, t, outPath); err != nil {
				zap.L().Error("Video encode failed", zap.String("quality", qName), zap.Error(err))
				w.rdb.HSet(ctx, statusKey, fmt.Sprintf("status_%s", qName), "fail")
				errCh <- fmt.Errorf("%s: %w", qName, err)
				return
			}
			w.rdb.HSet(ctx, statusKey, fmt.Sprintf("status_%s", qName), "success")
			completed.Add(1)
			zap.L().Info("Video encode done", zap.String("quality", qName))
		}(target, videoOutput, qualityName)
	}

	// 3b. 音频编码（与视频并发）
	var audioFiles []audioFileEntry
	audioErrCh := make(chan error, 1)
	zap.L().Info("Starting audio encode")
	go func() {
		files, err := w.encodeAudio(jobCtx, &info)
		if err == nil {
			audioFiles = files
		}
		audioErrCh <- err
	}()

	// 等待视频完成
	zap.L().Info("Waiting for video encodes to finish")
	encWg.Wait()
	close(errCh)

	// 等待音频完成
	if audioErr := <-audioErrCh; audioErr != nil {
		zap.L().Error("Audio encode failed", zap.Error(audioErr))
		w.writeStatus(statusKey, "failed", float64(completed.Load())/float64(totalQualities)*100, audioErr.Error())
		// 主音轨失败视为整体失败
		return fmt.Errorf("audio encode: %w", audioErr)
	}
	zap.L().Info("Audio encode done")

	// 收集视频编码错误
	var errs []string
	for err := range errCh {
		errs = append(errs, err.Error())
	}
	if len(errs) > 0 {
		w.writeStatus(statusKey, "partial_fail", float64(completed.Load())/float64(totalQualities)*100,
			fmt.Sprintf("errors: %s", strings.Join(errs, "; ")))
	} else {
		w.writeStatus(statusKey, "encoding_done", 100, "")
	}

	// 3c. 收集 MP4 init/index 范围（上传前文件仍在本地）
	for _, target := range targets {
		qName := target.Resolution + "_" + target.BitrateRate + "_" + target.FpsName
		videoFile := filepath.Join(w.workDir, fmt.Sprintf("alnitak_v_%d_%s.m4s", info.ResourceID, qName))

		var audioFile string
		if len(audioFiles) > 0 {
			audioFile = audioFiles[0].LocalPath
		}

		idx := w.collectQualityIndex(videoFile, audioFile, &info, qName)
		if idx != nil {
			data, _ := json.Marshal(idx)
			w.retryHSet(ctx, statusKey, fmt.Sprintf("idx_%s", qName), string(data))
		}
	}

	// 4. 上传转码产物到 OSS
	if w.storage != nil {
		w.writeStatus(statusKey, "uploading", 100, "")
		// 写入上传进度信息供主服务查询
		uploadInfo := map[string]interface{}{
			"ossType":  w.cfg.Storage.OssType,
			"progress": 0,
			"status":   "uploading",
		}
		if data, err := json.Marshal(uploadInfo); err == nil {
			w.retryHSet(ctx, statusKey, "upload", string(data))
		}
		if err := w.uploadOutputs(ctx, statusKey, &info, targets, audioFiles); err != nil {
			w.writeStatus(statusKey, "failed", 100, err.Error())
			return fmt.Errorf("upload outputs: %w", err)
		}
		// 上传完成后更新进度
		uploadInfo["progress"] = 100
		uploadInfo["status"] = "success"
		if data, err := json.Marshal(uploadInfo); err == nil {
			w.retryHSet(ctx, statusKey, "upload", string(data))
		}
	}

	// 5. 标记完成（关键写入，重试防丢包）
	w.retryHSet(ctx, statusKey, "status", "completed", "progress", "100")
	w.retryPublish(ctx, transcodingCompleteChannel, fmt.Sprintf("%d", info.ResourceID))
	// 通知服务端：本 Worker 空闲了，可以派发 pending 队列中的下一个任务
	w.retryPublish(ctx, transcodingDispatchChannel, "dispatch")

	zap.L().Info("Job completed",
		zap.Uint("resourceID", info.ResourceID),
		zap.Uint("videoID", info.VideoID),
		zap.Int("qualities", totalQualities),
	)
	return nil
}

// =============================================================================
// 进度上报
// =============================================================================

// writeStatus 写入转码进度到 Redis Hash（带重试防丢包）
func (w *Worker) writeStatus(key, status string, progress float64, errMsg string) {
	fields := map[string]interface{}{
		"status":   status,
		"progress": fmt.Sprintf("%.1f", progress),
		"updated":  time.Now().Unix(),
	}
	if errMsg != "" {
		fields["error"] = errMsg
	}
	w.retryHSet(context.Background(), key, fields)
}

// =============================================================================
// OSS 操作
// =============================================================================

// downloadInput 从 OSS 下载输入文件到本地临时路径。
// 主 OSS 失败时自动尝试备用 OSS（多源容灾）。
func (w *Worker) downloadInput(info *dto.TranscodingInfo, localPath string) error {
	if w.storage == nil {
		return fmt.Errorf("no OSS storage configured")
	}
	objectKey := fmt.Sprintf("video/%s/upload%s", info.DirName, info.Suffix)

	err := w.storage.GetObjectToFile(objectKey, localPath)
	if err != nil && w.backupStorage != nil {
		zap.L().Warn("Primary OSS download failed, trying backup",
			zap.String("objectKey", objectKey),
			zap.Error(err))
		return w.backupStorage.GetObjectToFile(objectKey, localPath)
	}
	return err
}

// probeInput 对已下载的源文件执行 ffprobe，补充 Width/Height 等元数据。
// 远程重转码场景下 main server 可能未提供这些字段，由 Worker 自行探测。
func (w *Worker) probeInput(info *dto.TranscodingInfo) error {
	result, err := ffmpeg.RunFFprobe(info.InputFile)
	if err != nil {
		return fmt.Errorf("ffprobe %s: %w", info.InputFile, err)
	}

	for _, s := range result.Streams {
		if s.CodecType == "video" {
			if s.Width > 0 {
				info.Width = s.Width
			}
			if s.Height > 0 {
				info.Height = s.Height
			}
			if s.CodecName != "" {
				info.CodecName = s.CodecName
			}
			if s.AvgFrameRate != "" {
				info.FPS = s.AvgFrameRate
				generate1080p60 := info.Generate1080p60 || w.cfg.Transcoding.Generate1080p60
				info.FPS30, info.FPS60 = ffmpeg.GetFpsInfo(info.FPS, generate1080p60)
			}
			if br, err := strconv.Atoi(s.BitRate); err == nil && br > 0 {
				info.VideoBitRate = br
			}
		}
	}

	if info.Duration <= 0 {
		if dur, err := strconv.ParseFloat(result.Format.Duration, 64); err == nil && dur > 0 {
			info.Duration = dur
		}
	}

	if info.Width == 0 || info.Height == 0 {
		return fmt.Errorf("could not determine video dimensions from: %s", info.InputFile)
	}

	zap.L().Info("Probed input file",
		zap.Int("width", info.Width),
		zap.Int("height", info.Height),
		zap.String("fps", info.FPS),
		zap.Int("bitrate", info.VideoBitRate),
	)
	return nil
}

// uploadOutputs 将转码产物（视频+音频）上传到 OSS（主+备）。
// audioFiles 由 encodeAudio 返回，含临时路径与 OSS key。
func (w *Worker) uploadOutputs(ctx context.Context, statusKey string, info *dto.TranscodingInfo, targets []ffmpeg.TranscodingTarget, audioFiles []audioFileEntry) error {
	type uploadFile struct {
		localPath string
		ossKey    string
	}
	var files []uploadFile

	// 收集视频文件
	for _, target := range targets {
		qName := target.Resolution + "_" + target.BitrateRate + "_" + target.FpsName
		videoFile := filepath.Join(w.workDir, fmt.Sprintf("alnitak_v_%d_%s.m4s", info.ResourceID, qName))
		if _, err := os.Stat(videoFile); os.IsNotExist(err) {
			continue
		}
		objectKey := fmt.Sprintf("video/%s/%s_video.m4s", info.DirName, qName)
		files = append(files, uploadFile{localPath: videoFile, ossKey: objectKey})
	}

	// 收集音频文件
	for _, af := range audioFiles {
		if _, err := os.Stat(af.LocalPath); os.IsNotExist(err) {
			continue
		}
		files = append(files, uploadFile{localPath: af.LocalPath, ossKey: af.OSSKey})
	}

	// 上传到主 OSS
	total := len(files)
	for i, f := range files {
		if err := w.storage.PutObjectFromFile(f.ossKey, f.localPath); err != nil {
			return fmt.Errorf("upload %s: %w", f.ossKey, err)
		}
		defer os.Remove(f.localPath)

		// 每上传一个文件更新一次进度
		pct := int(float64(i+1) / float64(total) * 100)
		uploadInfo := map[string]interface{}{
			"ossType":  w.cfg.Storage.OssType,
			"progress": pct,
			"status":   "uploading",
		}
		if data, err := json.Marshal(uploadInfo); err == nil {
			w.retryHSet(ctx, statusKey, "upload", string(data))
		}
		zap.L().Info("Upload progress",
			zap.String("resourceKey", statusKey),
			zap.Int("fileIndex", i+1),
			zap.Int("totalFiles", total),
			zap.Int("progress", pct))
	}

	// 上传到备用 OSS（多源容灾，失败不阻塞主流程）
	if w.backupStorage != nil {
		for _, f := range files {
			if err := w.backupStorage.PutObjectFromFile(f.ossKey, f.localPath); err != nil {
				zap.L().Error("Backup OSS upload failed",
					zap.String("key", f.ossKey),
					zap.Error(err),
				)
			}
		}
	}

	return nil
}

// audioFileEntry 音频上传条目
type audioFileEntry struct {
	LocalPath string // 本地临时文件路径
	OSSKey    string // OSS 目标路径
	Language  string // 语言代码
}

// =============================================================================
// 音频编码
// =============================================================================

// encodeAudio 编码音轨（单音轨或多音轨），返回上传条目列表。
//   - 多音轨模式：每个 AudioStream 独立编码，主音轨（第一条）失败视为整体失败。
//   - 单音轨模式：使用 info.AudioBitRate / AudioSampleRate / AudioChannels 回退。
func (w *Worker) encodeAudio(ctx context.Context, info *dto.TranscodingInfo) ([]audioFileEntry, error) {
	leadMs := ffmpeg.BFramePresentationLeadMs(info.FPS30)
	var files []audioFileEntry

	if len(info.AudioStreams) > 0 {
		// 多音轨模式
		for i, stream := range info.AudioStreams {
			audioFileName := ffmpeg.AudioFileNameForTrack(stream.Language)
			audioOutput := filepath.Join(w.workDir,
				fmt.Sprintf("alnitak_a_%d_%s.m4s", info.ResourceID, stream.Language))

			args := ffmpeg.AudioTrackEncodeArgs(info.InputFile, stream.StreamIndex,
				stream.BitRate, stream.SampleRate, stream.Channels,
				info.Duration, leadMs)
			args = append(args, "-y", audioOutput)

			if err := ffmpeg.EncodeAudio(ctx, args); err != nil {
				if i == 0 {
					return nil, fmt.Errorf("primary audio track %s: %w", stream.Language, err)
				}
				// 附加音轨失败仅记日志
				zap.L().Warn("Secondary audio track encode failed",
					zap.String("language", stream.Language),
					zap.Error(err))
				continue
			}

			ossKey := fmt.Sprintf("video/%s/%s", info.DirName, audioFileName)
			files = append(files, audioFileEntry{
				LocalPath: audioOutput,
				OSSKey:    ossKey,
				Language:  stream.Language,
			})
		}
		return files, nil
	}

	// 单音轨模式（向后兼容）
	audioOutput := filepath.Join(w.workDir,
		fmt.Sprintf("alnitak_a_%d.m4s", info.ResourceID))

	args := ffmpeg.AudioEncodeArgs(info.InputFile, info.AudioBitRate,
		info.AudioSampleRate, info.AudioChannels,
		info.Duration, leadMs)
	args = append(args, "-y", audioOutput)

	if err := ffmpeg.EncodeAudio(ctx, args); err != nil {
		return nil, fmt.Errorf("audio encode: %w", err)
	}

	files = append(files, audioFileEntry{
		LocalPath: audioOutput,
		OSSKey:    fmt.Sprintf("video/%s/audio.m4s", info.DirName),
		Language:  "und",
	})
	return files, nil
}

// =============================================================================
// 视频编码
// =============================================================================

// computeTargets 计算转码目标列表
func (w *Worker) computeTargets(info *dto.TranscodingInfo) []ffmpeg.TranscodingTarget {
	generate1080p60 := info.Generate1080p60 || w.cfg.Transcoding.Generate1080p60
	return ffmpeg.GetTranscodingTargets(
		info.Width, info.Height, info.VideoBitRate,
		info.FPS30, info.FPS60,
		generate1080p60,
	)
}

// encodeQuality 执行单画质视频编码（含 GPU→CPU 逐级降级）。
func (w *Worker) encodeQuality(ctx context.Context, info *dto.TranscodingInfo, target ffmpeg.TranscodingTarget, outputPath string) error {
	// 使用主服务下发的编码参数，保证与后台设置一致。
	// 旧任务（升级前入队）各字段为 false，等效 CPU+H.264。
	useGpu := info.UseGpu
	useHevc := info.UseH265
	useAv1 := info.UseAv1

	return w.doEncodeWithFallback(ctx, info, target, outputPath, useGpu, useAv1, useHevc)
}

// doEncodeWithFallback 执行编码并处理 GPU AV1→H.265→H.264→CPU 逐级降级。
// 全局信号量保护：
//   - nvencSem: 不超过 GeForce 驱动 8 路 NVENC 上限
//   - cpuSem:   CPU fallback 最多 2 路，防止线程打满
func (w *Worker) doEncodeWithFallback(ctx context.Context, info *dto.TranscodingInfo, target ffmpeg.TranscodingTarget, outputPath string, useGpu, useAv1, useHevc bool) error {
	encode := func(gpu, av1, hevc bool) error {
		// GPU 编码竞争 nvencSem，CPU 编码竞争 cpuSem
		if gpu {
			if w.nvencSem != nil {
				select {
				case w.nvencSem <- struct{}{}:
				case <-ctx.Done():
					return ctx.Err()
				}
				defer func() { <-w.nvencSem }()
			}
		} else {
			if w.cpuSem != nil {
				select {
				case w.cpuSem <- struct{}{}:
				case <-ctx.Done():
					return ctx.Err()
				}
				defer func() { <-w.cpuSem }()
			}
		}
		args := ffmpeg.VideoEncodeArgs(info.InputFile, target.Resolution, target.BitrateRate, target.FPS, info.Duration, gpu, av1, hevc)

		numCPU := runtime.NumCPU()
		threads := 4
		if !gpu {
			threads = (numCPU - 2) / w.concurrency
			if threads < 1 {
				threads = 1
			}
			if threads > 8 {
				threads = 8
			}
		}

		args = append(args, "-progress", "pipe:1", "-nostats",
			"-threads", strconv.Itoa(threads),
			"-y", outputPath)

		_, err := ffmpeg.EncodeVideo(ctx, args, info.Duration, func(pct float64) {
			statusKey := fmt.Sprintf("%s%d", transcodingStatusPrefix, info.ResourceID)
			qName := target.Resolution + "_" + target.BitrateRate + "_" + target.FpsName
			w.retryHSet(ctx, statusKey,
				fmt.Sprintf("progress_%s", qName), fmt.Sprintf("%.1f", pct),
				"updated", time.Now().Unix(),
			)
		})
		return err
	}

	err := encode(useGpu, useAv1, useHevc)
	if useGpu && err != nil && (strings.Contains(err.Error(), "GPU error") || strings.Contains(err.Error(), "nvenc")) {
		zap.L().Warn("GPU encode failed, attempting fallback",
			zap.String("quality", target.Resolution+"_"+target.FpsName),
			zap.Error(err))

		if ctx.Err() != nil {
			return ctx.Err()
		}

		// GPU AV1 → GPU H.265
		if useAv1 {
			zap.L().Info("GPU fallback: AV1 → H.265")
			err = encode(true, false, true)
			if err == nil {
				return nil
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}
		}

		// GPU H.265 → GPU H.264
		if useAv1 || useHevc {
			zap.L().Info("GPU fallback: H.265 → H.264")
			err = encode(true, false, false)
			if err == nil {
				return nil
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}
		}

		// GPU H.264 → CPU H.264
		zap.L().Info("GPU fallback: GPU → CPU")
		err = encode(false, false, false)
	}
	return err
}

// qualityIndexData 单画质索引信息（存储到 Redis 供服务端创建 VideoIndexFile 记录）。
type qualityIndexData struct {
	VideoInitRange  string `json:"vInit"`
	VideoIndexRange string `json:"vIndex"`
	AudioInitRange  string `json:"aInit"`
	AudioIndexRange string `json:"aIndex"`
	VideoCodec      string `json:"codec"`
}

// collectQualityIndex 编码完成后从本地文件解析 MP4 init/index range 和编码信息。
// 文件在 uploadOutputs 删除前仍有效，结果存入 Redis 供服务端后续创建索引记录。
func (w *Worker) collectQualityIndex(videoFile, audioFile string, info *dto.TranscodingInfo, qualityName string) *qualityIndexData {
	d := &qualityIndexData{}

	if vInit, vIndex, err := ffmpeg.GetMP4InitRange(videoFile); err == nil {
		d.VideoInitRange = vInit
		d.VideoIndexRange = vIndex
	} else {
		zap.L().Warn("collect video index failed", zap.String("quality", qualityName), zap.Error(err))
		return nil
	}

	if audioFile != "" {
		if aInit, aIndex, err := ffmpeg.GetMP4InitRange(audioFile); err == nil {
			d.AudioInitRange = aInit
			d.AudioIndexRange = aIndex
		} else {
			zap.L().Warn("collect audio index failed", zap.String("quality", qualityName), zap.Error(err))
		}
	}

	// 从已编码文件探测实际 codec（正确处理 GPU 降级场景）
	d.VideoCodec = ffmpeg.ProbeVideoActualCodec(videoFile)

	return d
}

// =============================================================================
// 取消
// =============================================================================

// handleCancel 处理取消信号
func (w *Worker) handleCancel(ctx context.Context, videoIDStr string) {
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		return
	}

	w.mu.RLock()
	cancel, ok := w.activeVideos[uint(videoID)]
	w.mu.RUnlock()

	if ok {
		cancel()
		zap.L().Info("Cancelled job", zap.Uint64("videoID", videoID))
	}
}
