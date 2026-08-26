package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 发布文章评论回复
func AddArticleComment(ctx *gin.Context) {
	// 获取参数
	var addCommentReq dto.AddCommentReq
	if err := ctx.Bind(&addCommentReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// cid 兼容 shortId 或数字 ID，统一解析为数字 ID
	cid, err := service.ParseArticleID(addCommentReq.Cid)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	if utils.VerifyStringLength(addCommentReq.Content, "=", 0) {
		resp.FailWithMessage(ctx, "评论或回复内容不能为空")
		return
	}

	comment, err := service.AddArticleComment(ctx, addCommentReq, cid)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"comment": comment})
}

// 获取评论
func GetArticleComment(ctx *gin.Context) {
	raw := ctx.Query("aid")
	cid, err := service.ParseArticleID(raw)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	if cid == 0 {
		resp.FailWithMessage(ctx, "参数有误")
		return
	}
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 100 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	comments, total, err := service.GetArticleComment(ctx, cid, page, pageSize)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"comments": comments, "total": total})
}

// 获取回复
func GetArticleReply(ctx *gin.Context) {
	vid := utils.StringToUint(ctx.Query("commentId"))
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	replies, err := service.GetArticleReply(ctx, vid, page, pageSize)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"replies": replies})
}

// 删除评论回复
func DeleteArticleComment(ctx *gin.Context) {
	// 获取参数
	id := utils.StringToUint(ctx.Param("id"))

	if err := service.DeleteArticleComment(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 获取文章评论列表
// shortId 为空时走"当前用户全部文章"逻辑
func GetArticleCommentList(ctx *gin.Context) {
	raw := ctx.Query("shortId")
	var aid uint
	if raw != "" {
		parsed, err := service.ParseArticleID(raw)
		if err != nil {
			resp.FailWithMessage(ctx, err.Error())
			return
		}
		aid = parsed
	}
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 100 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	comments, total, err := service.GetArticleCommentList(ctx, aid, page, pageSize)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"comments": comments, "total": total})
}
