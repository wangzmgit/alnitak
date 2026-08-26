package cache

import (
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

// TryAcquireShareVideoLimit 为 (userId, videoId) 维度获取一次分享计数的许可。
// 返回 true 表示成功占位，可以执行 shares+1；返回 false 表示冷却期内重复触发，应拒绝。
func TryAcquireShareVideoLimit(userId, videoId uint) bool {
	key := SHARE_VIDEO_LIMIT_KEY + utils.UintToString(userId) + ":" + utils.UintToString(videoId)
	return global.Redis.SetNX(key, 1, SHARE_LIMIT_EXPIRATION_TIME)
}

// TryAcquireShareArticleLimit 为 (userId, articleId) 维度获取一次分享计数的许可。
func TryAcquireShareArticleLimit(userId, articleId uint) bool {
	key := SHARE_ARTICLE_LIMIT_KEY + utils.UintToString(userId) + ":" + utils.UintToString(articleId)
	return global.Redis.SetNX(key, 1, SHARE_LIMIT_EXPIRATION_TIME)
}
