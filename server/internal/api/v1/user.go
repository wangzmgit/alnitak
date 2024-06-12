package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

func GetUserInfo(ctx *gin.Context) {
	userId := ctx.GetUint("userId")
	user := service.GetUserInfo(userId)

	resp.OkWithData(ctx, gin.H{"userInfo": user})
}

func EditUserInfo(ctx *gin.Context) {
	var editUserInfoReq dto.EditUserInfoReq
	if err := ctx.Bind(&editUserInfoReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(editUserInfoReq.Name, "=", 0) {
		resp.FailWithMessage(ctx, "用户名不能为空")
		return
	}

	if err := service.EditUserInfo(ctx, editUserInfoReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
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

// 获取用户列表(后台管理)
func GetUserListManage(ctx *gin.Context) {
	// 获取参数
	var userListReq dto.UserListReq
	if err := ctx.Bind(&userListReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if userListReq.PageSize > 100 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, roles := service.GetUserListManage(userListReq)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"list": roles, "total": total})
}

// 编辑用户信息(后台管理)
func EditUserInfoManage(ctx *gin.Context) {
	var editUserInfoManageReq dto.EditUserInfoManageReq
	if err := ctx.Bind(&editUserInfoManageReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(editUserInfoManageReq.Name, "=", 0) {
		resp.FailWithMessage(ctx, "用户名不能为空")
		return
	}

	if err := service.EditUserInfoManage(ctx, editUserInfoManageReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

// 设置用户角色
func EditUserRole(ctx *gin.Context) {
	// 获取参数
	var editUserRoleReq dto.EditUserRoleReq
	if err := ctx.Bind(&editUserRoleReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.EditUserRole(ctx, editUserRoleReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 删除用户
func DeleteUser(ctx *gin.Context) {
	// 获取参数
	id := utils.StringToUint(ctx.Param("id"))

	if err := service.DeleteUser(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}
