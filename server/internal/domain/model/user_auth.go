package model

import "gorm.io/gorm"

// 用户认证记录表（已通过认证的用户）
type UserAuth struct {
	gorm.Model
	Uid       uint   `gorm:"comment:用户ID;not null;index"`
	AuthTypeID uint  `gorm:"comment:认证类型ID;not null;index"`
	AuthTypeCode string `gorm:"type:varchar(32);comment:认证类型code;not null"`
	Title     string `gorm:"type:varchar(64);comment:认证头衔，如：知名UP主"`
	Desc      string `gorm:"type:varchar(255);comment:认证描述"`
	ExtraData string `gorm:"type:text;comment:扩展数据JSON，如：{platform:\"bilibili\", follow:10000}"`

	// 关联
	AuthType *AuthType `gorm:"foreignKey:AuthTypeID"`
}

func (table *UserAuth) TableName() string {
	return "user_auth"
}