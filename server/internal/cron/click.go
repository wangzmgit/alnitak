package cron

import (
	"strings"
	"time"

	"go.uber.org/zap"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 同步播放量到数据库
func SyncClicks() {
	start := time.Now() // 获取当前时间
	zap.L().Info("开始同步播放量", zap.String("module", "cron"))
	keys := cache.GetClicksKeys()
	for _, v := range keys {
		// 获取视频ID和播放量
		videoId := utils.StringToUint(v[strings.LastIndex(v, ":")+1:])
		clicks, err := cache.GetClicks(videoId)
		if err != nil {
			continue
		}
		// 写入数据库
		service.UpdateClicks(videoId, clicks)
		cache.SetClicks(videoId, 0)
	}
	zap.L().Info("播放量同步完成，耗时:"+time.Since(start).String(), zap.String("module", "cron"))
}
