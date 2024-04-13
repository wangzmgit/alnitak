package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectVideoRoutes(r *gin.RouterGroup) {
	videoGroup := r.Group("video")

	videoAuth := videoGroup.Group("")
	videoAuth.Use(middleware.Auth())
	{
		// 上传视频信息
		videoAuth.POST("uploadVideoInfo", api.UploadVideoInfo)
		// 获取上传视频状态信息
		videoAuth.GET("getVideoStatus", api.GetVideoStatus)
		// 提交审核
		videoAuth.POST("submitReview", api.SubmitReview)
		// 获取上传的视频
		videoAuth.GET("getUploadVideo", api.GetUploadVideoList)
		// 编辑视频信息
		videoAuth.PUT("editVideoInfo", api.EditVideoInfo)
		// 删除视频
		videoAuth.DELETE("deleteVideo/:id", api.DeleteVideo)
		// 获取所有的视频列表
		videoAuth.GET("getAllVideoList", api.GetAllVideoList)
		// 获取审核列表

		// 获取视频列表（后台管理）
		videoAuth.POST("getVideoListManage", api.GetVideoListManage)
		// 删除视频（后台管理）
		videoAuth.DELETE("deleteVideoManage/:id", api.DeleteVideoManage)
	}

	// 获取视频信息
	videoGroup.GET("getVideoById", api.GetVideoById)
	// 获取视频资源支持的分辨率信息
	videoGroup.GET("getResourceQuality", api.GetResourceQuality)
	// 获取视频文件
	videoGroup.GET("getVideoFile", api.GetVideoFile)
	// 获取视频切片
	videoGroup.GET("slice/:file", api.GetVideoSlice)
	// 获取用户视频
	videoGroup.GET("getVideoByUser", api.GetVideoByUser)
}
