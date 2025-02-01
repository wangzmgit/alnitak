package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectReviewRoutes(r *gin.RouterGroup) {
	reviewGroup := r.Group("review")

	reviewAuth := reviewGroup.Group("")
	reviewAuth.Use(middleware.Auth())
	{
		// 视频审核通过
		reviewAuth.POST("reviewVideoApproved", api.ReviewVideoApproved)
		// 视频审核不通过
		reviewAuth.POST("reviewVideoFailed", api.ReviewVideoFailed)
		// 获取视频审核记录
		reviewAuth.GET("getVideoReviewRecord", api.GetVideoReviewRecord)

		// 文章审核通过
		reviewAuth.POST("reviewArticleApproved", api.ReviewArticleApproved)
		// 文章审核不通过
		reviewAuth.POST("reviewArticleFailed", api.ReviewArticleFailed)
		// 获取文章审核记录
		reviewAuth.GET("getArticleReviewRecord", api.GetArticleReviewRecord)
	}
}
