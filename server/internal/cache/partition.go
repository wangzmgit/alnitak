package cache

import (
	"encoding/json"

	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func GetPartition() (partitions []vo.PartitionResp) {
	jsonStr := global.Redis.Get(PARTITION_KEY)
	if jsonStr == "" {
		return
	}

	// 反序列化
	if err := json.Unmarshal([]byte(jsonStr), &partitions); err != nil {
		utils.ErrorLog("分区反序列化失败", "cache", err.Error())
	}
	return
}

func SetPartition(partitions []vo.PartitionResp) {
	//先序列化
	pb, err := json.Marshal(partitions)
	if err != nil {
		utils.ErrorLog("分区序列化失败", "cache", err.Error())
		return
	}
	global.Redis.Set(PARTITION_KEY, pb, PARTITION_EXPRIRATION_TIME)
}

func DelPartition() {
	global.Redis.Del(PARTITION_KEY)
}
