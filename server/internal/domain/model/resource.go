package model

import "gorm.io/gorm"

type Resource struct {
	gorm.Model
	Vid       uint    `gorm:"comment:所属视频;index"`
	Uid       uint    `gorm:"comment:所属用户;index"`
	Title     string  `gorm:"type:varchar(50);comment:分P使用的标题"`
	CodecName string  `gorm:"type:varchar(10);comment:视频编码名称"`
	Duration  float64 `gorm:"comment:视频时长;default:0"`
	Status    int     `gorm:"comment:审核状态;not null;index"`
}

func (table *Resource) TableName() string {
	return "resource"
}
