package cache

import (
	"strconv"

	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func GetLoginTryCount(username string) int {
	s := global.Redis.Get(LOGIN_TRY_COUNT_KEY + username)
	count, err := strconv.Atoi(s)
	if err != nil {
		utils.ErrorLog("用户登录次数转换int类型失败", "cache", err.Error())
	}
	return count
}

func SetLoginTryCount(username string, count int) {
	global.Redis.Set(LOGIN_TRY_COUNT_KEY+username, count, LOGIN_TRY_COUNT_EXPRIRATION_TIME)
}

func IncrLoginTryCount(username string) {
	global.Redis.Incr(LOGIN_TRY_COUNT_KEY + username)
	global.Redis.Expire(LOGIN_TRY_COUNT_KEY+username, LOGIN_TRY_COUNT_EXPRIRATION_TIME)
}

func DelLoginTryCount(username string) {
	global.Redis.Del(LOGIN_TRY_COUNT_KEY + username)
}
