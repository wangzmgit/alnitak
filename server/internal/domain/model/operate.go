package model

import "gorm.io/gorm"

type Operate struct {
	gorm.Model
	Position string `gorm:"type:varchar(50);comment:'位置'"`
	Success  bool   `gorm:"comment:'是否成功'"`
	Param    string `gorm:"type:text;comment:'参数'"`
	Msg      string `gorm:"type:varchar(50);;comment:'信息'"`
	Error    string `gorm:"type:text;comment:'错误详情'"`
	Ip       string `gorm:"comment:'请求ip'"`
	UserID   uint   `gorm:"default:0;comment:'用户ID'"`
}

func (table *Operate) TableName() string {
	return "operate"
}
