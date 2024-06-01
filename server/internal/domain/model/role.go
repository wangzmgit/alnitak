package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20);not null;unique;comment:角色名"`
	Code     string `gorm:"type:varchar(20);not null;unique;comment:角色代码"`
	Desc     string `gorm:"type:varchar(100);comment:介绍"`
	HomePage string `gorm:"type:varchar(20);comment:角色首页"`
	Menus    []Menu `gorm:"many2many:role_menu;"` // 角色菜单多对多关系
}

func (table *Role) TableName() string {
	return "role"
}
