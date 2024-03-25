package model

import "gorm.io/gorm"

type History struct {
	gorm.Model
	Vid  uint    `gorm:"comment:'所在视频id';not null"`
	Uid  uint    `gorm:"comment:'所属用户ID';not null"`
	Part uint    `gorm:"comment:'分集';default:1"`
	Time float64 `gorm:"comment:'进度';not null"`

	Video Video `gorm:"-"` // 视频
}

func (table *History) TableName() string {
	return "history"
}
