package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectUserRoutes(r *gin.RouterGroup) {
	userApi := api.UserApi{}

	user := r.Group("user")
	user.GET("/getUserBaseInfo", userApi.GetUserBaseInfo)

	auth := user.Use(middleware.Auth())
	{
		// 用户获取个人信息
		auth.GET("/getUserInfo", userApi.GetUserInfo)
	}
}
