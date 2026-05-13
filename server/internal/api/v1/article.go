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
// 支持按分类筛选：category=all|published|pending|rejected（文章无转码状态）
func GetUploadArticleList(ctx *gin.Context) {
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))
	category := ctx.DefaultQuery("category", "all")

	if pageSize > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, articles := service.GetUploadArticleList(ctx, page, pageSize, category)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"total": total, "articles": articles})
}

// 获取文章状态
func GetArticleStatus(ctx *gin.Context) {
	raw := ctx.Query("aid")
	articleId, err := service.ParseArticleID(raw)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	if articleId == 0 {
		resp.FailWithMessage(ctx, "参数有误")
		return
	}

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
	raw := ctx.DefaultQuery("aid", "")
	aid, err := service.ParseArticleID(raw)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	if aid == 0 {
		resp.FailWithMessage(ctx, "参数有误")
		return
	}

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

// 获取所有的文章列表
func GetAllArticleList(ctx *gin.Context) {
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 30 {
		pageSize = 10
	}

	total, articles := service.GetAllArticleList(ctx, page, pageSize)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"total": total, "articles": articles})
}

// 获取用户文章
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

// 获取分区文章
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

// SearchArticle 搜索专栏（公开）
func SearchArticle(ctx *gin.Context) {
	var req dto.SearchKeywordPageReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}
	if req.PageSize > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 15
	}
	articles := service.SearchArticle(ctx, req)
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
