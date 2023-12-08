package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

type RoleApi struct{}

// 新增角色
func (r *RoleApi) AddRole(ctx *gin.Context) {
	// 获取参数
	var addRoleReq dto.AddRoleReq
	if err := ctx.Bind(&addRoleReq); err != nil {
		service.AddFailOperation(ctx, "新增角色", "请求参数有误", nil, nil)
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(addRoleReq.Name, "=", 0) { // 角色名格式验证
		service.AddFailOperation(ctx, "新增角色", "角色名为空", nil, nil)
		resp.FailWithMessage(ctx, "角色名不能为空")
		return
	}

	if utils.VerifyStringLength(addRoleReq.Code, "=", 0) { // 角色代码验证
		service.AddFailOperation(ctx, "新增角色", "角色代码为空", nil, nil)
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
func (r *RoleApi) GetRoleList(ctx *gin.Context) {
	// 获取参数
	var roleListReq dto.RoleListReq
	if err := ctx.Bind(&roleListReq); err != nil {
		service.AddFailOperation(ctx, "获取角色列表", "请求参数有误", nil, nil)
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if roleListReq.PageSize > 100 {
		service.AddFailOperation(ctx, "获取角色列表", "请求数量过多", gin.H{"数量": roleListReq.PageSize}, nil)
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, roles := service.GetRoleList(roleListReq.Page, roleListReq.PageSize)

	// 返回给前端
	service.AddOkOperation(ctx, "获取角色列表", "获取成功", nil)
	resp.OkWithData(ctx, gin.H{"list": roles, "total": total})
}

// 编辑角色
func (r *RoleApi) EditRole(ctx *gin.Context) {
	// 获取参数
	var editRoleReq dto.EditRoleReq
	if err := ctx.Bind(&editRoleReq); err != nil {
		service.AddFailOperation(ctx, "编辑角色", "请求参数有误", nil, nil)
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(editRoleReq.Name, "=", 0) { // 角色名格式验证
		service.AddFailOperation(ctx, "编辑角色", "角色名为空", nil, nil)
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
func (r *RoleApi) DeleteRole(ctx *gin.Context) {
	// 获取参数
	var idReq dto.IdReq
	if err := ctx.Bind(&idReq); err != nil {
		service.AddFailOperation(ctx, "删除角色", "请求参数有误", nil, nil)
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.DeleteRole(ctx, idReq.ID); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 编辑角色首页
func (r *RoleApi) EditRoleHome(ctx *gin.Context) {
	// 获取参数
	var editRoleHomeReq dto.EditRoleHomeReq
	if err := ctx.Bind(&editRoleHomeReq); err != nil {
		service.AddFailOperation(ctx, "编辑角色首页", "请求参数有误", nil, nil)
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
func (r *RoleApi) GetRoleInfo(ctx *gin.Context) {
	roleCode := ctx.GetString("roleCode")
	role := service.GetRoleInfo(roleCode)

	resp.OkWithData(ctx, gin.H{"role": role})
}
