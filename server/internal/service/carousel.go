package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

// 获取轮播图
func GetCarousel(ctx *gin.Context, partitionId uint) (carousels []vo.CarouselResp) {
	global.Mysql.Model(&model.Carousel{}).Select(vo.CAROUSEL_FIELD).
		Where("partition_id = ? and `use` = ?", partitionId, true).Scan(&carousels)

	return
}

// 获取全部轮播图
func GetCarouselList(ctx *gin.Context, carouselListReq dto.CarouselListReq) (carousels []vo.CarouselResp) {
	global.Mysql.Model(&model.Carousel{}).Select(vo.CAROUSEL_FIELD).Limit(carouselListReq.PageSize).
		Offset((carouselListReq.Page - 1) * carouselListReq.PageSize).Scan(&carousels)

	return
}

// 新增轮播图
func AddCarousel(ctx *gin.Context, addCarouselReq dto.AddCarouselReq) error {
	userId := ctx.GetUint("userId")

	if cache.GetUploadImage(addCarouselReq.Img) != userId {
		return errors.New("文件链接无效")
	}

	// 保存到数据库
	if err := global.Mysql.Create(&model.Carousel{
		Uid:         userId,
		Img:         addCarouselReq.Img,
		Title:       addCarouselReq.Title,
		Color:       addCarouselReq.Color,
		Url:         addCarouselReq.Url,
		Use:         addCarouselReq.Use,
		PartitionId: addCarouselReq.PartitionId,
	}).Error; err != nil {
		utils.ErrorLog("创建轮播图失败", "carousel", err.Error())
		return errors.New("创建失败")
	}

	return nil
}

// 编辑轮播图信息
func EditCarousel(ctx *gin.Context, editCarouselReq dto.EditCarouselReq) error {
	userId := ctx.GetUint("userId")

	oldCarousel, _ := FindCarouselById(editCarouselReq.ID)
	if cache.GetUploadImage(editCarouselReq.Img) != userId {
		// 查询是否与旧封面图一致
		if oldCarousel.Img != editCarouselReq.Img {
			return errors.New("文件链接无效")
		}
	}

	// 保存到数据库
	if err := global.Mysql.Model(&model.Carousel{}).Where("id = ?", editCarouselReq.ID).Updates(
		map[string]interface{}{
			"uid":          userId,
			"img":          editCarouselReq.Img,
			"title":        editCarouselReq.Title,
			"color":        editCarouselReq.Color,
			"url":          editCarouselReq.Url,
			"use":          editCarouselReq.Use,
			"partition_id": editCarouselReq.PartitionId,
		},
	).Error; err != nil {
		utils.ErrorLog("修改轮播图失败", "carousel", err.Error())
		return errors.New("修改失败")
	}

	return nil
}

// 删除API
func DeleteCarousel(ctx *gin.Context, id uint) error {
	if err := global.Mysql.Where("id = ?", id).Delete(&model.Carousel{}).Error; err != nil {
		utils.ErrorLog("删除轮播图失败", "api", err.Error())
		return errors.New("删除失败")
	}

	return nil
}

// 通过视频ID查询视频
func FindCarouselById(id uint) (carousel model.Carousel, err error) {
	err = global.Mysql.Where("`id` = ?", id).First(&carousel).Error
	return
}
