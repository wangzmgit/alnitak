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

func UploadVideoCreate(ctx *gin.Context) {
	video, err := ctx.FormFile("video")
	if err != nil {
		resp.FailWithMessage(ctx, "视频上传失败")
		return
	}

	resource, err := service.UploadVideoCreate(ctx, video)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"resource": resource})
}

func UploadVideoAdd(ctx *gin.Context) {
	// 获取视频ID
	vid := utils.StringToUint(ctx.Param("vid"))

	video, err := ctx.FormFile("video")
	if err != nil {
		resp.FailWithMessage(ctx, "视频上传失败")
		return
	}

	resource, err := service.UploadVideoAdd(ctx, vid, video)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"resource": resource})
}
