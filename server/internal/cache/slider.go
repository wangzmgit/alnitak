package cache

import (
	"strconv"

	"go.uber.org/zap"
	"interastral-peace.com/alnitak/internal/global"
)

func GetSliderX(captchaId string) int {
	s := global.Redis.Get(SLIDER_X_KEY + captchaId)
	x, err := strconv.Atoi(s)
	if err != nil {
		zap.L().Error("数据错误 | 缓存中滑块x坐标转换为int类型失败")
	}
	return x
}

func SetSliderX(captchaId string, x int) {
	global.Redis.Set(SLIDER_X_KEY+captchaId, x, SLIDER_X_EXPRIRATION_TIME)
}

func DelSlider(captchaId string) {
	global.Redis.Del(SLIDER_X_KEY + captchaId)
}
