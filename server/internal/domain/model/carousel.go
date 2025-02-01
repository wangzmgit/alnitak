package model

import "gorm.io/gorm"

type Carousel struct {
	gorm.Model
	Uid         uint   `gorm:"comment:创建者"`
	Img         string `gorm:"type:varchar(255);comment:图片链接"`
	Title       string `gorm:"type:varchar(30);comment:标题"`
	Url         string `gorm:"type:varchar(100);comment:指向的链接"`
	Color       string `gorm:"type:varchar(20);comment:主题色"`
	Use         bool   `gorm:"comment:是否启用"`
	PartitionId uint   `gorm:"default:0;comment:分区ID"`
}

func (table *Carousel) TableName() string {
	return "carousel"
}
