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
		commentAuth.POST("video/addComment", middleware.Ban(), middleware.RateLimiter(middleware.CommentRateLimit), api.AddVideoComment)
		// 删除评论或回复
		commentAuth.DELETE("video/deleteComment/:id", api.DeleteVideoComment)
		// 获取评论列表
		commentAuth.GET("video/getCommentList", api.GetVideoCommentList)
		// 点赞评论
		commentAuth.POST("video/like/:id", api.LikeComment)
		// 取消点赞评论
		commentAuth.DELETE("video/like/:id", api.UnlikeComment)
		// 点踩评论
		commentAuth.POST("video/dislike/:id", api.DislikeComment)
		// 取消点踩评论
		commentAuth.DELETE("video/dislike/:id", api.UndislikeComment)
		// 发表文章评论或回复
		commentAuth.POST("article/addComment", middleware.Ban(), middleware.RateLimiter(middleware.CommentRateLimit), api.AddArticleComment)
		// 删除评论或回复
		commentAuth.DELETE("article/deleteComment/:id", api.DeleteArticleComment)
		// 获取评论列表
		commentAuth.GET("article/getCommentList", api.GetArticleCommentList)
		// 点赞评论（文章）
		commentAuth.POST("article/like/:id", api.LikeComment)
		// 取消点赞评论（文章）
		commentAuth.DELETE("article/like/:id", api.UnlikeComment)
		// 点踩评论（文章）
		commentAuth.POST("article/dislike/:id", api.DislikeComment)
		// 取消点踩评论（文章）
		commentAuth.DELETE("article/dislike/:id", api.UndislikeComment)
	}

	// 公共只读接口：允许匿名访问，登录时解析 userId 以便回填 liked
	commentPublic := commentGroup.Group("")
	commentPublic.Use(middleware.OptionalAuth())
	{
		commentPublic.GET("video/getComment", api.GetVideoComment)
		commentPublic.GET("video/getReply", api.GetVideoReply)
		commentPublic.GET("article/getComment", api.GetArticleComment)
		commentPublic.GET("article/getReply", api.GetArticleReply)
	}
}
