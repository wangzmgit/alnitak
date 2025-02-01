package model

import "gorm.io/gorm"

type Partition struct {
	gorm.Model
	Name     string `gorm:"varchar(20);comment:分区名称"`
	Type     int    `gorm:"size:1;default:0;comment:类型:0视频、1文章"`
	ParentId uint   `gorm:"default:0;comment:所属分区ID"`
}

func (table *Partition) TableName() string {
	return "partition"
}
