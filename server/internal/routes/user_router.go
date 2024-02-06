package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectUserRoutes(r *gin.RouterGroup) {
	userGroup := r.Group("user")
	userGroup.GET("/getUserBaseInfo", api.GetUserBaseInfo)

	userAuth := userGroup.Use(middleware.Auth())
	{
		// 用户获取个人信息
		userAuth.GET("/getUserInfo", api.GetUserInfo)
	}
}
