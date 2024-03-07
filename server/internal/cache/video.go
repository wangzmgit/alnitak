package cache

import (
	"encoding/json"

	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func GetVideoInfo(id uint) (video vo.VideoResp) {
	jsonStr := global.Redis.Get(VIDEO_INFO_KEY + utils.UintToString(id))
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
