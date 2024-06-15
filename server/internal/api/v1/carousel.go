package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 获取轮播图
func GetCarousel(ctx *gin.Context) {
	partitionId := utils.StringToUint(ctx.Query("partitionId"))

	carousels := service.GetCarousel(ctx, partitionId)

	// 返回
	resp.OkWithData(ctx, gin.H{"carousels": carousels})
}

// 获取全部轮播图
func GetCarouselList(ctx *gin.Context) {
	// 获取参数
	var carouselListReq dto.CarouselListReq
	if err := ctx.Bind(&carouselListReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	total, carousels := service.GetCarouselList(ctx, carouselListReq)

	// 返回
	resp.OkWithData(ctx, gin.H{"total": total, "list": carousels})
}

// 上传轮播图信息
func AddCarousel(ctx *gin.Context) {
	// 获取参数
	var addCarouselReq dto.AddCarouselReq
	if err := ctx.Bind(&addCarouselReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if utils.VerifyStringLength(addCarouselReq.Title, "=", 0) { //
		resp.FailWithMessage(ctx, "标题不能为空")
		return
	}

	if err := service.AddCarousel(ctx, addCarouselReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 编辑轮播图信息
func EditCarousel(ctx *gin.Context) {
	// 获取参数
	var editCarouselReq dto.EditCarouselReq
	if err := ctx.Bind(&editCarouselReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if utils.VerifyStringLength(editCarouselReq.Title, "=", 0) { //
		resp.FailWithMessage(ctx, "标题不能为空")
		return
	}

	if err := service.EditCarousel(ctx, editCarouselReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 删除轮播图
func DeleteCarousel(ctx *gin.Context) {
	// 获取参数
	id := utils.StringToUint(ctx.Param("id"))

	if err := service.DeleteCarousel(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}
