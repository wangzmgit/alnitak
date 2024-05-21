package model

import "gorm.io/gorm"

type LikeVideo struct {
	gorm.Model
	Uid    uint `gorm:"comment:用户ID;not null;index"`
	Vid    uint `gorm:"comment:视频ID;not null"`
	IsLike bool `gorm:"comment:是否点赞;default:false"` //是否点赞
}

func (table *LikeVideo) TableName() string {
	return "like_video"
}
