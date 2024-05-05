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
		reviewAuth.POST("reviewVideoApproved", api.ReviewVideoApproved)
		reviewAuth.POST("reviewVideoFailed", api.ReviewVideoFailed)
		reviewAuth.GET("getVideoReviewRecord", api.GetVideoReviewRecord)

		reviewAuth.POST("reviewArticleApproved", api.ReviewArticleApproved)
		reviewAuth.POST("reviewArticleFailed", api.ReviewArticleFailed)
		reviewAuth.GET("getArticleReviewRecord", api.GetArticleReviewRecord)
	}
}
