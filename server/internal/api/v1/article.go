package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 上传文章信息
func UploadArticleInfo(ctx *gin.Context) {
	var uploadArticleReq dto.UploadArticleReq
	if err := ctx.Bind(&uploadArticleReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(uploadArticleReq.Title, "=", 0) {
		resp.FailWithMessage(ctx, "标题不能为空")
		return
	}

	if err := service.UploadArticleInfo(ctx, uploadArticleReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

// 编辑文章信息
func EditArticleInfo(ctx *gin.Context) {
	var editArticleReq dto.EditArticleReq
	if err := ctx.Bind(&editArticleReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(editArticleReq.Title, "=", 0) {
		resp.FailWithMessage(ctx, "标题不能为空")
		return
	}

	if err := service.EditArticleInfo(ctx, editArticleReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

// 获取自己的文章
func GetUploadArticleList(ctx *gin.Context) {
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, articles := service.GetUploadArticleList(ctx, page, pageSize)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"total": total, "articles": articles})
}

// 获取视频状态
func GetArticleStatus(ctx *gin.Context) {
	articleId := utils.StringToUint(ctx.Query("aid"))

	article, err := service.GetArticleStatus(ctx, articleId)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"article": article})
}

// 获取文章信息
func GetArticleById(ctx *gin.Context) {
	aid := utils.StringToUint(ctx.DefaultQuery("aid", "0"))

	article, err := service.GetArticleById(ctx, aid)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"article": article})
}

// 删除文章
func DeleteArticle(ctx *gin.Context) {
	// 获取参数
	id := utils.StringToUint(ctx.Param("id"))

	if err := service.DeleteArticle(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 获取所有的视频列表
func GetAllArticleList(ctx *gin.Context) {
	articles := service.GetAllArticleList(ctx)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"articles": articles})
}

// 获取用户视频
func GetArticleByUser(ctx *gin.Context) {
	userId := utils.StringToUint(ctx.Query("userId"))
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, articles := service.GetArticleByUser(ctx, userId, page, pageSize)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"total": total, "articles": articles})
}

// 获取分区视频
func GetRandomArticleList(ctx *gin.Context) {
	size := utils.StringToInt(ctx.Query("size"))

	if size > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	articles := service.GetRandomArticleList(ctx, size)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"articles": articles})
}

// 获取待审核文章列表
func GetReviewArticleList(ctx *gin.Context) {
	// 获取参数
	var reviewListReq dto.ReviewArticleListReq
	if err := ctx.Bind(&reviewListReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if reviewListReq.PageSize > 100 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, articles := service.GetReviewArticleList(reviewListReq)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"list": articles, "total": total})
}

// 获取文章列表(后台管理)
func GetArticleListManage(ctx *gin.Context) {
	// 获取参数
	var articleListReq dto.ArticleListReq
	if err := ctx.Bind(&articleListReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if articleListReq.PageSize > 100 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, articles := service.GetArticleListManage(articleListReq)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"list": articles, "total": total})
}

// 删除文章(后台管理)
func DeleteArticleManage(ctx *gin.Context) {
	// 获取参数
	id := utils.StringToUint(ctx.Param("id"))

	if err := service.DeleteArticleManage(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}
