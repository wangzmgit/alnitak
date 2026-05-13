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
		videoAuth.POST("uploadVideoInfo", middleware.Ban(), api.UploadVideoInfo)
		// 获取上传视频状态信息
		videoAuth.GET("getVideoStatus", api.GetVideoStatus)
		// 获取上传的视频
		videoAuth.GET("getUploadVideo", api.GetUploadVideoList)
		// 编辑视频信息
		videoAuth.PUT("editVideoInfo", middleware.Ban(), api.EditVideoInfo)
		// 删除视频
		videoAuth.DELETE("deleteVideo/:id", api.DeleteVideo)
		// 获取所有的视频列表
		videoAuth.GET("getAllVideoList", api.GetAllVideoList)
		// 获取审核列表（后台管理）
		videoAuth.POST("getReviewList", api.GetReviewList)
		// 获取审核资源列表（后台管理）
		videoAuth.GET("getReviewResourceList", api.GetReviewResourceList)
		// 获取视频列表（后台管理）
		videoAuth.POST("getVideoListManage", api.GetVideoListManage)
		// 删除视频（后台管理）
		videoAuth.DELETE("deleteVideoManage/:id", api.DeleteVideoManage)
		// 获取视频资源支持的分辨率信息（后台管理）
		videoAuth.GET("getResourceQualityManage", api.GetResourceQualityManage)
		// 获取视频文件（后台管理）
		videoAuth.GET("getVideoFileManage", api.GetVideoFileManage)
		// 获取处理失败的视频列表（后台管理）
		videoAuth.POST("getFailedVideoList", api.GetFailedVideoList)
		// 获取处理中视频列表（后台管理）
		videoAuth.POST("getProcessingVideoList", api.GetProcessingVideoList)
		// 重新转码视频（后台管理专用）
		videoAuth.POST("reTranscodeVideo", api.ReTranscodeVideo)
		// 字幕（UGC 稿件作者）
		videoAuth.POST("subtitle/upload", middleware.Ban(), api.PostSubtitleUpload)
		videoAuth.PUT("subtitle/:id", middleware.Ban(), api.PutSubtitle)
		videoAuth.DELETE("subtitle/:id", middleware.Ban(), api.DeleteSubtitle)
	}

	videoPublicOpt := videoGroup.Group("")
	videoPublicOpt.Use(middleware.OptionalAuth())
	{
		videoPublicOpt.GET("subtitle/list", api.GetSubtitleList)
	}

	// 获取视频信息
	videoGroup.GET("getVideoById", api.GetVideoById)
	// 获取视频资源支持的分辨率信息
	videoGroup.GET("getResourceQuality", api.GetResourceQuality)
	// 获取视频文件
	videoGroup.GET("getVideoFile", api.GetVideoFile)
	// 获取视频切片（兼容模式：SegmentList）
	videoGroup.GET("slice/:file", api.GetVideoSlice)
	// 获取视频流（B站风格：SegmentBase，支持字节范围请求）
	videoGroup.GET("stream/:file", api.GetVideoStream)
	// 获取用户视频
	videoGroup.GET("getVideoByUser", api.GetVideoByUser)
	// 获取热门视频
	videoGroup.GET("getHotVideo", api.GetHotVideo)
	// 获取最近上传的视频
	videoGroup.GET("getLatestVideo", api.GetLatestVideo)
	// 获取分区视频
	videoGroup.GET("getVideoListByPartition", api.GetVideoListByPartition)
	// 获取相关推荐视频
	videoGroup.GET("getRelatedVideoList", api.GetRelatedVideoList)
	// 搜索视频
	videoGroup.POST("searchVideo", api.SearchVideo)
}
