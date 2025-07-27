package service

import (
	"time"

	"go.uber.org/zap"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
)

// CleanupOldOperateLogs 清理旧的操作日志
// 保留最近30天的日志，删除更早的记录
func CleanupOldOperateLogs() {
	// 计算30天前的时间
	cutoffTime := time.Now().AddDate(0, 0, -30)

	// 删除30天前的操作日志
	result := global.Mysql.Where("created_at < ?", cutoffTime).Delete(&model.Operate{})

	if result.Error != nil {
		zap.L().Error("清理操作日志失败",
			zap.String("module", "operate_cleanup"),
			zap.Error(result.Error))
		return
	}

	if result.RowsAffected > 0 {
		zap.L().Info("操作日志清理完成",
			zap.String("module", "operate_cleanup"),
			zap.Int64("deleted_rows", result.RowsAffected),
			zap.Time("cutoff_time", cutoffTime))
	}
}

// CleanupOldOperateLogsBatch 批量清理旧的操作日志
// 分批删除以避免长时间锁表
func CleanupOldOperateLogsBatch() {
	cutoffTime := time.Now().AddDate(0, 0, -30)
	batchSize := 1000
	totalDeleted := int64(0)

	for {
		// 分批删除
		result := global.Mysql.Where("created_at < ?", cutoffTime).
			Limit(batchSize).
			Delete(&model.Operate{})

		if result.Error != nil {
			zap.L().Error("批量清理操作日志失败",
				zap.String("module", "operate_cleanup"),
				zap.Error(result.Error))
			break
		}

		totalDeleted += result.RowsAffected

		// 如果删除的记录数少于批次大小，说明已经清理完成
		if result.RowsAffected < int64(batchSize) {
			break
		}

		// 短暂休息，避免对数据库造成过大压力
		time.Sleep(100 * time.Millisecond)
	}

	if totalDeleted > 0 {
		zap.L().Info("批量操作日志清理完成",
			zap.String("module", "operate_cleanup"),
			zap.Int64("total_deleted", totalDeleted),
			zap.Time("cutoff_time", cutoffTime))
	}
}
