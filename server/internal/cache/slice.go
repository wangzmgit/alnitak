package cache

import "interastral-peace.com/alnitak/internal/global"

func GetVideoSlice(key string) string {
	return global.Redis.Get(VIDEO_SLICE_STATUS + key)
}

func SetVideoSlice(key, value string) {
	global.Redis.Set(VIDEO_SLICE_STATUS+key, value, VIDEO_SLICE_EXPRIRATION_TIME)
}

func DelVideoSlice(key string) {
	global.Redis.Del(VIDEO_SLICE_STATUS + key)
}
