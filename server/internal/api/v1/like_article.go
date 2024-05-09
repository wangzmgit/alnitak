package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

func LikeArticle(ctx *gin.Context) {
	// 获取参数
	var likeReq dto.LikeArticleReq
	if err := ctx.Bind(&likeReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.LikeArticle(ctx, likeReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

func CancelLikeArticle(ctx *gin.Context) {
	var likeReq dto.LikeArticleReq
	if err := ctx.Bind(&likeReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.CancelLikeArticle(ctx, likeReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

func HasLikeArticle(ctx *gin.Context) {
	videoId := utils.StringToUint(ctx.Query("vid"))
	like, err := service.HasLikeArticle(ctx, videoId)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"like": like})
}
