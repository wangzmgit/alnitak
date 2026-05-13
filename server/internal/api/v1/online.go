package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/service"
)

// 视频Websocket连接(统计在线人数)
func GetVideoOnlineConnect(ctx *gin.Context) {
	raw := ctx.Query("vid")
	clientId := ctx.Query("clientId")
	rid := ctx.Query("rid") // 资源短ID，用于精准推送弹幕
	if raw == "" || clientId == "" {
		return
	}

	// 兼容短ID和数字ID：与 GetVideoById 同一解析逻辑
	videoId, err := service.ParseVideoID(raw)
	if err != nil || videoId == 0 {
		return
	}

	// 升级为websocket长链接
	service.GetVideoOnlineConnect(ctx, videoId, clientId, rid)
}
