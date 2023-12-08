package cache

import (
	"time"

	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func IsRefreshTokenExist(userId uint, token string) bool {
	return global.Redis.ZScore(REFRESH_TOKEN_KEY+utils.UintToString(userId), token) != 0
}

func SetRefreshToken(id uint, token string) {
	key := REFRESH_TOKEN_KEY + utils.UintToString(id)
	if global.Redis.ZCard(key) >= MAX_LOGIN_LIMIT {
		// 保留MAX_LOGIN_LIMIT - 1个token
		global.Redis.ZRemRangeByRank(key, 0, 1-MAX_LOGIN_LIMIT)
	}

	global.Redis.ZAdd(key, float64(time.Now().Add(REFRESH_TOKEN_EXPRIRATION_TIME).Unix()), token)
}

func DelRefreshToken(id uint, token string) {
	global.Redis.ZRem(REFRESH_TOKEN_KEY+utils.UintToString(id), token)
}
