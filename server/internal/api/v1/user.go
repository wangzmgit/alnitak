package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

func GetUserInfo(ctx *gin.Context) {
	userId := ctx.GetUint("userId")
	user := service.GetUserInfo(userId)

	var ban *vo.BanResp
	if user.Status == 1 {
		b := service.FindUserLastBan(userId)
		ban = &b
	}

	resp.OkWithData(ctx, gin.H{"userInfo": user, "ban": ban})
}

func EditUserInfo(ctx *gin.Context) {
	var editUserInfoReq dto.EditUserInfoReq
	if err := ctx.Bind(&editUserInfoReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if utils.VerifyStringLength(editUserInfoReq.Name, "=", 0) {
		resp.FailWithMessage(ctx, "用户名不能为空")
		return
	}

	if err := service.EditUserInfo(ctx, editUserInfoReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 获取用户基本信息
func GetUserBaseInfo(ctx *gin.Context) {
	userId := utils.StringToUint(ctx.Query("userId"))

	user := service.GetUserBaseInfo(userId)
	if user.ID == 0 {
		resp.FailWithMessage(ctx, "用户不存在")
		return
	}

	resp.OkWithData(ctx, gin.H{"userInfo": user})
}

// SearchUser 搜索用户 / UP（公开）
func SearchUser(ctx *gin.Context) {
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
	users := service.SearchUser(ctx, req)
	resp.OkWithData(ctx, gin.H{"users": users})
}

// 获取用户列表(后台管理)
func GetUserListManage(ctx *gin.Context) {

	var userListReq dto.UserListReq
	if err := ctx.Bind(&userListReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if userListReq.PageSize > 50 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, roles := service.GetUserListManage(userListReq)

	resp.OkWithData(ctx, gin.H{"list": roles, "total": total})
}

// 编辑用户信息(后台管理)
func EditUserInfoManage(ctx *gin.Context) {
	var editUserInfoManageReq dto.EditUserInfoManageReq
	if err := ctx.Bind(&editUserInfoManageReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if utils.VerifyStringLength(editUserInfoManageReq.Name, "=", 0) {
		resp.FailWithMessage(ctx, "用户名不能为空")
		return
	}

	if err := service.EditUserInfoManage(ctx, editUserInfoManageReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 设置用户角色
func EditUserRole(ctx *gin.Context) {
	var editUserRoleReq dto.EditUserRoleReq
	if err := ctx.Bind(&editUserRoleReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.EditUserRole(ctx, editUserRoleReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 封禁用户
func BanUser(ctx *gin.Context) {
	var banUserReq dto.BanUserReq
	if err := ctx.Bind(&banUserReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.BanUser(ctx, banUserReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 解封用户
func UnBanUser(ctx *gin.Context) {
	var idReq dto.IdReq
	if err := ctx.Bind(&idReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.UnBanUser(ctx, idReq.ID); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 获取封禁记录
func GetUserBanRecord(ctx *gin.Context) {
	userId := utils.StringToUint(ctx.Query("uid"))
	userBanList := service.GetUserBanRecord(userId)

	resp.OkWithData(ctx, gin.H{"list": userBanList})
}

// 删除用户
func DeleteUser(ctx *gin.Context) {
	id := utils.StringToUint(ctx.Param("id"))

	if err := service.DeleteUser(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}
