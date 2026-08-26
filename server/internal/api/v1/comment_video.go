package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 发布评论回复
func AddVideoComment(ctx *gin.Context) {
	// 获取参数
	var addCommentReq dto.AddCommentReq
	if err := ctx.Bind(&addCommentReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// cid 兼容 shortId 或数字 ID，统一解析为数字 ID
	cid, err := service.ParseVideoID(addCommentReq.Cid)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	if utils.VerifyStringLength(addCommentReq.Content, "=", 0) {
		resp.FailWithMessage(ctx, "评论或回复内容不能为空")
		return
	}

	comment, err := service.AddVideoComment(ctx, addCommentReq, cid)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"comment": comment})
}

// 获取评论
func GetVideoComment(ctx *gin.Context) {
	cid, err := service.ParseVideoID(ctx.Query("vid"))
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 100 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	comments, total, err := service.GetVideoComment(ctx, cid, page, pageSize)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"comments": comments, "total": total})
}

// 获取回复
func GetVideoReply(ctx *gin.Context) {
	vid := utils.StringToUint(ctx.Query("commentId"))
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	replies, err := service.GetVideoReply(ctx, vid, page, pageSize)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"replies": replies})
}

// 删除评论回复
func DeleteVideoComment(ctx *gin.Context) {
	// 获取参数
	id := utils.StringToUint(ctx.Param("id"))

	if err := service.DeleteVideoComment(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 获取视频评论列表
// shortId 为空时走"当前用户全部视频"逻辑
func GetVideoCommentList(ctx *gin.Context) {
	raw := ctx.Query("shortId")
	var vid uint
	if raw != "" {
		parsed, err := service.ParseVideoID(raw)
		if err != nil {
			resp.FailWithMessage(ctx, err.Error())
			return
		}
		vid = parsed
	}
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 100 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	comments, total, err := service.GetVideoCommentList(ctx, vid, page, pageSize)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"comments": comments, "total": total})
}

// 点赞评论
func LikeComment(ctx *gin.Context) {
	commentId := utils.StringToUint(ctx.Param("id"))
	if commentId == 0 {
		resp.FailWithMessage(ctx, "评论ID不能为空")
		return
	}

	if err := service.LikeComment(ctx, commentId); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 取消点赞评论
func UnlikeComment(ctx *gin.Context) {
	commentId := utils.StringToUint(ctx.Param("id"))
	if commentId == 0 {
		resp.FailWithMessage(ctx, "评论ID不能为空")
		return
	}

	if err := service.UnlikeComment(ctx, commentId); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 点踩评论
func DislikeComment(ctx *gin.Context) {
	commentId := utils.StringToUint(ctx.Param("id"))
	if commentId == 0 {
		resp.FailWithMessage(ctx, "评论ID不能为空")
		return
	}

	if err := service.DislikeComment(ctx, commentId); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 取消点踩评论
func UndislikeComment(ctx *gin.Context) {
	commentId := utils.StringToUint(ctx.Param("id"))
	if commentId == 0 {
		resp.FailWithMessage(ctx, "评论ID不能为空")
		return
	}

	if err := service.UndislikeComment(ctx, commentId); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

