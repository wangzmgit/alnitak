package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 审核通过
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

// 审核不通过
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

// 获取审核记录
func GetReviewRecord(ctx *gin.Context) {
	id := utils.StringToUint(ctx.Query("vid"))

	review, err := service.GetReviewRecord(ctx, id)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"review": review})
}
