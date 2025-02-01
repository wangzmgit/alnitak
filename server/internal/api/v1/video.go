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

	if err := service.UploadVideoInfo(ctx, uploadVideoReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
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

// 获取用户视频
func GetVideoByUser(ctx *gin.Context) {
	userId := utils.StringToUint(ctx.Query("userId"))
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, videos := service.GetVideoByUser(ctx, userId, page, pageSize)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"total": total, "videos": videos})
}

// 获取视频列表(后台管理)
func GetVideoListManage(ctx *gin.Context) {
	// 获取参数
	var videoListReq dto.VideoListReq
	if err := ctx.Bind(&videoListReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if videoListReq.PageSize > 100 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, videos := service.GetVideoListManage(videoListReq)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"list": videos, "total": total})
}

// 删除视频(后台管理)
func DeleteVideoManage(ctx *gin.Context) {
	// 获取参数
	id := utils.StringToUint(ctx.Param("id"))

	if err := service.DeleteVideoManage(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 获取待审核视频列表
func GetReviewList(ctx *gin.Context) {
	// 获取参数
	var reviewListReq dto.ReviewListReq
	if err := ctx.Bind(&reviewListReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if reviewListReq.PageSize > 100 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, videos := service.GetReviewList(reviewListReq)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"list": videos, "total": total})
}

// 获取待审核视频资源
func GetReviewResourceList(ctx *gin.Context) {
	// 获取参数
	vid := utils.StringToUint(ctx.Query("vid"))
	resources := service.GetReviewResourceList(vid)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"resources": resources})
}

// 获取热门视频
func GetHotVideo(ctx *gin.Context) {
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	videos := service.GetHotVideo(ctx, page, pageSize)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"videos": videos})
}

// 获取分区视频
func GetVideoListByPartition(ctx *gin.Context) {
	size := utils.StringToInt(ctx.Query("size"))
	partitionId := utils.StringToUint(ctx.Query("partitionId"))

	if size > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	videos := service.GetVideoListByPartition(ctx, size, partitionId)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"videos": videos})
}

// 获取相关推荐视频
func GetRelatedVideoList(ctx *gin.Context) {
	videoId := utils.StringToUint(ctx.Query("vid"))

	videos := service.GetRelatedVideoList(ctx, videoId)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"videos": videos})
}

// 搜索视频
func SearchVideo(ctx *gin.Context) {
	// 获取参数
	var searchVideoReq dto.SearchVideoReq
	if err := ctx.Bind(&searchVideoReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if searchVideoReq.PageSize > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	videos := service.SearchVideo(ctx, searchVideoReq)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"videos": videos})
}
