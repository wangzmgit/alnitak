package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string    `gorm:"type:varchar(30);comment:用户名;not null;index"`
	Email      string    `gorm:"type:varchar(30);comment:邮箱;not null;index"`
	Password   string    `gorm:"type:varchar(128);comment:密码;not null"`
	Avatar     string    `gorm:"type:varchar(255);comment:头像"`
	SpaceCover string    `gorm:"type:varchar(255);comment:空间封面"`
	Gender     int       `gorm:"size:1;default:0;comment:性别:0未知、1男、3女"`
	Birthday   time.Time `gorm:"default:'1970-01-01';comment:生日"`
	Sign       string    `gorm:"type:varchar(50);comment:个性签名;default:'这个人很懒，什么都没有留下'"`
	ClientIp   string    `gorm:"type:varchar(20);comment:客户端IP"`
	Role       string    `gorm:"default:'001';comment:角色ID"`
}

func (table *User) TableName() string {
	return "user"
}
