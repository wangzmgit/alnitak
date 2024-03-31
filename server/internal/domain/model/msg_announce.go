package model

import "gorm.io/gorm"

type Announce struct {
	gorm.Model
	Title   string `gorm:"type:varchar(50);comment:'标题';not null"`
	Content string `gorm:"type:varchar(200);comment:'内容';"`
	Url     string `gorm:"type:varchar(100);comment:'链接';"`
}

func (table *Announce) TableName() string {
	return "msg_announce"
}
