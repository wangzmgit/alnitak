package model

import "gorm.io/gorm"

type CollectVideo struct {
	gorm.Model
	Uid          uint `gorm:"comment:用户ID;not null;index:idx_collect_video_collection_id_uid"`
	Vid          uint `gorm:"comment:视频ID;not null"`
	CollectionID uint `gorm:"comment:所属收藏夹ID;default:0;index:idx_collect_video_collection_id_uid"`
}

func (table *CollectVideo) TableName() string {
	return "collect_video"
}
