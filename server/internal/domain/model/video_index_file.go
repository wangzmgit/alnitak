package model

import "gorm.io/gorm"

type VideoIndexFile struct {
	gorm.Model
	ResourceID uint   `gorm:"index;comment:视频资源ID;"`
	Quality    string `gorm:"comment:视频质量;"`
	DirName    string `gorm:"type:varchar(20);comment:目录名称;"`
	Content    string `gorm:"type:text;comment:文件内容;"`
}

func (table *VideoIndexFile) TableName() string {
	return "video_index_file"
}
