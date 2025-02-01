package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectDanmakuRoutes(r *gin.RouterGroup) {
	danmakuGroup := r.Group("danmaku")

	danmakuAuth := danmakuGroup.Group("")
	danmakuAuth.Use(middleware.Auth())
	{
		// 发送弹幕
		danmakuAuth.POST("sendDanmaku", api.SendDanmaku)
	}

	// 获取弹幕列表
	danmakuGroup.GET("getDanmaku", api.GetDanmaku)
}
