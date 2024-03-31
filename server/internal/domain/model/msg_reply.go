package model

import "gorm.io/gorm"

type ReplyMessage struct {
	gorm.Model
	Vid                uint   `gorm:"comment:'所在视频id';not null"`
	Uid                uint   `gorm:"comment:'所属用户ID';not null"`
	Sid                uint   `gorm:"comment:'关联用户ID';not null"`
	Content            string `gorm:"comment:'内容';not null"`
	TargetReplyContent string `gorm:"comment:'上级回复内容'"`
	RootContent        string `gorm:"comment:'根评论内容'"`
	CommentId          uint   `gorm:"comment:'评论ID';"`
}

func (table *ReplyMessage) TableName() string {
	return "msg_reply"
}
