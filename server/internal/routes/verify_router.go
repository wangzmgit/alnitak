package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
)

func CollectVerifyRoutes(r *gin.RouterGroup) {

	verifyGroup := r.Group("verify")
	{
		// 获取滑块验证
		verifyGroup.GET("captcha/get", api.GetSliderCaptcha)
		// 验证滑块
		verifyGroup.POST("captcha/validate", api.ValidateSlider)
		// 获取邮箱验证码
		verifyGroup.POST("getEmailCode", api.SendRegisterEmailCode)
	}

}
