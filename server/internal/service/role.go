package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func AddRole(ctx *gin.Context, addRoleReq dto.AddRoleReq) error {
	if role, _ := FindRoleByCode(addRoleReq.Code); role.ID != 0 {
		return errors.New("角色代码已存在")
	}

	// 保存到数据库
	return global.Mysql.Create(&model.Role{
		Name: addRoleReq.Name,
		Code: addRoleReq.Code,
		Desc: addRoleReq.Desc,
	}).Error
}

// 获取角色列表
func GetRoleList(page, pageSize int) (total int64, roles []vo.RoleResp) {
	global.Mysql.Model(&model.Role{}).Count(&total)
	global.Mysql.Model(&model.Role{}).Limit(pageSize).Offset((page - 1) * pageSize).Scan(&roles)
	return
}

// 获取所有的视频列表
func GetAllRoleList(ctx *gin.Context) (roles []vo.AllRoleResp) {
	global.Mysql.Model(&model.Role{}).Select("`name`,`code`").Scan(&roles)

	return
}

// 编辑角色
func EditRole(ctx *gin.Context, editRoleReq dto.EditRoleReq) error {
	if err := global.Mysql.Model(&model.Role{}).Where("id = ?", editRoleReq.ID).Updates(
		map[string]interface{}{
			"name": editRoleReq.Name,
			"desc": editRoleReq.Desc,
		},
	).Error; err != nil {
		utils.ErrorLog("修改角色失败", "role", err.Error())
		return errors.New("修改角色失败")
	}
	return nil
}

// 删除角色
func DeleteRole(ctx *gin.Context, id uint) error {
	// TODO: 查询是否有使用该角色的用户

	if err := global.Mysql.Where("id = ?", id).Delete(&model.Role{}).Error; err != nil {
		utils.ErrorLog("删除角色失败", "role", err.Error())
		return errors.New("删除角色失败")
	}

	return nil
}

// 编辑角色首页
func EditRoleHome(ctx *gin.Context, editRoleHomeReq dto.EditRoleHomeReq) error {
	if err := global.Mysql.Model(&model.Role{}).Where("id = ?", editRoleHomeReq.ID).Updates(
		map[string]interface{}{
			"home_page": editRoleHomeReq.Home,
		},
	).Error; err != nil {
		utils.ErrorLog("修改角色首页失败", "role", err.Error())
		return errors.New("修改首页失败")
	}
	return nil
}

func GetRoleInfo(code string) (role vo.RoleResp) {
	global.Mysql.Model(&model.Role{}).Where("`code` = ?", code).Scan(&role)

	return
}

// 通过角色ID查询角色
func FindRoleById(id uint) (role model.Role, err error) {
	err = global.Mysql.Where("`id` = ?", id).First(&role).Error
	return
}

// 通过角色代码查询角色
func FindRoleByCode(code string) (role model.Role, err error) {
	err = global.Mysql.Where("`code` = ?", code).First(&role).Error
	return
}
