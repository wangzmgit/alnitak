package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectVerifyRoutes(r *gin.RouterGroup) {

	verifyGroup := r.Group("verify")
	{
		// 获取滑块验证
		verifyGroup.GET("captcha/get", api.GetSliderCaptcha)
		// 验证滑块
		verifyGroup.POST("captcha/validate", api.ValidateSlider)
		// 获取邮箱验证码（限制每 IP 每分钟 3 次，防短信轰炸）
		verifyGroup.POST("getEmailCode", middleware.RateLimiter(middleware.EmailCodeRateLimit), api.SendRegisterEmailCode)
	}

}
