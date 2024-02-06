package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 上传视频信息
func UploadVideoInfo(ctx *gin.Context) {
	var uploadVideoReq dto.UploadVideoReq
	if err := ctx.Bind(&uploadVideoReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(uploadVideoReq.Title, "=", 0) {
		resp.FailWithMessage(ctx, "视频标题不能为空")
		return
	}

	videoId, err := service.UploadVideoInfo(ctx, uploadVideoReq)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"videoId": videoId})
}

// 获取视频状态
func GetVideoStatus(ctx *gin.Context) {
	videoId := utils.StringToUint(ctx.Query("vid"))

	video, err := service.GetVideoStatus(ctx, videoId)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"video": video})
}
