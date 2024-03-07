package cache

import (
	"interastral-peace.com/alnitak/internal/global"
)

const (
	VIDEO_SLICE_NOT_USED = "1" // 未使用
)

func GetVideoSliceStatus(key string) string {
	return global.Redis.Get(VIDEO_SLICE_STATUS + key)
}

func SetVideoSliceStatus(key string) {
	global.Redis.Set(VIDEO_SLICE_STATUS+key, VIDEO_SLICE_NOT_USED, VIDEO_SLICE_EXPRIRATION_TIME)
}

func DelVideoSliceStatus(key string) {
	global.Redis.Del(VIDEO_SLICE_STATUS + key)
}
