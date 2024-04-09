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
		commentAuth.POST("/addComment", api.AddComment)
		// 删除评论或回复
		commentAuth.DELETE("/deleteComment/:id", api.DeleteComment)
		// 获取评论列表
		commentAuth.GET("getCommentList", api.GetCommentList)
	}

	// 获取评论
	commentGroup.GET("getComment", api.GetComment)
	// 获取回复
	commentGroup.GET("getReply", api.GetReply)
}
