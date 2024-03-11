package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectArchiveRoutes(r *gin.RouterGroup) {
	archiveGroup := r.Group("archive")

	archiveAuth := archiveGroup.Group("")
	archiveAuth.Use(middleware.Auth())
	{
		// 点赞
		archiveAuth.POST("like", api.Like)
		// 取消赞
		archiveAuth.POST("cancelLike", api.CancelLike)
		// 是否点赞
		archiveAuth.GET("hasLike", api.HasLike)

		// 收藏
		archiveAuth.POST("collect", api.Collect)
		// 已收藏的文件夹
		archiveAuth.GET("getCollectInfo", api.GetCollectedInfo)
		// 是否收藏
		archiveAuth.GET("hasCollect", api.HasCollect)

	}

	// 获取点赞收藏数据
	archiveGroup.GET("stat", api.GetArchiveStat)
}
