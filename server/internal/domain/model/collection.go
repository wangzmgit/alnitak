package model

import "gorm.io/gorm"

//收藏夹
type Collection struct {
	gorm.Model
	Uid   uint   `gorm:"comment:所属用户;not null"`
	Name  string `gorm:"comment:收藏夹名称;type:varchar(20);"`
	Desc  string `gorm:"comment:简介;type:varchar(150);"`
	Cover string `gorm:"comment:封面;size:255;"`
	Open  bool   `gorm:"comment:是否公开;default:false"`
}

func (table *Collection) TableName() string {
	return "collection"
}
