package model

import "gorm.io/gorm"

// BackupUploadFailure 备用 OSS 上传失败记录，用于持久化存储失败信息以便后续重试。
// 上传成功或重试成功后删除对应记录。
type BackupUploadFailure struct {
	gorm.Model
	ObjectKey  string `gorm:"type:varchar(512);not null;comment:OSS对象键"`
	FilePath   string `gorm:"type:varchar(1024);not null;comment:本地文件路径"`
	Module     string `gorm:"type:varchar(50);not null;comment:来源模块(image/cover/subtitle/video)"`
	ErrMsg     string `gorm:"type:text;comment:上次错误信息"`
	RetryCount int    `gorm:"default:0;comment:已重试次数"`
}

func (table *BackupUploadFailure) TableName() string {
	return "backup_upload_failure"
}
