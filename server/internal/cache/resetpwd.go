package cache

import (
	"strconv"

	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

// 状态1: 验证成功

func GetResetPwdCheckStatus(email string) int {
	s := global.Redis.Get(RESET_PWD_CHECK_KEY + email)
	status, err := strconv.Atoi(s)
	if err != nil {
		utils.ErrorLog("重置密码状态转换为int类型失败", "cache", err.Error())
	}
	return status
}

func SetResetPwdCheckStatus(email string, status int) {
	global.Redis.Set(RESET_PWD_CHECK_KEY+email, status, RESET_PWD_CHECK_EXPRIRATION_TIME)
}

func DelResetPwdCheckStatus(email string) {
	global.Redis.Del(RESET_PWD_CHECK_KEY + email)
}
