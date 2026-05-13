package model

import "gorm.io/gorm"

type LikeMessage struct {
	gorm.Model
	Cid        uint `gorm:"comment:内容ID（视频/文章ID）;not null"`
	CommentId  uint `gorm:"comment:评论ID（点赞评论时使用）;default:0"`
	Uid        uint `gorm:"comment:所属用户ID（被点赞者）;not null"`
	Sid        uint `gorm:"comment:发送用户ID（点赞者）;not null"`
	Type       int  `gorm:"size:1;default:0;comment:类型:0视频、1文章、2合集、3评论"`
	ParentType int  `gorm:"size:1;default:0;comment:评论所属类型:0视频、1文章"`
}

func (table *LikeMessage) TableName() string {
	return "msg_like"
}
