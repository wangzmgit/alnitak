package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 获取消息Websocket连接
func GetWhisperConnect(ctx *gin.Context) {
	// 升级为websocket长链接
	service.GetWhisperConnect(ctx)
}

// 发送私信
func SendWhisper(ctx *gin.Context) {
	//获取参数
	var whisperReq dto.WhisperReq
	if err := ctx.Bind(&whisperReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	//验证数据
	if whisperReq.Fid == 0 {
		resp.FailWithMessage(ctx, "接收者不存在")
		return
	}

	if utils.VerifyStringLength(whisperReq.Content, "=", 0) { // 角色名格式验证
		resp.FailWithMessage(ctx, "消息内容不能为空")
		return
	}

	if err := service.SendWhisper(ctx, whisperReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 获取消息列表
func GetWhisperList(ctx *gin.Context) {
	messages := service.GetWhisperList(ctx)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"messages": messages})
}

// 获取消息详细信息
func GetMessageDetails(ctx *gin.Context) {
	fid := utils.StringToUint(ctx.Query("fid"))
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	messages := service.GetMessageDetails(ctx, fid, page, pageSize)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"messages": messages})
}

// 已读消息
func ReadWhisper(ctx *gin.Context) {
	// 获取参数
	var idReq dto.IdReq
	if err := ctx.Bind(&idReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.ReadWhisper(ctx, idReq.ID); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}
