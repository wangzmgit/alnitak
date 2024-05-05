package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectCommentRoutes(r *gin.RouterGroup) {
	commentGroup := r.Group("comment")

	commentAuth := commentGroup.Group("")
	commentAuth.Use(middleware.Auth())
	{
		// 发表评论或回复
		commentAuth.POST("addVideoComment", api.AddVideoComment)
		// 删除评论或回复
		commentAuth.DELETE("deleteVideoComment/:id", api.DeleteVideoComment)
		// 获取评论列表
		commentAuth.GET("getVideoCommentList", api.GetVideoCommentList)
		// 发表文章评论或回复
		commentAuth.POST("addArticleComment", api.AddArticleComment)
		// 删除评论或回复
		commentAuth.DELETE("deleteArticleComment/:id", api.DeleteArticleComment)
		// 获取评论列表
		commentAuth.GET("getArticleCommentList", api.GetArticleCommentList)
	}

	// 获取视频评论
	commentGroup.GET("getVideoComment", api.GetVideoComment)
	// 获取视频回复
	commentGroup.GET("getVideoReply", api.GetVideoReply)
	// 获取文章评论
	commentGroup.GET("getArticleComment", api.GetArticleComment)
	// 获取文章回复
	commentGroup.GET("getArticleReply", api.GetArticleReply)
}
