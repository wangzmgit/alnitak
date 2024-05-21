package model

import "gorm.io/gorm"

type LikeArticle struct {
	gorm.Model
	Uid    uint `gorm:"comment:用户ID;not null;index"`
	Aid    uint `gorm:"comment:内容ID;not null"`
	IsLike bool `gorm:"comment:是否点赞;default:false"` //是否点赞
}

func (table *LikeArticle) TableName() string {
	return "like_article"
}
