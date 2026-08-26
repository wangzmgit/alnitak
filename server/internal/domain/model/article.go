package model

import (
	"gorm.io/gorm"
	"interastral-peace.com/alnitak/internal/domain/types"
)

type Article struct {
	gorm.Model
	Title       string `gorm:"type:varchar(50);comment:标题;not null;index"`
	Cover       string `gorm:"type:varchar(255);comment:封面图;not null"`
	Content     string `gorm:"type:text;comment:内容;not null"`
	ContentDesc string `gorm:"type:varchar(300);comment:内容简介"`
	Uid         uint   `gorm:"comment:用户ID;not null;index"`
	Copyright   types.CopyrightType `gorm:"comment:版权类型 0=未知 1=原创 2=转载 3=PGC授权;not null;default:0"`
	Clicks      int64  `gorm:"comment:点击量;default:0"`
	Status      int    `gorm:"comment:审核状态;not null"`
	PartitionId uint   `gorm:"comment:分区ID;default:0"`
	Tags        string `gorm:"type:varchar(100);comment:标签;"`
	Shares      int64  `gorm:"comment:分享数;default:0"`

	ShortID string `gorm:"type:varchar(16);comment:短ID;uniqueIndex" json:"shortId"`

	Author User `gorm:"-"`
}

func (table *Article) TableName() string {
	return "article"
}
