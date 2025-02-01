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
		// 视频点赞
		archiveAuth.POST("video/like", api.LikeVideo)
		// 视频取消赞
		archiveAuth.POST("video/cancelLike", api.CancelLikeVideo)
		// 视频是否点赞
		archiveAuth.GET("video/hasLike", api.HasLikeVideo)

		// 文章点赞
		archiveAuth.POST("article/like", api.LikeArticle)
		// 文章取消赞
		archiveAuth.POST("article/cancelLike", api.CancelLikeArticle)
		// 文章是否点赞
		archiveAuth.GET("article/hasLike", api.HasLikeArticle)

		// 视频收藏
		archiveAuth.POST("video/collect", api.CollectVideo)
		// 获取视频已收藏的文件夹
		archiveAuth.GET("video/getCollectInfo", api.GetVideoCollectedInfo)
		// 视频是否收藏
		archiveAuth.GET("video/hasCollect", api.HasCollectVideo)

		// 文章收藏
		archiveAuth.POST("article/collect", api.CollectArticle)
		// 文章取消收藏
		archiveAuth.POST("article/cancelCollect", api.CancelCollectArticle)
		// 文章是否收藏
		archiveAuth.GET("article/hasCollect", api.HasCollectArticle)
	}

	// 获取视频点赞收藏数据
	archiveGroup.GET("video/stat", api.GetVideoArchiveStat)
	// 获取文章点赞收藏数据
	archiveGroup.GET("article/stat", api.GetArticleArchiveStat)
}
