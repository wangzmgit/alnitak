package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 审核通过(视频)
func ReviewVideoApproved(ctx *gin.Context) {
	var reviewVideoReq dto.ReviewVideoReq
	if err := ctx.Bind(&reviewVideoReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.ReviewVideoApproved(ctx, reviewVideoReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

// 审核不通过(视频)
func ReviewVideoFailed(ctx *gin.Context) {
	var reviewVideoReq dto.ReviewVideoReq
	if err := ctx.Bind(&reviewVideoReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.ReviewVideoFailed(ctx, reviewVideoReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

// 获取审核记录(视频)
func GetVideoReviewRecord(ctx *gin.Context) {
	id := utils.StringToUint(ctx.Query("vid"))

	review, err := service.GetVideoReviewRecord(ctx, id)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"review": review})
}

// 审核通过(文章)
func ReviewArticleApproved(ctx *gin.Context) {
	var reviewArticleReq dto.ReviewArticleReq
	if err := ctx.Bind(&reviewArticleReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.ReviewArticleApproved(ctx, reviewArticleReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

// 审核不通过(文章)
func ReviewArticleFailed(ctx *gin.Context) {
	var reviewArticleReq dto.ReviewArticleReq
	if err := ctx.Bind(&reviewArticleReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.ReviewArticleFailed(ctx, reviewArticleReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

// 获取审核记录(文章)
func GetArticleReviewRecord(ctx *gin.Context) {
	id := utils.StringToUint(ctx.Query("aid"))

	review, err := service.GetArticleReviewRecord(ctx, id)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"review": review})
}
