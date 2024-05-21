package model

import (
	"gorm.io/gorm"
)

// 评论
type Comment struct {
	gorm.Model
	Cid           uint   `gorm:"comment:内容ID;index"`
	Uid           uint   `gorm:"comment:所属用户;index"`
	Content       string `gorm:"type:varchar(200);comment:评论内容"`
	AtUsernames   string `gorm:"type:varchar(100);comment:提及的用户名"`
	AtUserIds     string `gorm:"type:varchar(100);comment:提及的用户的ID"`
	ReplyUserID   uint   `gorm:"comment:回复的用户的ID"`
	ReplyUserName string `gorm:"type:varchar(20);comment:回复用户的用户名"`
	ParentId      uint   `gorm:"default:0;comment:所属评论ID"`
	Type          int    `gorm:"size:1;default:0;comment:类型:0视频、1文章"`
}

func (table *Comment) TableName() string {
	return "comment"
}
