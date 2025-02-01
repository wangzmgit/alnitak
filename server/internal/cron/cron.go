package cron

import "github.com/jasonlvhit/gocron"

func StartCronTask() {
	c := gocron.NewScheduler()

	// 每3小时刷新同步播放量数据
	c.Every(3).Hours().Do(SyncClicks)

	// 每3小时刷新一次热点
	c.Every(3).Hours().Do(RefreshPopular)

	<-c.Start()
}
