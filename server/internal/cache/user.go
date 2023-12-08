package cache

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func GetUserInfo(id uint) (user vo.UserInfoResp) {
	jsonStr := global.Redis.Get(USER_INFO_KEY + utils.UintToString(id))
	// 反序列化
	if err := json.Unmarshal([]byte(jsonStr), &user); err != nil {
		zap.L().Error(fmt.Sprintf("Cache | 用户信息反序列化失败 | 错误信息: %s", err.Error()))
	}
	return
}

func SetUser(user vo.UserInfoResp) {
	//先序列化user
	ub, err := json.Marshal(user)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Cache | 用户信息序列化失败 | 错误信息: %s", err.Error()))
		return
	}

	global.Redis.Set(USER_INFO_KEY+utils.UintToString(user.ID), ub, USER_INFO_EXPIRATION_TIME)
}
