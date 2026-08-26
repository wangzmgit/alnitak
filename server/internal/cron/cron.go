package cron

import (
	"context"

	"github.com/jasonlvhit/gocron"
	"go.uber.org/zap"
	"interastral-peace.com/alnitak/internal/service"
)

// safeDo 包装定时任务函数，自动 recover panic 避免单个任务炸掉整个 cron goroutine
func safeDo(fn func()) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				zap.L().Error("cron任务panic恢复",
					zap.Any("recover", r),
					zap.Stack("stack"))
			}
		}()
		fn()
	}
}

func StartCronTask() {
	c := gocron.NewScheduler()

	// 每3小时刷新同步播放量数据
	c.Every(3).Hours().Do(safeDo(SyncClicks))

	// 每3小时刷新一次热点
	c.Every(3).Hours().Do(safeDo(RefreshPopular))

	// 每天晚上12点解封用户
	c.Every(1).Day().At("00:00").Do(safeDo(UnbanUser))

	// 每天凌晨3点清理孤立资源
	c.Every(1).Day().At("03:00").Do(safeDo(CleanupOrphanedResources))

	// 每30分钟重试备用OSS上传失败记录
	c.Every(30).Minutes().Do(safeDo(RetryBackupFailures))

	// 每5分钟恢复卡住的转码任务（Enqueue失败但状态未更新）
	c.Every(5).Minutes().Do(safeDo(RecoverStuckTranscoding))

	// 每10秒从 pending 队列派发任务到 Stream（兜底 Worker 完成通知丢失的场景）
	c.Every(10).Seconds().Do(safeDo(DispatchPendingJobs))

	<-c.Start()
}

// DispatchPendingJobs 兜底：定期检查 pending 队列，有空闲 Worker 则派发
func DispatchPendingJobs() {
	t := service.GetCurrentTranscoder()
	if t == nil {
		return
	}
	t.DispatchPending(context.Background())
}

// ListenDispatch 启动监听 Worker 完成通知的 goroutine，实时触发 pending 派发
func StartDispatchListener() {
	t := service.GetCurrentTranscoder()
	if t == nil {
		return
	}
	go t.ListenDispatch(context.Background())
	zap.L().Info("Started dispatch listener for pending transcoding jobs")
}
