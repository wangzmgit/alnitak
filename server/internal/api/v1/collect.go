package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 添加收藏
func Collect(ctx *gin.Context) {
	// 获取参数
	var addCollectReq dto.AddCollectReq
	if err := ctx.Bind(&addCollectReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.Collect(ctx, addCollectReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 是否收藏
func HasCollect(ctx *gin.Context) {
	videoId := utils.StringToUint(ctx.Query("vid"))
	collect, err := service.HasCollect(ctx, videoId)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"collect": collect})
}

// 获取已收藏的收藏夹
func GetCollectedInfo(ctx *gin.Context) {
	videoId := utils.StringToUint(ctx.Query("vid"))

	collectionIds := service.GetCollectedInfo(ctx, videoId)

	// 返回
	resp.OkWithData(ctx, gin.H{"collectionIds": collectionIds})
}

// 获取收藏视频列表
func GetCollectVideo(ctx *gin.Context) {
	collectionId := utils.StringToUint(ctx.Query("cid"))
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, videos, err := service.GetCollectVideo(ctx, collectionId, page, pageSize)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"videos": videos, "total": total})
}
