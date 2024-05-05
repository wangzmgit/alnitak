package cron

import (
	"strings"
	"time"

	"go.uber.org/zap"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

func SyncClicks() {
	syncVideoClicks()
	syncArticleClicks()
}

// 同步视频播放量到数据库
func syncVideoClicks() {
	start := time.Now() // 获取当前时间
	zap.L().Info("开始同步播放量", zap.String("module", "cron"))
	keys := cache.GetVideoClicksKeys()
	for _, v := range keys {
		// 获取视频ID和播放量
		videoId := utils.StringToUint(v[strings.LastIndex(v, ":")+1:])
		clicks, err := cache.GetVideoClicks(videoId)
		if err != nil || clicks == 0 {
			continue
		}
		// 写入数据库
		service.UpdateVideoClicks(videoId, clicks)

		// 更新缓存中的播放量
		video := cache.GetVideoInfo(videoId)
		video.Clicks += clicks
		cache.SetVideoInfo(video)

		// 清除当前增加的播放量
		cache.SetVideoClicks(videoId, 0)
	}
	zap.L().Info("播放量同步完成，耗时:"+time.Since(start).String(), zap.String("module", "cron"))
}

// 同步文章点击量到数据库
func syncArticleClicks() {
	start := time.Now() // 获取当前时间
	zap.L().Info("开始同步文章点击量", zap.String("module", "cron"))
	keys := cache.GetArticleClicksKeys()
	for _, v := range keys {
		// 获取视频ID和播放量
		articleId := utils.StringToUint(v[strings.LastIndex(v, ":")+1:])
		clicks, err := cache.GetArticleClicks(articleId)
		if err != nil || clicks == 0 {
			continue
		}
		// 写入数据库
		service.UpdateArticleClicks(articleId, clicks)

		// 清除当前增加的播放量
		cache.SetArticleClicks(articleId, 0)
	}
	zap.L().Info("文章点击量同步完成，耗时:"+time.Since(start).String(), zap.String("module", "cron"))
}
