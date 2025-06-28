package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

func AddHistory(ctx *gin.Context) {
	var historyReq dto.HistoryReq
	if err := ctx.Bind(&historyReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.AddHistory(ctx, historyReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

func GetHistoryList(ctx *gin.Context) {
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	videos, err := service.GetHistoryList(ctx, page, pageSize)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"videos": videos})
}

func GetHistoryProgress(ctx *gin.Context) {
	videoId := utils.StringToUint(ctx.Query("vid"))
	part := utils.StringToUint(ctx.Query("part"))

	progress, realPart, err := service.GetHistoryProgress(ctx, videoId, part)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"progress": progress, "part": realPart})
}
