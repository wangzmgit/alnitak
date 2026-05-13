package model

import "gorm.io/gorm"

type Danmaku struct {
	gorm.Model
	VideoShortID      string  `gorm:"type:varchar(16);comment:视频短ID;index"`
	Part              uint    `gorm:"comment:分集ID;default:1;index"`
	ResourceShortID  string  `gorm:"type:varchar(16);comment:资源短ID（精准绑定）;index"`
	Time              float32 `gorm:"comment:时间;not null"`
	Type              int     `gorm:"comment:类型:0滚动;1顶部;2底部;default:0"`
	Color             string  `gorm:"type:varchar(10);comment:颜色;default:'#fff'"`
	Text              string  `gorm:"type:varchar(100);comment:内容;not null"`
	Uid               uint    `gorm:"comment:用户ID;not null"`
}

func (table *Danmaku) TableName() string {
	return "danmaku"
}
