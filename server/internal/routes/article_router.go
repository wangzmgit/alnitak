package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectArticleRoutes(r *gin.RouterGroup) {
	articleGroup := r.Group("article")

	articleAuth := articleGroup.Group("")
	articleAuth.Use(middleware.Auth())
	{
		// 上传文章信息
		articleAuth.POST("uploadArticleInfo", api.UploadArticleInfo)
		// 获取文章状态信息
		articleAuth.GET("getArticleStatus", api.GetArticleStatus)
		// 获取上传的文章
		articleAuth.GET("getUploadArticle", api.GetUploadArticleList)
		// 编辑文章信息
		articleAuth.PUT("editArticleInfo", api.EditArticleInfo)
		// 删除文章
		articleAuth.DELETE("deleteArticle/:id", api.DeleteArticle)
		// 获取所有的文章列表
		articleAuth.GET("getAllArticleList", api.GetAllArticleList)
		// 获取审核列表（后台管理）
		articleAuth.POST("getReviewArticleList", api.GetReviewArticleList)
		// 获取文章列表（后台管理）
		articleAuth.POST("getArticleListManage", api.GetArticleListManage)
		// 删除文章（后台管理）
		articleAuth.DELETE("deleteArticleManage/:id", api.DeleteArticleManage)
	}

	// 获取文章信息
	articleGroup.GET("getArticleById", api.GetArticleById)
	// 获取用户文章
	articleGroup.GET("getArticleByUser", api.GetArticleByUser)
	// 获取分区视频
	articleGroup.GET("getRandomArticleList", api.GetRandomArticleList)
}
