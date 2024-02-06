package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectVideoRoutes(r *gin.RouterGroup) {
	videoGroup := r.Group("video")

	videoAuth := videoGroup.Use(middleware.Auth())
	{
		// 上传视频信息
		videoAuth.POST("/uploadVideoInfo", api.UploadVideoInfo)
		// 获取上传视频状态信息
		videoAuth.GET("/getVideoStatus", api.GetVideoStatus)
	}
}
