package model

import "gorm.io/gorm"

type LikeMessage struct {
	gorm.Model
	Vid uint `gorm:"comment:'所在视频id';not null"`
	Uid uint `gorm:"comment:'所属用户ID';not null"`
	Sid uint `gorm:"comment:'发送用户ID';not null"`
}

func (table *LikeMessage) TableName() string {
	return "msg_like"
}
