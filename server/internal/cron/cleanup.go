package cron

import (
	"fmt"

	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// CleanupOrphanedResources 定时清理孤立资源
func CleanupOrphanedResources() {
	utils.InfoLog("开始执行资源清理任务", "cron")

	result := service.ExecuteCleanup()

	utils.InfoLog(fmt.Sprintf("资源清理任务完成: 视频目录=%d, 图片=%d, 字幕=%d, 视频文件记录=%d, 索引文件记录=%d, 图片文件记录=%d, Resource记录=%d, 错误=%d",
		result.CleanedVideoDirs,
		result.CleanedImages,
		result.CleanedSubtitles,
		result.CleanedVideoFiles,
		result.CleanedIndexFiles,
		result.CleanedImageFiles,
		result.CleanedResources,
		len(result.Errors)), "cron")

	if len(result.Errors) > 0 {
		for _, err := range result.Errors {
			utils.ErrorLog("资源清理错误", "cron", err)
		}
	}
}
