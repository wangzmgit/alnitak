package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 关注
func Follow(ctx *gin.Context) {
	var idReq dto.IdReq
	if err := ctx.Bind(&idReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.Follow(ctx, idReq.ID); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 取消关注
func UnFollow(ctx *gin.Context) {
	var idReq dto.IdReq
	if err := ctx.Bind(&idReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.UnFollow(ctx, idReq.ID); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 获取用户关系
func GetUserRelation(ctx *gin.Context) {
	targetUid := utils.StringToUint(ctx.Query("uid"))

	relation := service.GetUserRelation(ctx, targetUid)

	// 返回
	resp.OkWithData(ctx, gin.H{"relation": relation})
}
