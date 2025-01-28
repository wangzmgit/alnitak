package model

import "gorm.io/gorm"

type VideoFile struct {
	gorm.Model
	Uid          uint   `gorm:"comment:用户ID;not null;index"`
	OriginalName string `gorm:"type:varchar(100);comment:原始文件名;"`
	DirName      string `gorm:"type:varchar(20);comment:目录名称;index"`
	Hash         string `gorm:"type:varchar(64);comment:文件hash;"`
	ChunksCount  int    `gorm:"comment:分片数量;"`
}

func (table *VideoFile) TableName() string {
	return "video_file"
}
