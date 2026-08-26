package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectBackupRoutes(r *gin.RouterGroup) {
	backupGroup := r.Group("backup")
	backupGroup.Use(middleware.Auth())
	{
		// 获取上传失败记录列表
		backupGroup.GET("failures", api.ListBackupFailures)
		// 重试单条失败记录
		backupGroup.POST("retry/:id", api.RetryBackupUpload)
		// 重试所有失败记录
		backupGroup.POST("retryAll", api.RetryAllBackupUploads)
	}
}
