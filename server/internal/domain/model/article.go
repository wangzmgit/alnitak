package model

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title       string `gorm:"type:varchar(50);comment:标题;not null;index"`
	Cover       string `gorm:"type:varchar(255);comment:封面图;not null"`
	Content     string `gorm:"type:text;comment:内容;not null"`
	ContentDesc string `gorm:"type:varchar(300);comment:内容简介"`
	Uid         uint   `gorm:"comment:用户ID;not null;index"`
	Copyright   bool   `gorm:"comment:是否为原创;not null"`
	Clicks      int64  `gorm:"comment:点击量;default:0"`
	Status      int    `gorm:"comment:审核状态;not null"`
	PartitionId uint   `gorm:"comment:分区ID;deult:0"`
	Tags        string `gorm:"type:varchar(100);comment:标签;"`

	Author User `gorm:"-"`
}

func (table *Article) TableName() string {
	return "article"
}
