package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectResourceRoutes(r *gin.RouterGroup) {
	resource := r.Group("resource")
	{
		//需要用户登录
		auth := resource.Group("")
		auth.Use(middleware.Auth())
		{
			auth.PUT("modifyTitle", api.ModifyResourceTitle)
			auth.DELETE("deleteResource/:id", api.DeleteResource)
		}
	}
}
