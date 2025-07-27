package cron

import (
	"github.com/jasonlvhit/gocron"
	"interastral-peace.com/alnitak/internal/service"
)

func StartCronTask() {
	c := gocron.NewScheduler()

	// 每3小时刷新同步播放量数据
	c.Every(3).Hours().Do(SyncClicks)

	// 每3小时刷新一次热点
	c.Every(3).Hours().Do(RefreshPopular)

	// 每天凌晨2点清理旧的操作日志
	c.Every(1).Day().At("02:00").Do(service.CleanupOldOperateLogsBatch)

	<-c.Start()
}
