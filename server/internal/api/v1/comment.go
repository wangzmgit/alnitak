package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 发布评论回复
func AddComment(ctx *gin.Context) {
	// 获取参数
	var addCommentReq dto.AddCommentReq
	if err := ctx.Bind(&addCommentReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if utils.VerifyStringLength(addCommentReq.Content, "=", 0) {
		resp.FailWithMessage(ctx, "评论或回复内容不能为空")
		return
	}

	comment, err := service.AddComment(ctx, addCommentReq)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"comment": comment})
}

// 获取评论
func GetComment(ctx *gin.Context) {
	vid := utils.StringToUint(ctx.Query("vid"))
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 100 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	comments, total, err := service.GetComment(ctx, vid, page, pageSize)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"comments": comments, "total": total})
}

// 获取回复
func GetReply(ctx *gin.Context) {
	vid := utils.StringToUint(ctx.Query("commentId"))
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	replies, err := service.GetReply(ctx, vid, page, pageSize)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"replies": replies})
}

// 删除评论回复
func DeleteComment(ctx *gin.Context) {
	// 获取参数
	id := utils.StringToUint(ctx.Param("id"))

	if err := service.DeleteComment(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 获取视频评论列表
func GetCommentList(ctx *gin.Context) {
	vid := utils.StringToUint(ctx.Query("vid"))
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 100 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	comments, total, err := service.GetCommentList(ctx, vid, page, pageSize)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"comments": comments, "total": total})
}
