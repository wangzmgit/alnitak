package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
)

// 获取邮箱配置信息
func GetEmailConfig(ctx *gin.Context) {
	config := service.GetEmailConfig()

	resp.OkWithData(ctx, gin.H{"config": config})
}

// 配置邮箱
func SetEmailConfig(ctx *gin.Context) {
	var emailConfigReq dto.EmailConfigReq
	if err := ctx.Bind(&emailConfigReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.SetEmailConfig(emailConfigReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

// 获取存储配置信息
func GetStorageConfig(ctx *gin.Context) {
	config := service.GetStorageConfig()

	resp.OkWithData(ctx, gin.H{"config": config})
}

// 修改存储配置
func SetStorageConfig(ctx *gin.Context) {
	var storageConfigReq dto.StorageConfigReq
	if err := ctx.Bind(&storageConfigReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.SetStorageConfig(storageConfigReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

// 获取其他配置信息
func GetOtherConfig(ctx *gin.Context) {
	config := service.GetOtherConfig()

	resp.OkWithData(ctx, gin.H{"config": config})
}

// 修改其他配置信息
func SetOtherConfig(ctx *gin.Context) {
	var otherConfigReq dto.OtherConfigReq
	if err := ctx.Bind(&otherConfigReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.SetOtherConfig(otherConfigReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}
