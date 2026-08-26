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
		// 图片上传（速率 + 10MB body 限制）
		uploadGroup.POST("image",
			middleware.MaxBodySize(middleware.UploadImgMaxBody),
			middleware.RateLimiter(middleware.UploadImgRateLimit),
			api.UploadImg)

		// 新建视频 / 新增分P（速率 + 1MB body 限制）
		uploadGroup.POST("video/:vid",
			middleware.MaxBodySize(middleware.UploadJSONMaxBody),
			middleware.RateLimiter(middleware.UploadCreateRateLimit),
			api.UploadVideoAdd)
		uploadGroup.POST("video",
			middleware.MaxBodySize(middleware.UploadJSONMaxBody),
			middleware.RateLimiter(middleware.UploadCreateRateLimit),
			api.UploadVideoCreate)

		// 查询分片（速率 + 1MB body 限制）
		uploadGroup.POST("checkVideo",
			middleware.MaxBodySize(middleware.UploadJSONMaxBody),
			middleware.RateLimiter(middleware.UploadCheckRateLimit),
			api.UploadVideoCheck)

		// 分片上传视频（速率 + 50MB body 限制）
		uploadGroup.POST("chunkVideo",
			middleware.MaxBodySize(middleware.UploadChunkMaxBody),
			middleware.RateLimiter(middleware.UploadChunkRateLimit),
			api.UploadVideoChunk)

		// 合并视频分片（速率 + 1MB body 限制）
		uploadGroup.POST("mergeVideo",
			middleware.MaxBodySize(middleware.UploadJSONMaxBody),
			middleware.RateLimiter(middleware.UploadChunkRateLimit),
			api.UploadVideoMerge)

		// ── 直传 OSS 接口 ──

		// 图片：获取预签名 URL（1MB body）
		uploadGroup.POST("presignImage",
			middleware.MaxBodySize(middleware.UploadJSONMaxBody),
			middleware.RateLimiter(middleware.UploadImgRateLimit),
			api.PresignImageUpload)

		// 图片：确认上传完成（1MB body）
		uploadGroup.POST("confirmImage",
			middleware.MaxBodySize(middleware.UploadJSONMaxBody),
			middleware.RateLimiter(middleware.UploadImgRateLimit),
			api.ConfirmImageUpload)

		// 视频：初始化分片直传（1MB body）— 返回第一批20个预签名URL
		uploadGroup.POST("initVideo",
			middleware.MaxBodySize(middleware.UploadJSONMaxBody),
			middleware.RateLimiter(middleware.UploadCreateRateLimit),
			api.InitVideoUpload)

		// 视频：续签下一批分片预签名 URL（1MB body）
		uploadGroup.POST("presignChunks",
			middleware.MaxBodySize(middleware.UploadJSONMaxBody),
			middleware.RateLimiter(middleware.UploadChunkRateLimit),
			api.PresignUploadChunks)

		// 视频：完成分片直传（1MB body）
		uploadGroup.POST("completeVideo",
			middleware.MaxBodySize(middleware.UploadJSONMaxBody),
			middleware.RateLimiter(middleware.UploadChunkRateLimit),
			api.CompleteVideoUpload)
	}
}
