package model

import "gorm.io/gorm"

type LikeArticle struct {
	gorm.Model
	Uid    uint `gorm:"comment:用户ID;not null;index:idx_like_article_aid_uid"`
	Aid    uint `gorm:"comment:内容ID;not null;index:idx_like_article_aid_uid"`
	IsLike bool `gorm:"comment:是否点赞;default:false"` //是否点赞
}

func (table *LikeArticle) TableName() string {
	return "like_article"
}
