package service

import (
	"gorm.io/gorm"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
)

// 获取播放量
func GetVideoClicks(videoId uint) int64 {
	clicks, err := cache.GetVideoClicks(videoId)
	if err != nil {
		cache.SetVideoClicks(videoId, 0)
	}

	return clicks
}

// 增加播放量
func AddVideoClicks(videoId uint, ip string) {
	// 一个ip在同一个视频下，每30分钟可重新增加1播放量
	if cache.GetVideoClicksLimit(videoId, ip) == "" {
		cache.AddVideoClicks(videoId)
		cache.SetVideoClicksLimit(videoId, ip)
	}
}

// 更新播放量
func UpdateVideoClicks(videoId uint, clicks int64) error {
	return global.Mysql.Model(&model.Video{}).Where("id = ?", videoId).
		UpdateColumn("clicks", gorm.Expr("clicks + ?", clicks)).Error
}

// 获取播放量
func GetArticleClicks(articleId uint) int64 {
	clicks, err := cache.GetArticleClicks(articleId)
	if err != nil {
		cache.SetArticleClicks(articleId, 0)
	}

	return clicks
}

// 增加播放量
func AddArticleClicks(articleId uint, ip string) {
	// 一个ip在同一个视频下，每30分钟可重新增加1播放量
	if cache.GetArticleClicksLimit(articleId, ip) == "" {
		cache.AddArticleClicks(articleId)
		cache.SetArticleClicksLimit(articleId, ip)
	}
}

// 更新文章点击量
func UpdateArticleClicks(articleId uint, clicks int64) error {
	return global.Mysql.Model(&model.Article{}).Where("id = ?", articleId).
		UpdateColumn("clicks", gorm.Expr("clicks + ?", clicks)).Error
}
