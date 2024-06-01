package cache

import (
	"encoding/json"

	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func GetVideoPartition() (partitions []vo.PartitionResp) {
	jsonStr := global.Redis.Get(VIDEO_PARTITION_KEY)
	if jsonStr == "" {
		return
	}

	// 反序列化
	if err := json.Unmarshal([]byte(jsonStr), &partitions); err != nil {
		utils.ErrorLog("分区反序列化失败", "cache", err.Error())
	}
	return
}

func SetVideoPartition(partitions []vo.PartitionResp) {
	//先序列化
	pb, err := json.Marshal(partitions)
	if err != nil {
		utils.ErrorLog("分区序列化失败", "cache", err.Error())
		return
	}
	global.Redis.Set(VIDEO_PARTITION_KEY, pb, VIDEO_PARTITION_EXPRIRATION_TIME)
}

func DelVideoPartition() {
	global.Redis.Del(VIDEO_PARTITION_KEY)
}

func GetArticlePartition() (partitions []vo.PartitionResp) {
	jsonStr := global.Redis.Get(ARTICLE_PARTITION_KEY)
	if jsonStr == "" {
		return
	}

	// 反序列化
	if err := json.Unmarshal([]byte(jsonStr), &partitions); err != nil {
		utils.ErrorLog("分区反序列化失败", "cache", err.Error())
	}
	return
}

func SetArticlePartition(partitions []vo.PartitionResp) {
	//先序列化
	pb, err := json.Marshal(partitions)
	if err != nil {
		utils.ErrorLog("分区序列化失败", "cache", err.Error())
		return
	}
	global.Redis.Set(ARTICLE_PARTITION_KEY, pb, ARTICLE_PARTITION_EXPRIRATION_TIME)
}

func DelArticlePartition() {
	global.Redis.Del(ARTICLE_PARTITION_KEY)
}
