package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

func UploadImg(ctx *gin.Context) {
	img, err := ctx.FormFile("image")
	if err != nil {
		resp.FailWithMessage(ctx, "文件上传失败")
		return
	}

	url, err := service.UploadImg(ctx, img)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"url": url})
}

func UploadVideo(ctx *gin.Context) {
	// 获取视频ID
	vid := utils.StringToUint(ctx.Param("vid"))
	if vid == 0 {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	video, err := ctx.FormFile("video")
	if err != nil {
		resp.FailWithMessage(ctx, "视频上传失败")
		return
	}

	if err := service.UploadVideo(ctx, vid, video); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}
