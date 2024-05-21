package model

import "gorm.io/gorm"

type CollectVideo struct {
	gorm.Model
	Uid          uint `gorm:"comment:用户ID;not null"`
	Vid          uint `gorm:"comment:视频ID;not null"`
	CollectionID uint `gorm:"comment:所属收藏夹ID;default:0"`
}

func (table *CollectVideo) TableName() string {
	return "collect_video"
}
