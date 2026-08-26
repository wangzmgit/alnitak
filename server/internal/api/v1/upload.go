package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
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
	// 获取参数
	var videoFileReq dto.VideoFileReq
	if err := ctx.Bind(&videoFileReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	resource, err := service.UploadVideoCreate(ctx, videoFileReq)
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

	// 获取参数
	var videoFileReq dto.VideoFileReq
	if err := ctx.Bind(&videoFileReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	resource, err := service.UploadVideoAdd(ctx, vid, videoFileReq)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"resource": resource})
}

func UploadVideoChunk(ctx *gin.Context) {
	video, err := ctx.FormFile("video")
	if err != nil {
		resp.FailWithMessage(ctx, "视频上传失败")
		return
	}

	if err := service.UploadVideoChunk(ctx, video); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

func UploadVideoMerge(ctx *gin.Context) {
	var videoFileReq dto.VideoFileReq
	if err := ctx.Bind(&videoFileReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.UploadVideoMerge(ctx, videoFileReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

func UploadVideoCheck(ctx *gin.Context) {
	var videoFileReq dto.VideoFileReq
	if err := ctx.Bind(&videoFileReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	result, err := service.UploadVideoCheck(ctx, videoFileReq)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithData(ctx, gin.H{"chunks": result.Chunks, "fileID": result.FileID})
}

// ── 直传 OSS Handlers ──

func PresignImageUpload(ctx *gin.Context) {
	var req dto.PresignImageReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	data, err := service.PresignImageUpload(ctx, req)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithData(ctx, data)
}

func ConfirmImageUpload(ctx *gin.Context) {
	var req dto.ConfirmImageReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	url, err := service.ConfirmImageUpload(ctx, req)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithData(ctx, gin.H{"url": url})
}

func InitVideoUpload(ctx *gin.Context) {
	var req dto.InitVideoUploadReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	data, err := service.InitVideoUpload(ctx, req)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithData(ctx, data)
}

func CompleteVideoUpload(ctx *gin.Context) {
	var req dto.CompleteVideoUploadReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.CompleteVideoUpload(ctx, req); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

func PresignUploadChunks(ctx *gin.Context) {
	var req dto.PresignUploadChunksReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	data, err := service.PresignUploadChunks(ctx, req)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithData(ctx, data)
}
