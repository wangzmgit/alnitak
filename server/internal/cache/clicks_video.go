package cache

import (
	"errors"
	"strconv"
	"time"

	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func GetVideoClicksLimit(videoId uint, ip string) string {
	s := global.Redis.Get(VIDEO_CLICKS_LIMIT_KEY + utils.UintToString(videoId) + ":" + ip)

	return s
}

func SetVideoClicksLimit(videoId uint, ip string) {
	global.Redis.Set(VIDEO_CLICKS_LIMIT_KEY+utils.UintToString(videoId)+":"+ip, 1,
		VIDEO_CLICKS_LIMIT_EXPRIRATION_TIME)
}

func GetVideoClicks(videoId uint) (int64, error) {
	s := global.Redis.Get(VIDEO_CLICKS_KEY + utils.UintToString(videoId))
	if s == "" {
		return 0, errors.New("no record")
	}

	count, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		utils.ErrorLog("播放量转换为int64类型失败", "cache", err.Error())
		return 0, err
	}

	return count, nil
}

func SetVideoClicks(videoId uint, count int64) {
	global.Redis.Set(VIDEO_CLICKS_KEY+utils.UintToString(videoId), count, VIDEO_CLICKS_EXPRIRATION_TIME)
}

// 删除播放量
func DelVideoClicks(videoId uint) {
	global.Redis.Del(VIDEO_CLICKS_KEY + utils.UintToString(videoId))
}

func AddVideoClicks(videoId uint) {
	global.Redis.Incr(VIDEO_CLICKS_KEY + utils.UintToString(videoId))
}

func GetVideoClicksKeys() []string {
	return global.Redis.Keys(VIDEO_CLICKS_KEY + "*")
}

// 获取点击量过期时间
func VideoClickTTL(videoId uint) time.Duration {
	return global.Redis.TTL(VIDEO_CLICKS_KEY + utils.UintToString(videoId))
}
