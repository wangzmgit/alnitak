package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

func LikeVideo(ctx *gin.Context) {
	// 获取参数
	var likeReq dto.LikeVideoReq
	if err := ctx.Bind(&likeReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.LikeVideo(ctx, likeReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

func CancelLikeVideo(ctx *gin.Context) {
	var likeReq dto.LikeVideoReq
	if err := ctx.Bind(&likeReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.CancelLikeVideo(ctx, likeReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

func HasLikeVideo(ctx *gin.Context) {
	videoId := utils.StringToUint(ctx.Query("vid"))
	like, err := service.HasLikeVideo(ctx, videoId)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"like": like})
}
