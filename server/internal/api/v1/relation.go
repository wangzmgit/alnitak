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
func Unfollow(ctx *gin.Context) {
	var idReq dto.IdReq
	if err := ctx.Bind(&idReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.Unfollow(ctx, idReq.ID); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 获取用户关系
func GetUserRelation(ctx *gin.Context) {
	targetUid := utils.StringToUint(ctx.Query("userId"))

	relation := service.GetUserRelation(ctx, targetUid)

	// 返回
	resp.OkWithData(ctx, gin.H{"relation": relation})
}

// 获取关注列表
func GetFollowings(ctx *gin.Context) {
	userId := utils.StringToUint(ctx.Query("userId"))
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 50 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	users, err := service.GetFollowings(ctx, userId, page, pageSize)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"users": users})
}

// 获取粉丝列表
func GetFollowers(ctx *gin.Context) {
	userId := utils.StringToUint(ctx.Query("userId"))
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 50 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	users, err := service.GetFollowers(ctx, userId, page, pageSize)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"users": users})
}

// 获取关注和粉丝数
func GetFollowCount(ctx *gin.Context) {
	userId := utils.StringToUint(ctx.Query("userId"))

	following, follower := service.GetFollowCount(ctx, userId)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"following": following, "follower": follower})
}
