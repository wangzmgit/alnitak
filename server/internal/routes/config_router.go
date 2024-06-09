package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectConfigRoutes(r *gin.RouterGroup) {
	configGroup := r.Group("config")

	configAuth := configGroup.Group("")
	configAuth.Use(middleware.Auth())
	{
		// 获取邮箱配置
		configAuth.GET("getEmailConfig", api.GetEmailConfig)
		// 修改邮箱配置
		configAuth.POST("setEmailConfig", api.SetEmailConfig)
		// 获取存储配置
		configAuth.GET("getStorageConfig", api.GetStorageConfig)
		// 修改存储配置
		configAuth.POST("setStorageConfig", api.SetStorageConfig)
		// 获取其他配置
		configAuth.GET("getOtherConfig", api.GetOtherConfig)
		// 修改其他配置
		configAuth.POST("setOtherConfig", api.SetOtherConfig)
	}

}
