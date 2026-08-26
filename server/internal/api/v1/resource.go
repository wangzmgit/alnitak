package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 编辑资源标题
func ModifyResourceTitle(ctx *gin.Context) {
	// 获取参数
	var modifyTitleReq dto.ModifyResourceTitleReq
	if err := ctx.Bind(&modifyTitleReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(modifyTitleReq.Title, "=", 0) { // 角色名格式验证
		resp.FailWithMessage(ctx, "资源名不能为空")
		return
	}

	// 更新数据库
	if err := service.ModifyResourceTitle(ctx, modifyTitleReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 删除资源
func DeleteResource(ctx *gin.Context) {
	// 获取参数
	id := utils.StringToUint(ctx.Param("id"))

	// 获取是否删除弹幕选项（默认为false）
	deleteDanmaku := ctx.PostForm("deleteDanmaku") == "true"

	if err := service.DeleteResource(ctx, id, deleteDanmaku); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 获取视频资源支持的分辨率信息
func GetResourceQuality(ctx *gin.Context) {
	resourceId, parseErr := service.ParseResourceID(ctx.Query("resourceId"))
	if parseErr != nil {
		resp.FailWithMessage(ctx, parseErr.Error())
		return
	}

	info, err := service.GetResourceQuality(ctx, resourceId)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端（包含 quality、supportsDash 和 qualityOrder）
	resp.OkWithData(ctx, gin.H{
		"quality":      info.Quality,
		"supportsDash": info.SupportsDash,
		"qualityOrder": info.QualityOrder,
	})
}

// 获取视频资源支持的分辨率信息(后台管理)
func GetResourceQualityManage(ctx *gin.Context) {
	// 禁用缓存，确保审核时看到最新数据
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	resourceId := utils.StringToUint(ctx.Query("resourceId"))

	info, err := service.GetResourceQualityManage(ctx, resourceId)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端（包含 quality、supportsDash 和 qualityOrder）
	resp.OkWithData(ctx, gin.H{
		"quality":      info.Quality,
		"supportsDash": info.SupportsDash,
		"qualityOrder": info.QualityOrder,
	})
}

// 资源排序
func ReorderResources(ctx *gin.Context) {
	var req dto.ReorderResourceReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if req.Vid == 0 || len(req.ResourceIds) == 0 {
		resp.FailWithMessage(ctx, "参数不能为空")
		return
	}

	if err := service.ReorderResources(ctx, req); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 检查资源是否需要替换
func CheckReplaceResource(ctx *gin.Context) {
	// 获取参数
	var replaceReq dto.ReplaceResourceReq
	if err := ctx.Bind(&replaceReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if replaceReq.ResourceID == 0 {
		resp.FailWithMessage(ctx, "资源ID不能为空")
		return
	}
	if replaceReq.Hash == "" {
		resp.FailWithMessage(ctx, "文件Hash不能为空")
		return
	}

	// 检查是否需要替换
	if err := service.CheckReplaceResource(ctx, replaceReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// hash不同，可以替换
	resp.Ok(ctx)
}

// 替换资源
func ReplaceResource(ctx *gin.Context) {
	// 获取参数
	var replaceReq dto.ReplaceResourceReq
	if err := ctx.Bind(&replaceReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 支持从query参数获取resourceId（用于分片上传完成后的回调）
	if replaceReq.ResourceID == 0 {
		replaceReq.ResourceID = utils.StringToUint(ctx.Query("resourceId"))
	}

	// 参数校验
	if replaceReq.ResourceID == 0 {
		resp.FailWithMessage(ctx, "资源ID不能为空")
		return
	}
	if replaceReq.Hash == "" {
		resp.FailWithMessage(ctx, "文件Hash不能为空")
		return
	}

	// 替换资源
	resource, err := service.ReplaceResource(ctx, replaceReq)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.OkWithData(ctx, gin.H{"resource": resource})
}
