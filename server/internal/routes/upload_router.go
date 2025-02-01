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
		uploadGroup.POST("checkVideo", api.UploadVideoCheck)
		uploadGroup.POST("chunkVideo", api.UploadVideoChunk) // 分片上传视频
		uploadGroup.POST("mergeVideo", api.UploadVideoMerge) // 合并视频分片

		uploadGroup.POST("video/:vid", api.UploadVideoAdd)
		uploadGroup.POST("video", api.UploadVideoCreate)
	}
}
