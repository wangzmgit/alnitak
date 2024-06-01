package model

import (
	"gorm.io/gorm"
)

type Menu struct {
	gorm.Model
	Name      string `gorm:"type:varchar(50);comment:菜单名称"`
	Path      string `gorm:"type:varchar(100);comment:菜单访问路径"`
	Component string `gorm:"type:varchar(100);comment:前端组件路径" `
	Desc      string `gorm:"type:varchar(100);comment:介绍" `
	Sort      uint   `gorm:"type:int(3);default:1;comment:菜单排序"`
	Children  []Menu `gorm:"-" json:"children"`    // 子菜单集合
	Roles     []Role `gorm:"many2many:role_menu;"` // 角色菜单多对多关系
	ParentId  uint   `gorm:"default:0;comment:父菜单ID"`
	// Meta内容
	Title     string `gorm:"type:varchar(50);comment:菜单标题"`
	Icon      string `gorm:"type:varchar(50);comment:菜单图标"`
	Hidden    bool   `gorm:"default:false;comment:在菜单中隐藏"`
	KeepAlive bool   `gorm:"default:false;comment:缓存页面"`
}

func (table *Menu) TableName() string {
	return "menu"
}
