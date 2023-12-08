package cache

import (
	"strconv"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"interastral-peace.com/alnitak/internal/global"
)

// 创建验证状态
func CreateCaptchaStatus() string {
	captchaId := uuid.New().String()
	SetCaptchaStatus(captchaId, global.CAPTCHA_STATUS_NOT_USED)
	return captchaId
}

// 获取验证状态
func GetCaptchaStatus(captchaId string) int { // 0 不存在 1未验证  2验证成功
	s := global.Redis.Get(CAPTCHA_STATUS_KEY + captchaId)
	status, err := strconv.Atoi(s)
	if err != nil {
		zap.L().Error("数据错误 | 人机验证状态转换为int类型失败")
	}
	return status
}

func SetCaptchaStatus(captchaId string, status int) {
	global.Redis.Set(CAPTCHA_STATUS_KEY+captchaId, status, CAPTCHA_STATUS_EXPRIRATION_TIME)
}

func DelCaptchaStatus(captchaId string) {
	global.Redis.Del(CAPTCHA_STATUS_KEY + captchaId)
}
