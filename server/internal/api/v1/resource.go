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

	if err := service.DeleteResource(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 获取视频资源支持的分辨率信息
func GetResourceQuality(ctx *gin.Context) {
	resourceId := utils.StringToUint(ctx.Query("resourceId"))

	quality, err := service.GetResourceQuality(ctx, resourceId)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"quality": quality})
}
