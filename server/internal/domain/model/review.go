package model

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	Cid    uint   `gorm:"comment:内容ID;not null;index"`
	Status int    `gorm:"comment:审核状态;not null"`
	Remark string `gorm:"type:varchar(200);comment:备注;not null;index"`
	Uid    uint   `gorm:"comment:审核人;"`
	Type   int    `gorm:"size:1;default:0;comment:类型:0视频、1文章"`
}

func (table *Review) TableName() string {
	return "review"
}
