package model

import "gorm.io/gorm"

type LikeMessage struct {
	gorm.Model
	Cid  uint `gorm:"comment:内容ID;not null"`
	Uid  uint `gorm:"comment:所属用户ID;not null"`
	Sid  uint `gorm:"comment:发送用户ID;not null"`
	Type int  `gorm:"size:1;default:0;comment:类型:0视频、1文章"`
}

func (table *LikeMessage) TableName() string {
	return "msg_like"
}
