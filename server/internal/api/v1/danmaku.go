package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 获取弹幕
func GetDanmaku(ctx *gin.Context) {
	raw := ctx.Query("vid")
	vid, err := service.ParseVideoID(raw)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	if vid == 0 {
		resp.FailWithMessage(ctx, "参数有误")
		return
	}

	part := utils.StringToUint(ctx.Query("part"))
	rid := ctx.Query("rid")

	// 优先通过 rid 获取弹幕，支持回退到 part
	danmaku, err := service.GetDanmakuByRid(ctx, vid, rid, part)
	if err == nil && len(danmaku) > 0 {
		resp.OkWithData(ctx, gin.H{"danmaku": danmaku})
		return
	}

	// 回退到旧的按 part 查询逻辑
	danmaku, err = service.GetDanmaku(ctx, vid, part)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"danmaku": danmaku})
}

// 发送弹幕
func SendDanmaku(ctx *gin.Context) {
	var danmakuReq dto.DanmakuReq
	if err := ctx.Bind(&danmakuReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(danmakuReq.Text, "=", 0) { // 弹幕内容不能为空
		resp.FailWithMessage(ctx, "弹幕内容不能为空")
		return
	}

	if err := service.SendDanmaku(ctx, danmakuReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}
