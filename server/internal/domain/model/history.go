package model

import "gorm.io/gorm"

type History struct {
	gorm.Model
	Vid  uint    `gorm:"comment:所在视频id;not null;uniqueIndex:idx_history_unique"`
	Uid  uint    `gorm:"comment:所属用户ID;not null;uniqueIndex:idx_history_unique"`
	Part uint    `gorm:"comment:分集;default:1;uniqueIndex:idx_history_unique"`
	Time float64 `gorm:"comment:进度;not null"`
}

func (table *History) TableName() string {
	return "history"
}
