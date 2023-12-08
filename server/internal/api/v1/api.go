package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

type ApiApi struct{}

// 获取API列表
func (a *ApiApi) GetApiList(ctx *gin.Context) {
	// 获取参数
	var apiListReq dto.ApiListReq
	if err := ctx.Bind(&apiListReq); err != nil {
		service.AddFailOperation(ctx, "获取API列表", "请求参数有误", nil, nil)
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if apiListReq.PageSize > 100 {
		service.AddFailOperation(ctx, "获取API列表", "请求数量过多", gin.H{"数量": apiListReq.PageSize}, nil)
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, apis := service.GetApiList(apiListReq.Page, apiListReq.PageSize)

	// 返回给前端
	service.AddOkOperation(ctx, "获取API列表", "获取成功", nil)
	resp.OkWithData(ctx, gin.H{"list": apis, "total": total})
}

// 获取全部API列表
func (a *ApiApi) GetAllApiList(ctx *gin.Context) {
	total, apis := service.GetAllApiList()

	// 返回给前端
	service.AddOkOperation(ctx, "获取全部API列表", "获取成功", nil)
	resp.OkWithData(ctx, gin.H{"list": apis, "total": total})
}

// 添加API
func (a *ApiApi) AddApi(ctx *gin.Context) {
	// 获取参数
	var addApiReq dto.AddApiReq
	if err := ctx.Bind(&addApiReq); err != nil {
		service.AddFailOperation(ctx, "添加API", "请求参数有误", nil, nil)
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(addApiReq.Path, "=", 0) { // 路径不能为空
		service.AddFailOperation(ctx, "添加API", "API路径为空", nil, nil)
		resp.FailWithMessage(ctx, "路径不能为空")
		return
	}

	if utils.VerifyStringLength(addApiReq.Method, "=", 0) { // 请求方法不能为空
		service.AddFailOperation(ctx, "添加API", "请求方法为空", nil, nil)
		resp.FailWithMessage(ctx, "请求方法不能为空")
		return
	}

	if utils.VerifyStringLength(addApiReq.Category, "=", 0) { // 分类名不能为空
		service.AddFailOperation(ctx, "添加API", "分类名为空", nil, nil)
		resp.FailWithMessage(ctx, "分类名不能为空")
		return
	}

	if err := service.AddApi(ctx, addApiReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 编辑API
func (a *ApiApi) EditApi(ctx *gin.Context) {
	// 获取参数
	var editApiReq dto.EditApiReq
	if err := ctx.Bind(&editApiReq); err != nil {
		service.AddFailOperation(ctx, "编辑API", "请求参数有误", nil, nil)
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(editApiReq.Path, "=", 0) { // 路径不能为空
		service.AddFailOperation(ctx, "编辑API", "API路径为空", nil, nil)
		resp.FailWithMessage(ctx, "路径不能为空")
		return
	}

	if utils.VerifyStringLength(editApiReq.Method, "=", 0) { // 请求方法不能为空
		service.AddFailOperation(ctx, "编辑API", "请求方法为空", nil, nil)
		resp.FailWithMessage(ctx, "请求方法不能为空")
		return
	}

	if utils.VerifyStringLength(editApiReq.Category, "=", 0) { // 分类名不能为空
		service.AddFailOperation(ctx, "编辑API", "分类名为空", nil, nil)
		resp.FailWithMessage(ctx, "分类名不能为空")
		return
	}

	if err := service.EditApi(ctx, editApiReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 删除APi
func (a *ApiApi) DeleteApi(ctx *gin.Context) {
	// 获取参数
	id := utils.StringToUint(ctx.Param("id"))

	if err := service.DeleteApi(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 获取角色API列表
func (a *ApiApi) GetRoleApi(ctx *gin.Context) {
	roleCode := ctx.Query("code")

	// 参数校验
	if utils.VerifyStringLength(roleCode, "=", 0) { // 角色代码为空
		service.AddFailOperation(ctx, "获取角色API", "角色代码为空", nil, nil)
		resp.FailWithMessage(ctx, "角色代码不能为空")
		return
	}

	// 获取角色菜单
	apis := service.SelectRoleApiList(roleCode)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"list": apis})
}

// 编辑角色Api
func (a *ApiApi) EditRoleApi(ctx *gin.Context) {
	var editRoleApiReq dto.EditRoleApiReq
	if err := ctx.Bind(&editRoleApiReq); err != nil {
		service.AddFailOperation(ctx, "编辑角色Api", "角色代码为空", nil, nil)
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.EditRoleApi(ctx, editRoleApiReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}
