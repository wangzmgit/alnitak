package service

import (
	"gorm.io/gorm"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
)

// 获取播放量
func GetVideoClicks(videoId uint) int64 {
	clicks, err := cache.GetClicks(videoId)
	if err != nil {
		cache.SetClicks(videoId, 0)
	}

	return clicks
}

// 增加播放量
func AddVideoClicks(videoId uint, ip string) {
	// 一个ip在同一个视频下，每30分钟可重新增加1播放量
	if cache.GetClicksLimit(videoId, ip) == "" {
		cache.AddClicks(videoId)
		cache.SetClicksLimit(videoId, ip)
	}
}

// 更新播放量
func UpdateClicks(videoId uint, clicks int64) error {
	return global.Mysql.Model(&model.Video{}).Where("id = ?", videoId).
		UpdateColumn("clicks", gorm.Expr("clicks + ?", clicks)).Error
}
