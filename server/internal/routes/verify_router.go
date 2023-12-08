package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
)

func CollectVerifyRoutes(r *gin.RouterGroup) {
	emailApi := api.EmailApi{}
	captchaApi := api.CaptchaApi{}

	verify := r.Group("verify")
	{
		// 获取滑块验证
		verify.GET("captcha/get", captchaApi.GetSliderCaptcha)
		// 验证滑块
		verify.POST("captcha/validate", captchaApi.ValidateSlider)
		// 获取邮箱验证码
		verify.POST("getEmailCode", emailApi.SendRegisterEmailCode)
	}

}
