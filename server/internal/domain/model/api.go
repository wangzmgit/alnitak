package model

import "gorm.io/gorm"

type Api struct {
	gorm.Model
	Method   string `gorm:"type:varchar(20);comment:请求方式"`
	Path     string `gorm:"type:varchar(100);comment:访问路径"`
	Category string `gorm:"type:varchar(50);comment:所属类别"`
	Desc     string `gorm:"type:varchar(100);comment:说明"`
}

func (table *Api) TableName() string {
	return "api"
}
