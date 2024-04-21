package model

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	Vid    uint   `gorm:"comment:'视频ID';not null;index"`
	Status int    `gorm:"comment:'审核状态';not null"`
	Remark string `gorm:"type:varchar(200);comment:'备注';not null;index"`
}

func (table *Review) TableName() string {
	return "review"
}
