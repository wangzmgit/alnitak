package model

import "gorm.io/gorm"

type CollectArticle struct {
	gorm.Model
	Uid       uint `gorm:"comment:用户ID;not null;index"`
	Aid       uint `gorm:"comment:内容ID;not null"`
	IsCollect bool `gorm:"comment:是否收藏;default:false"` //是否点赞
}

func (table *CollectArticle) TableName() string {
	return "collect_article"
}
