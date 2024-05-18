package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectHistoryRoutes(r *gin.RouterGroup) {
	historyGroup := r.Group("history")

	historyAuth := historyGroup.Group("")
	historyAuth.Use(middleware.Auth())
	{
		// 记录历史记录
		historyAuth.POST("video/addHistory", api.AddHistory)
		// 获取历史记录
		historyAuth.GET("video/getHistory", api.GetHistoryList)
		// 获取播放进度
		historyAuth.GET("video/getProgress", api.GetHistoryProgress)
	}
}
