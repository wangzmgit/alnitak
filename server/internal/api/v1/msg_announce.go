package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 获取公告
func GetAnnounce(ctx *gin.Context) {
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, announces := service.GetAnnounce(ctx, page, pageSize)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"total": total, "announces": announces})
}

// 添加公告
func AddAnnounce(ctx *gin.Context) {
	//获取参数
	var addAnnounceReq dto.AddAnnounceReq
	if err := ctx.Bind(&addAnnounceReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(addAnnounceReq.Title, "=", 0) { // 角色名格式验证
		resp.FailWithMessage(ctx, "公告标题不能为空")
		return
	}

	if utils.VerifyStringLength(addAnnounceReq.Content, "=", 0) { // 角色名格式验证
		resp.FailWithMessage(ctx, "公告内容不能为空")
		return
	}

	if err := service.AddAnnounce(ctx, addAnnounceReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 删除公告
func DeleteAnnounce(ctx *gin.Context) {
	// 获取参数
	var idReq dto.IdReq
	if err := ctx.Bind(&idReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.DeleteAnnounce(ctx, idReq.ID); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}
