package cache

import (
	"encoding/json"

	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func GetUserInfo(id uint) (user vo.UserInfoResp) {
	jsonStr := global.Redis.Get(USER_INFO_KEY + utils.UintToString(id))
	// 反序列化
	if err := json.Unmarshal([]byte(jsonStr), &user); err != nil {
		utils.ErrorLog("用户信息反序列化失败", "cache", err.Error())
	}
	return
}

func SetUser(user vo.UserInfoResp) {
	//先序列化user
	ub, err := json.Marshal(user)
	if err != nil {
		utils.ErrorLog("用户信息序列化失败", "cache", err.Error())
		return
	}

	global.Redis.Set(USER_INFO_KEY+utils.UintToString(user.ID), ub, USER_INFO_EXPIRATION_TIME)
}
