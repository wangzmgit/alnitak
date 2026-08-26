package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
)

// ListBackupFailures 获取备用 OSS 上传失败记录列表。
func ListBackupFailures(ctx *gin.Context) {
	failures := service.ListBackupFailures()
	resp.OkWithData(ctx, failures)
}

// RetryBackupUpload 重试单条备用 OSS 上传失败记录。
func RetryBackupUpload(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		resp.FailWithMessage(ctx, "无效的ID")
		return
	}

	if err := service.RetryBackupUpload(uint(id)); err != nil {
		resp.FailWithMessage(ctx, "重试失败: "+err.Error())
		return
	}

	resp.OkWithMessage(ctx, "重试成功")
}

// RetryAllBackupUploads 重试所有备用 OSS 上传失败记录。
func RetryAllBackupUploads(ctx *gin.Context) {
	success, failed := service.RetryAllBackupUploads()
	resp.OkWithDetailed(ctx, gin.H{
		"success": success,
		"failed":  failed,
	}, "重试完成")
}
