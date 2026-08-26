package cron

import "interastral-peace.com/alnitak/internal/service"

// RetryBackupFailures 定时重试备用 OSS 上传失败的记录。
func RetryBackupFailures() {
	service.AutoRetryBackupFailures()
}
