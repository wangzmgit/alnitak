package service

import (
	"fmt"
	"time"

	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

// UploadToBackupWithRetry 上传文件到备用 OSS，带指数退避重试。
// 重试全部失败后写入 backup_upload_failure 表持久化，供后续手动/自动重试。
func UploadToBackupWithRetry(objectKey, filePath, module string) {
	backup := global.StorageBackup
	if backup == nil {
		return
	}

	var lastErr error
	for attempt := 0; attempt <= ossUploadMaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(ossUploadBackoff[attempt-1])
		}
		lastErr = backup.PutObjectFromFile(objectKey, filePath)
		if lastErr == nil {
			// 成功，清除之前的失败记录
			clearBackupFailure(objectKey)
			return
		}
	}

	// 所有重试均失败，持久化记录
	utils.ErrorLog(fmt.Sprintf("【备用OSS上传失败】%s (重试%d次)", objectKey, ossUploadMaxRetries), module, lastErr.Error())
	saveBackupFailure(objectKey, filePath, module, lastErr.Error())
}

// RecordBackupUploadFailure 直接记录一次上传失败（适用于调用方已自行处理重试的场景，如转码产物上传）。
// 不影响调用方已有的重试逻辑，仅在最终失败时落库。
func RecordBackupUploadFailure(objectKey, filePath, module, errMsg string) {
	saveBackupFailure(objectKey, filePath, module, errMsg)
}

// clearBackupFailure 删除指定 objectKey 的失败记录（上传成功后调用）。
func clearBackupFailure(objectKey string) {
	global.Mysql.Where("object_key = ?", objectKey).Delete(&model.BackupUploadFailure{})
}

// saveBackupFailure 写入/更新失败记录（幂等：相同 objectKey 会覆盖）。
func saveBackupFailure(objectKey, filePath, module, errMsg string) {
	// 先删旧记录，再插新记录（简单幂等）
	global.Mysql.Where("object_key = ?", objectKey).Delete(&model.BackupUploadFailure{})
	global.Mysql.Create(&model.BackupUploadFailure{
		ObjectKey: objectKey,
		FilePath:  filePath,
		Module:    module,
		ErrMsg:    errMsg,
	})
}

// RetryBackupUpload 重试单条失败记录。成功则删除记录，失败则更新错误信息与重试次数。
func RetryBackupUpload(id uint) error {
	var record model.BackupUploadFailure
	if err := global.Mysql.First(&record, id).Error; err != nil {
		return err
	}

	backup := global.StorageBackup
	if backup == nil {
		return fmt.Errorf("备用 OSS 未配置")
	}

	if err := backup.PutObjectFromFile(record.ObjectKey, record.FilePath); err != nil {
		// 更新错误信息与重试次数
		global.Mysql.Model(&record).Updates(map[string]interface{}{
			"err_msg":     err.Error(),
			"retry_count": record.RetryCount + 1,
		})
		return err
	}

	// 成功，删除记录
	global.Mysql.Delete(&record)
	return nil
}

// RetryAllBackupUploads 重试所有失败记录。返回成功与失败数量。
func RetryAllBackupUploads() (success, failed int) {
	var records []model.BackupUploadFailure
	global.Mysql.Find(&records)
	for _, r := range records {
		if err := RetryBackupUpload(r.ID); err != nil {
			failed++
		} else {
			success++
		}
	}
	return
}

// ListBackupFailures 返回所有上传失败记录，按创建时间倒序。
func ListBackupFailures() []model.BackupUploadFailure {
	var records []model.BackupUploadFailure
	global.Mysql.Order("created_at desc").Find(&records)
	return records
}

// AutoRetryBackupFailures 供 cron 调用的自动重试：重试所有失败记录。
func AutoRetryBackupFailures() {
	success, failed := RetryAllBackupUploads()
	total := success + failed
	if total > 0 {
		utils.InfoLog(fmt.Sprintf("【备用OSS自动重试】完成: 成功%d, 失败%d", success, failed), "backup")
	}
}
