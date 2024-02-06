package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectUploadRoutes(r *gin.RouterGroup) {

	uploadGroup := r.Group("upload")
	uploadGroup.Use(middleware.Auth())
	{
		uploadGroup.POST("image", api.UploadImg)
		uploadGroup.POST("video/:vid", api.UploadVideo)
	}
}
