package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectPlayRoutes(r *gin.RouterGroup) {
	g := r.Group("play")
	{
		// 可选登录：带 Token 时 grant 内写入 uid，便于后续会员/版权策略
		g.POST("grant", middleware.OptionalAuth(), api.PostPlayGrant)
		g.GET(":resourceShortId", api.GetPlayByResource)
	}

	// 音轨信息（多音轨支持）
	audio := r.Group("audio")
	{
		audio.GET("tracks/:resourceId", api.GetAudioTracks)
	}
}
