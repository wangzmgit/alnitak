package model

import "gorm.io/gorm"

type Relation struct {
	gorm.Model
	Uid       uint `gorm:"comment:用户ID;not null;index"`
	TargetUid uint `gorm:"comment:目标用户ID;not null;index"`
	Relation  int  `gorm:"comment:关系:0未关注、1关注、2互相关注;default:1"`
}

func (table *Relation) TableName() string {
	return "relation"
}
