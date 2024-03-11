package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 添加收藏夹
func AddCollection(ctx *gin.Context) {
	// 获取参数
	var addCollectionReq dto.AddCollectionReq
	if err := ctx.Bind(&addCollectionReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(addCollectionReq.Name, "=", 0) {
		resp.FailWithMessage(ctx, "收藏夹名称不能为空")
		return
	}

	if err := service.AddCollection(ctx, addCollectionReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 编辑收藏夹
func EditCollection(ctx *gin.Context) {
	var editCollectionReq dto.EditCollectionReq
	if err := ctx.Bind(&editCollectionReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(editCollectionReq.Name, "=", 0) {
		resp.FailWithMessage(ctx, "收藏夹名称不能为空")
		return
	}

	if err := service.EditCollection(ctx, editCollectionReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 删除收藏夹
func DeleteCollection(ctx *gin.Context) {
	// 获取参数
	id := utils.StringToUint(ctx.Param("id"))

	if err := service.DeleteCollection(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 获取收藏夹列表
func GetCollectionList(ctx *gin.Context) {
	page := utils.StringToInt(ctx.Query("page"))
	pageSize := utils.StringToInt(ctx.Query("pageSize"))

	if pageSize > 30 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, collections := service.GetCollectionList(ctx, page, pageSize)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"total": total, "collections": collections})
}

// 获取收藏夹信息
func GetCollectionInfo(ctx *gin.Context) {
	id := utils.StringToUint(ctx.Query("id"))

	collection, err := service.GetCollectionInfo(ctx, id)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"collection": collection})
}
