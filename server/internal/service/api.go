package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
)

func GetApiList(page, pageSize int) (total int64, apis []vo.ApiResp) {
	global.Mysql.Model(&model.Api{}).Count(&total)
	global.Mysql.Model(&model.Api{}).Limit(pageSize).Offset((page - 1) * pageSize).Scan(&apis)
	return
}

func GetAllApiList() (total int64, apis []vo.ApiResp) {
	global.Mysql.Model(&model.Api{}).Count(&total)
	global.Mysql.Model(&model.Api{}).Scan(&apis)
	return
}

// 新增API
func AddApi(ctx *gin.Context, addApiReq dto.AddApiReq) error {
	// 保存到数据库
	if err := global.Mysql.Create(&model.Api{
		Path:     addApiReq.Path,
		Method:   addApiReq.Method,
		Category: addApiReq.Category,
		Desc:     addApiReq.Desc,
	}).Error; err != nil {
		AddFailOperation(ctx, "新增API", "更新数据库失败", nil, err)
		return errors.New("新增API失败")
	}

	return nil
}

// 编辑API
func EditApi(ctx *gin.Context, editApiReq dto.EditApiReq) error {
	if err := global.Mysql.Model(&model.Api{}).Where("id = ?", editApiReq.ID).Updates(
		map[string]interface{}{
			"path":     editApiReq.Path,
			"method":   editApiReq.Method,
			"category": editApiReq.Category,
			"desc":     editApiReq.Desc,
		},
	).Error; err != nil {
		AddFailOperation(ctx, "编辑API", "更新数据库失败", nil, err)
		return errors.New("编辑API失败")
	}
	return nil
}

// 删除API
func DeleteApi(ctx *gin.Context, id uint) error {
	if err := global.Mysql.Where("id = ?", id).Delete(&model.Api{}).Error; err != nil {
		AddFailOperation(ctx, "删除API", "更新数据库失败", nil, err)
		return errors.New("删除API失败")
	}

	return nil
}

// 获取角色API列表
func SelectRoleApiList(code string) (apis []vo.ApiResp) {
	apiPaths := global.Mysql.Model(&model.CasbinRule{}).Select("v1").Where("v0 = ?", code)
	global.Mysql.Model(&model.Api{}).Where("path in (?)", apiPaths).Scan(&apis)

	return
}

// 编辑角色Api
func EditRoleApi(ctx *gin.Context, editRoleApiReq dto.EditRoleApiReq) error {
	role, err := FindRoleById(editRoleApiReq.Id)
	if err != nil {
		AddFailOperation(ctx, "编辑角色API", "获取角色信息失败", nil, err)
		return errors.New("获取角色信息失败")
	}

	// 先删除后添加
	var delApis []model.Api
	global.Mysql.Where("id in (?)", editRoleApiReq.RemoveIds).Find(&delApis)
	for _, v := range delApis {
		global.Casbin.DeletePolicy(role.Code, v.Path, v.Method)
	}

	var addApis []model.Api
	global.Mysql.Where("id in (?)", editRoleApiReq.AddIds).Find(&addApis)

	for _, v := range addApis {
		global.Casbin.AddPolicy(role.Code, v.Path, v.Method)
	}

	return nil
}
