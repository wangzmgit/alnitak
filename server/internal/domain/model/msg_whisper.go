package model

import "gorm.io/gorm"

type Whisper struct {
	gorm.Model
	Uid     uint   `gorm:"comment:用户ID;not null;"`
	Fid     uint   `gorm:"comment:好友ID;not null;"`
	FromId  uint   `gorm:"comment:发送者;not null;"`
	ToId    uint   `gorm:"comment:接受者;not null;"`
	Content string `gorm:"size:255;comment:内容;"`
	Status  bool   `gorm:"comment:已读状态;default:false"`
}

func (table *Whisper) TableName() string {
	return "msg_whisper"
}
