package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectMessageRoutes(r *gin.RouterGroup) {
	messageGroup := r.Group("message")
	messageAuth := messageGroup.Group("")
	messageAuth.Use(middleware.Auth())
	{
		// 添加公告
		messageAuth.POST("addAnnounce", api.AddAnnounce)
		// 删除公告
		messageAuth.DELETE("deleteAnnounce/:id", api.DeleteAnnounce)

		// 获取点赞消息
		messageAuth.GET("getLikeMsg", api.GetLikeMessage)
		// 获取AT消息
		messageAuth.GET("getAtMsg", api.GetAtMessage)
		// 获取回复消息
		messageAuth.GET("getReplyMsg", api.GetReplyMessage)

		// 获取私信列表
		messageAuth.GET("getWhisperList", api.GetWhisperList)
		// 获取私信详情
		messageAuth.GET("getWhisperDetails", api.GetMessageDetails)
		// 发送私信
		messageAuth.POST("sendWhisper", api.SendWhisper)
		// 已读私信
		messageAuth.POST("readWhisper", api.ReadWhisper)
	}

	// 获取公告
	messageGroup.GET("getAnnounce", api.GetAnnounce)
	// 连接消息ws
	messageGroup.GET("ws", middleware.WsAuth(), api.GetWhisperConnect)
}
