package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

type MenuApi struct{}

// 添加菜单
func (m *MenuApi) AddMenu(ctx *gin.Context) {
	// 获取参数
	var addMenuReq dto.AddMenuReq
	if err := ctx.Bind(&addMenuReq); err != nil {
		service.AddFailOperation(ctx, "添加菜单", "请求参数有误", nil, nil)
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(addMenuReq.Name, "=", 0) { // 路由Name不能为空
		service.AddFailOperation(ctx, "添加菜单", "路由Name为空", nil, nil)
		resp.FailWithMessage(ctx, "路由Name不能为空")
		return
	}

	if utils.VerifyStringLength(addMenuReq.Path, "=", 0) { // 路由Path不能为空
		service.AddFailOperation(ctx, "添加菜单", "路由Path为空", nil, nil)
		resp.FailWithMessage(ctx, "路由Path不能为空")
		return
	}

	if utils.VerifyStringLength(addMenuReq.Title, "=", 0) { // 菜单标题不能为空
		service.AddFailOperation(ctx, "添加菜单", "菜单标题为空", nil, nil)
		resp.FailWithMessage(ctx, "菜单标题不能为空")
		return
	}

	if utils.VerifyStringLength(addMenuReq.Icon, "=", 0) { // 菜单图标不能为空
		service.AddFailOperation(ctx, "添加菜单", " 菜单图标为空", nil, nil)
		resp.FailWithMessage(ctx, "菜单图标不能为空")
		return
	}

	if err := service.AddMenu(ctx, addMenuReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 编辑菜单
func (m *MenuApi) EditMenu(ctx *gin.Context) {
	// 获取参数
	var editMenuReq dto.EditMenuReq
	if err := ctx.Bind(&editMenuReq); err != nil {
		service.AddFailOperation(ctx, "编辑菜单", "请求参数有误", nil, nil)
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(editMenuReq.Name, "=", 0) { // 路由Name不能为空
		service.AddFailOperation(ctx, "编辑菜单", "路由Name为空", nil, nil)
		resp.FailWithMessage(ctx, "路由Name不能为空")
		return
	}

	if utils.VerifyStringLength(editMenuReq.Path, "=", 0) { // 路由Path不能为空
		service.AddFailOperation(ctx, "编辑菜单", "路由Path为空", nil, nil)
		resp.FailWithMessage(ctx, "路由Path不能为空")
		return
	}

	if utils.VerifyStringLength(editMenuReq.Title, "=", 0) { // 菜单标题不能为空
		service.AddFailOperation(ctx, "编辑菜单", "菜单标题为空", nil, nil)
		resp.FailWithMessage(ctx, "菜单标题不能为空")
		return
	}

	if utils.VerifyStringLength(editMenuReq.Icon, "=", 0) { // 菜单图标不能为空
		service.AddFailOperation(ctx, "编辑菜单", " 菜单图标为空", nil, nil)
		resp.FailWithMessage(ctx, "菜单图标不能为空")
		return
	}

	if err := service.EditMenu(ctx, editMenuReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 获取菜单树
func (m *MenuApi) GetMenuTree(ctx *gin.Context) {
	menus := service.GetMenuTree()

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"menus": vo.MenuListToMenuResp(menus)})
}

// 获取用户菜单
func (m *MenuApi) GetUserMenuTree(ctx *gin.Context) {
	roleCode := ctx.GetString("roleCode")
	menus := service.GetMenuTreeByRoleCode(roleCode)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"menus": vo.MenuListToMenuTree(menus, 0)})
}

// 删除菜单
func (m *MenuApi) DeleteMenu(ctx *gin.Context) {
	// 获取参数
	id := utils.StringToUint(ctx.Param("id"))

	if err := service.DeleteMenu(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 获取角色菜单列表
func (m *MenuApi) GetRoleMenu(ctx *gin.Context) {
	roleCode := ctx.Query("code")
	menus := service.GetMenuTreeByRoleCode(roleCode)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"menus": vo.MenuListToMenuResp(menus)})
}

// 编辑角色菜单
func (m *MenuApi) EditRoleMenu(ctx *gin.Context) {
	var editRoleMenuReq dto.EditRoleMenuReq
	if err := ctx.Bind(&editRoleMenuReq); err != nil {
		service.AddFailOperation(ctx, "编辑角色菜单", "角色代码为空", nil, nil)
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.EditRoleMenu(ctx, editRoleMenuReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}
