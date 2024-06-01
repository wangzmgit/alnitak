package cache

import (
	"strconv"

	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func GetUploadImage(url string) uint {
	s := global.Redis.Get(UPLOAD_IMAGE_KEY + url)
	if s == "" {
		return 0
	}

	userId, err := strconv.Atoi(s)
	if err != nil {
		utils.ErrorLog("用户id转换为uint类型失败", "cache", err.Error())
	}
	return uint(userId)
}

func SetUploadImage(url string, userID uint) {
	global.Redis.Set(UPLOAD_IMAGE_KEY+url, userID, UPLOAD_IMAGE_EXPRIRATION_TIME)
}

func DelUploadImage(url string) {
	global.Redis.Del(UPLOAD_IMAGE_KEY + url)
}
