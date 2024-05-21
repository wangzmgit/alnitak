package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Title       string  `gorm:"type:varchar(50);comment:标题;not null;index"`
	Cover       string  `gorm:"type:varchar(255);cmment:封面图;not null"`
	Desc        string  `gorm:"type:varchar(200);comment:视频简介;default:什么都没有~"`
	Uid         uint    `gorm:"comment:用户ID;not null;index"`
	Copyright   bool    `gorm:"comment:是否为原创;not null"`
	Clicks      int64   `gorm:"comment:点击量;default:0"`
	Status      int     `gorm:"comment:审核状态;not null"`
	PartitionId uint    `gorm:"comment:分区ID;deult:0"`
	Tags        string  `gorm:"type:varchar(100);comment:标签;"`
	Duration    float64 `gorm:"comment:视频时长;default:0"`

	Author User `gorm:"-"`
}

func (table *Video) TableName() string {
	return "video"
}
