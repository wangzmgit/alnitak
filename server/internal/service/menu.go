package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

// 新增菜单
func AddMenu(ctx *gin.Context, addMenuReq dto.AddMenuReq) error {
	if menu := FindMenusByName(addMenuReq.Name); menu.ID != 0 {
		return errors.New("路由Name已存在")
	}

	// 保存到数据库
	if err := global.Mysql.Create(&model.Menu{
		Name:      addMenuReq.Name,
		Path:      addMenuReq.Path,
		Component: addMenuReq.Component,
		Desc:      addMenuReq.Desc,
		Sort:      addMenuReq.Sort,
		ParentId:  addMenuReq.ParentId,
		Icon:      addMenuReq.Icon,
		Title:     addMenuReq.Title,
		Hidden:    addMenuReq.Hidden,
		KeepAlive: addMenuReq.KeepAlive,
	}).Error; err != nil {
		utils.ErrorLog("新增菜单更新数据库失败", "menu", err.Error())
		return errors.New("新增菜单失败")
	}
	return nil
}

// 编辑菜单
func EditMenu(ctx *gin.Context, editMenuReq dto.EditMenuReq) error {
	if menu := FindMenusByName(editMenuReq.Name); menu.ID != 0 && menu.ID != editMenuReq.ID {
		return errors.New("路由Name已存在")
	}

	if err := global.Mysql.Model(&model.Menu{}).Where("id = ?", editMenuReq.ID).Updates(
		map[string]interface{}{
			"name":       editMenuReq.Name,
			"path":       editMenuReq.Path,
			"component":  editMenuReq.Component,
			"desc":       editMenuReq.Desc,
			"sort":       editMenuReq.Sort,
			"icon":       editMenuReq.Icon,
			"title":      editMenuReq.Title,
			"hidden":     editMenuReq.Hidden,
			"keep_alive": editMenuReq.KeepAlive,
		},
	).Error; err != nil {
		utils.ErrorLog("编辑菜单更新数据库失败", "menu", err.Error())
		return errors.New("编辑菜单失败")
	}
	return nil
}

// 删除菜单
func DeleteMenu(ctx *gin.Context, id uint) error {
	if err := global.Mysql.Where("id = ?", id).Delete(&model.Menu{}).Error; err != nil {
		utils.ErrorLog("删除菜单更新数据库失败", "menu", err.Error())

		return errors.New("删除菜单失败")
	}

	return nil
}

// 获取菜单树
func GetMenuTree() []model.Menu {
	menus := FindMenusByParentId(0)
	for i := 0; i < len(menus); i++ {
		menus[i].Children = FindMenusByParentId(menus[i].ID)
	}

	return menus
}

// 通过角色ID获取角色菜单树
func GetMenuTreeByRoleId(roleId uint) []model.Menu {
	var role model.Role
	global.Mysql.Where("id = ?", roleId).Preload("Menus").First(&role)

	return role.Menus
}

// 通过角色Code获取角色菜单树
func GetMenuTreeByRoleCode(code string) []model.Menu {
	var role model.Role
	global.Mysql.Where("code = ?", code).Preload("Menus").First(&role)

	return role.Menus
}

// 更新角色菜单
func EditRoleMenu(ctx *gin.Context, editMenuReq dto.EditRoleMenuReq) error {
	role, err := FindRoleById(editMenuReq.Id)
	if err != nil {
		utils.ErrorLog("更新角色菜单获取角色信息失败", "menu", err.Error())
		return errors.New("获取角色信息失败")
	}

	role.Menus = FindMenuListByIds(editMenuReq.MenuIds)
	global.Mysql.Model(&role).Association("Menus").Replace(role.Menus)

	return nil
}

// 通过parentId查询菜单
func FindMenusByParentId(parentId uint) (menus []model.Menu) {
	global.Mysql.Where("parent_id = ?", parentId).Order("sort").Find(&menus)

	return
}

// 通过路由Name查询菜单
func FindMenusByName(name string) (menu model.Menu) {
	global.Mysql.Where("name = ?", name).First(&menu)

	return
}

// 通过ID数组查询菜单列表
func FindMenuListByIds(ids []uint) (menus []model.Menu) {
	global.Mysql.Where("id in (?)", ids).Find(&menus)
	return
}
