package cache

import (
	"encoding/json"
	"strconv"

	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func GetVideoInfo(id uint) (video vo.VideoResp) {
	jsonStr := global.Redis.Get(VIDEO_INFO_KEY + utils.UintToString(id))
	if jsonStr == "" {
		return
	}
	// 反序列化
	if err := json.Unmarshal([]byte(jsonStr), &video); err != nil {
		utils.ErrorLog("视频信息反序列化失败", "cache", err.Error())
	}
	return
}

func SetVideoInfo(video vo.VideoResp) {
	//先序列化video
	ub, err := json.Marshal(video)
	if err != nil {
		utils.ErrorLog("视频信息序列化失败", "cache", err.Error())
		return
	}

	global.Redis.Set(VIDEO_INFO_KEY+utils.UintToString(video.ID), ub, VIDEO_INFO_EXPIRATION_TIME)
}

func DelVideoInfo(id uint) {
	global.Redis.Del(VIDEO_INFO_KEY + utils.UintToString(id))
}

func SetVideoId(partitionId, videoId uint) {
	global.Redis.SAdd(ALL_VIDEO_KEY+strconv.Itoa(int(partitionId)), videoId)
}

func DelAllVideoId() {
	keys := global.Redis.Keys(ALL_VIDEO_KEY + "*")
	for _, key := range keys {
		global.Redis.Del(key)
	}
}

func DelVideoId(partitionId, videoId uint) {
	global.Redis.SRem(ALL_VIDEO_KEY+strconv.Itoa(int(partitionId)), videoId)
}

func SetHotVideoId(videoId uint) {
	global.Redis.SAdd(HOT_VIDEO_KEY, videoId)
}

func DelHotVideoId() {
	global.Redis.Del(HOT_VIDEO_KEY)
}

func GetHotVideoId() []string {
	return global.Redis.SMembers(HOT_VIDEO_KEY)
}

func GetVideoIdByPartition(partitionId uint, count int64) []string {
	return global.Redis.SRandMemberN(ALL_VIDEO_KEY+strconv.Itoa(int(partitionId)), count)
}
