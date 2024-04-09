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

// 提交审核
func SubmitReview(ctx *gin.Context) {
	//获取参数
	var idReq dto.IdReq
	if err := ctx.Bind(&idReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	err := service.SubmitReview(ctx, idReq.ID)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

// 获取自己的视频
func GetUploadVideoList(ctx *gin.Context) {
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, videos := service.GetUploadVideoList(ctx, page, pageSize)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"total": total, "videos": videos})
}

// 获取视频信息
func GetVideoById(ctx *gin.Context) {
	vid := utils.StringToUint(ctx.DefaultQuery("vid", "0"))

	video, err := service.GetVideoById(ctx, vid)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"video": video})
}

// 获取所有的视频列表
func GetAllVideoList(ctx *gin.Context) {
	videos := service.GetAllVideoList(ctx)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"videos": videos})
}

// 编辑视频信息
func EditVideoInfo(ctx *gin.Context) {
	var editVideoReq dto.EditVideoReq
	if err := ctx.Bind(&editVideoReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(editVideoReq.Title, "=", 0) {
		resp.FailWithMessage(ctx, "视频标题不能为空")
		return
	}

	if err := service.EditVideoInfo(ctx, editVideoReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

// 删除视频
func DeleteVideo(ctx *gin.Context) {
	// 获取参数
	id := utils.StringToUint(ctx.Param("id"))

	if err := service.DeleteVideo(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}
