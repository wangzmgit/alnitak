package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 新增角色
func AddRole(ctx *gin.Context) {
	// 获取参数
	var addRoleReq dto.AddRoleReq
	if err := ctx.Bind(&addRoleReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(addRoleReq.Name, "=", 0) { // 角色名格式验证
		resp.FailWithMessage(ctx, "角色名不能为空")
		return
	}

	if utils.VerifyStringLength(addRoleReq.Code, "=", 0) { // 角色代码验证
		resp.FailWithMessage(ctx, "角色代码不能为空")
		return
	}

	if err := service.AddRole(ctx, addRoleReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 获取角色列表
func GetRoleList(ctx *gin.Context) {
	// 获取参数
	var roleListReq dto.RoleListReq
	if err := ctx.Bind(&roleListReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if roleListReq.PageSize > 100 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, roles := service.GetRoleList(roleListReq.Page, roleListReq.PageSize)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"list": roles, "total": total})
}

// 获取全部角色
func GetAllRoleList(ctx *gin.Context) {
	roles := service.GetAllRoleList(ctx)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"roles": roles})
}

// 编辑角色
func EditRole(ctx *gin.Context) {
	// 获取参数
	var editRoleReq dto.EditRoleReq
	if err := ctx.Bind(&editRoleReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(editRoleReq.Name, "=", 0) { // 角色名格式验证
		resp.FailWithMessage(ctx, "角色名不能为空")
		return
	}

	// 更新数据库
	if err := service.EditRole(ctx, editRoleReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 删除角色
func DeleteRole(ctx *gin.Context) {
	// 获取参数
	id := utils.StringToUint(ctx.Param("id"))

	if err := service.DeleteRole(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 编辑角色首页
func EditRoleHome(ctx *gin.Context) {
	// 获取参数
	var editRoleHomeReq dto.EditRoleHomeReq
	if err := ctx.Bind(&editRoleHomeReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.EditRoleHome(ctx, editRoleHomeReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 获取个人角色信息
func GetRoleInfo(ctx *gin.Context) {
	roleCode := ctx.GetString("roleCode")
	role := service.GetRoleInfo(roleCode)

	resp.OkWithData(ctx, gin.H{"role": role})
}
