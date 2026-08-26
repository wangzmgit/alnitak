package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 创建合集
func AddPlaylist(ctx *gin.Context) {
	var req dto.AddPlaylistReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if utils.VerifyStringLength(req.Title, "=", 0) {
		resp.FailWithMessage(ctx, "合集标题不能为空")
		return
	}

	if utils.VerifyStringLength(req.Cover, "=", 0) {
		resp.FailWithMessage(ctx, "合集封面不能为空")
		return
	}

	id, err := service.AddPlaylist(ctx, req)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithData(ctx, gin.H{"id": id})
}

// 编辑合集
func EditPlaylist(ctx *gin.Context) {
	var req dto.EditPlaylistReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if utils.VerifyStringLength(req.Title, "=", 0) {
		resp.FailWithMessage(ctx, "合集标题不能为空")
		return
	}

	if utils.VerifyStringLength(req.Cover, "=", 0) {
		resp.FailWithMessage(ctx, "合集封面不能为空")
		return
	}

	if err := service.EditPlaylist(ctx, req); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 删除合集
func DeletePlaylist(ctx *gin.Context) {
	id := utils.StringToUint(ctx.Param("id"))
	if id == 0 {
		resp.FailWithMessage(ctx, "参数有误")
		return
	}

	if err := service.DeletePlaylist(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 获取自己的合集列表
func GetMyPlaylist(ctx *gin.Context) {
	total, list := service.GetMyPlaylist(ctx)
	resp.OkWithData(ctx, gin.H{"total": total, "playlists": list})
}

// 获取合集详情
func GetPlaylistInfo(ctx *gin.Context) {
	id := utils.StringToUint(ctx.Query("id"))
	if id == 0 {
		resp.FailWithMessage(ctx, "参数有误")
		return
	}

	playlist, err := service.GetPlaylistInfo(ctx, id)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithData(ctx, gin.H{"playlist": playlist})
}

// 获取用户的公开合集列表
func GetUserPlaylist(ctx *gin.Context) {
	uid := utils.StringToUint(ctx.Query("uid"))
	page := utils.StringToInt(ctx.DefaultQuery("page", "1"))
	pageSize := utils.StringToInt(ctx.DefaultQuery("pageSize", "20"))
	if pageSize > 50 {
		pageSize = 50
	}
	if uid == 0 {
		resp.FailWithMessage(ctx, "参数有误")
		return
	}

	total, list := service.GetUserPlaylist(ctx, uid, page, pageSize)
	resp.OkWithData(ctx, gin.H{"total": total, "playlists": list})
}

// 添加视频到合集
func AddPlaylistVideo(ctx *gin.Context) {
	var req dto.PlaylistVideoAddReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.AddPlaylistVideo(ctx, req); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 从合集移除视频
func DelPlaylistVideo(ctx *gin.Context) {
	var req dto.PlaylistVideoDelReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.DelPlaylistVideo(ctx, req); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 调整合集视频排序
func SortPlaylistVideo(ctx *gin.Context) {
	var req dto.PlaylistVideoSortReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.SortPlaylistVideo(ctx, req); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 获取合集视频列表
func GetPlaylistVideoList(ctx *gin.Context) {
	playlistID := utils.StringToUint(ctx.Query("playlistId"))
	page := utils.StringToInt(ctx.DefaultQuery("page", "1"))
	pageSize := utils.StringToInt(ctx.DefaultQuery("pageSize", "20"))
	if pageSize > 50 {
		pageSize = 50
	}
	if playlistID == 0 {
		resp.FailWithMessage(ctx, "参数有误")
		return
	}

	total, list := service.GetPlaylistVideoList(ctx, playlistID, page, pageSize)
	resp.OkWithData(ctx, gin.H{"total": total, "videos": list})
}

// 获取合集视频列表（包含多分P展开和当前视频分P）
func GetPlaylistVideoListWithParts(ctx *gin.Context) {
	raw := ctx.Query("vid")
	videoID, err := service.ParseVideoID(raw)
	if err != nil || videoID == 0 {
		resp.FailWithMessage(ctx, "参数有误")
		return
	}

	data := service.GetPlaylistVideoListWithParts(ctx, videoID)
	resp.OkWithData(ctx, data)
}

// 获取全站合集列表（后台管理）
func GetPlaylistListManage(ctx *gin.Context) {
	page := utils.StringToInt(ctx.DefaultQuery("page", "1"))
	pageSize := utils.StringToInt(ctx.DefaultQuery("pageSize", "10"))
	if pageSize > 50 {
		pageSize = 50
	}

	total, list := service.GetPlaylistListManage(page, pageSize)
	resp.OkWithData(ctx, gin.H{"total": total, "list": list})
}

// 删除合集（后台管理）
func DeletePlaylistManage(ctx *gin.Context) {
	id := utils.StringToUint(ctx.Param("id"))
	if id == 0 {
		resp.FailWithMessage(ctx, "参数有误")
		return
	}

	if err := service.DeletePlaylistManage(id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 获取待审核合集列表
func GetReviewPlaylistList(ctx *gin.Context) {
	page := utils.StringToInt(ctx.DefaultQuery("page", "1"))
	pageSize := utils.StringToInt(ctx.DefaultQuery("pageSize", "10"))
	if pageSize > 50 {
		pageSize = 50
	}

	total, list := service.GetReviewPlaylistList(page, pageSize)
	resp.OkWithData(ctx, gin.H{"total": total, "list": list})
}

// 合集审核通过
func ReviewPlaylistApproved(ctx *gin.Context) {
	var req dto.ReviewPlaylistReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.ReviewPlaylistApproved(ctx, req); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 合集审核不通过
func ReviewPlaylistFailed(ctx *gin.Context) {
	var req dto.ReviewPlaylistReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.ReviewPlaylistFailed(ctx, req); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 获取合集审核记录
func GetPlaylistReviewRecord(ctx *gin.Context) {
	id := utils.StringToUint(ctx.Query("id"))
	if id == 0 {
		resp.FailWithMessage(ctx, "参数有误")
		return
	}

	review, err := service.GetPlaylistReviewRecord(ctx, id)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithData(ctx, gin.H{"review": review})
}

// GetMyPlaylistVideoIds 获取当前用户所有合集中的视频ID映射
func GetMyPlaylistVideoIds(ctx *gin.Context) {
	result := service.GetMyPlaylistVideoIds(ctx)
	resp.OkWithData(ctx, gin.H{"videoPlaylistMap": result})
}

// GetVideoPlaylists 获取视频所属的合集列表
func GetVideoPlaylists(ctx *gin.Context) {
	raw := ctx.Query("vid")
	videoId, err := service.ParseVideoID(raw)
	if err != nil || videoId == 0 {
		resp.FailWithMessage(ctx, "参数有误")
		return
	}

	total, lists := service.GetVideoPlaylists(videoId)

	resp.OkWithData(ctx, gin.H{"total": total, "playlists": lists})
}

// GetVideoPrimaryPlaylist 获取视频的主要合集
func GetVideoPrimaryPlaylist(ctx *gin.Context) {
	raw := ctx.Query("vid")
	videoId, err := service.ParseVideoID(raw)
	if err != nil || videoId == 0 {
		resp.FailWithMessage(ctx, "参数有误")
		return
	}
	userId := ctx.GetUint("userId") // 可能为0，表示未登录

	playlist, err := service.GetVideoPrimaryPlaylist(videoId, userId)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithData(ctx, playlist)
}
