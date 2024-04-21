package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectResourceRoutes(r *gin.RouterGroup) {
	resourceGroup := r.Group("resource")

	resourceAuth := resourceGroup.Group("")
	resourceAuth.Use(middleware.Auth())
	{
		resourceAuth.PUT("modifyTitle", api.ModifyResourceTitle)
		resourceAuth.DELETE("deleteResource/:id", api.DeleteResource)
	}
}
