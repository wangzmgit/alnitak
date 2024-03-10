package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 获取分区
func GetPartitionList(ctx *gin.Context) {
	partitions := service.GetPartitionList()

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"partitions": partitions})
}

// 添加分区
func AddPartition(ctx *gin.Context) {
	// 获取参数
	var addPartitionReq dto.AddPartitionReq
	if err := ctx.Bind(&addPartitionReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(addPartitionReq.Name, "=", 0) { // 分区名格式验证
		resp.FailWithMessage(ctx, "分区名不能为空")
		return
	}

	if err := service.AddPartition(ctx, addPartitionReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 删除分区
func DeletePartition(ctx *gin.Context) {
	// 获取参数
	id := utils.StringToUint(ctx.Param("id"))

	if err := service.DeletePartition(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}
