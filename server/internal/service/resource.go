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

func ModifyResourceTitle(ctx *gin.Context, modifyTitleReq dto.ModifyResourceTitleReq) error {
	var resource model.Resource
	userId := ctx.GetUint("userId")
	global.Mysql.Model(&model.Resource{}).Where("id = ? and uid = ?", modifyTitleReq.ID, userId).First(&resource)
	if resource.ID == 0 {
		return errors.New("资源不存在")
	}

	if err := global.Mysql.Model(&model.Resource{}).Where("id = ?", modifyTitleReq.ID).Updates(
		map[string]interface{}{
			"title": modifyTitleReq.Title,
		},
	).Error; err != nil {
		utils.ErrorLog("更新资源标题失败", "resource", err.Error())
		return errors.New("更新资源标题失败")
	}
	return nil
}

// 删除资源
func DeleteResource(ctx *gin.Context, id uint) error {
	var resource model.Resource
	userId := ctx.GetUint("userId")
	global.Mysql.Model(&model.Resource{}).Where("id = ? and uid = ?", id, userId).First(&resource)
	if resource.ID == 0 {
		return errors.New("资源不存在")
	}

	var count int64
	global.Mysql.Model(&model.Resource{}).Where("vid = ?", resource.Vid).Count(&count)
	if count <= 1 {
		return errors.New("至少需要一个视频")
	}

	if err := global.Mysql.Where("id = ?", id).Delete(&model.Resource{}).Error; err != nil {
		utils.ErrorLog("删除资源失败", "resource", err.Error())
		return errors.New("删除资源失败")
	}

	// 更新视频信息缓存
	cache.DelVideoInfo(resource.Vid)
	VideoWriteCache(resource.Vid)

	return nil
}

// 获取视频资源
func GetVideoResourceByStatus(videoId uint, status int) (resources []vo.ResourceResp) {
	global.Mysql.Model(&model.Resource{}).Where("vid = ? and status = ?", videoId, status).Scan(&resources)

	return
}

// 获取视频资源
func GetReviewResourceList(videoId uint) (resources []vo.ResourceResp) {
	global.Mysql.Model(&model.Resource{}).Where("vid = ?", videoId).Scan(&resources)

	return
}

// 获取视频资源支持的分辨率信息
func GetResourceQuality(ctx *gin.Context, id uint) ([]string, error) {
	if !IsResourceExist(id) {
		return nil, errors.New("资源不存在")
	}

	var quality []string
	if err := global.Mysql.Model(&model.VideoIndexFile{}).Where("resource_id = ?", id).
		Pluck("quality", &quality).Error; err != nil {
		utils.ErrorLog("分辨率信息获取失败", "resource", err.Error())
		return nil, errors.New("获取失败")
	}

	return quality, nil
}

// 获取视频资源支持的分辨率信息(后台管理)
func GetResourceQualityManage(ctx *gin.Context, id uint) ([]string, error) {
	var quality []string
	if err := global.Mysql.Model(&model.VideoIndexFile{}).Where("resource_id = ?", id).
		Pluck("quality", &quality).Error; err != nil {
		utils.ErrorLog("分辨率信息获取失败", "resource", err.Error())
		return nil, errors.New("获取失败")
	}

	return quality, nil
}

func IsResourceExist(resourceId uint) bool {
	var resource model.Resource
	global.Mysql.Model(&model.Resource{}).Select("id").
		Where("id = ? and `status` = ?", resourceId, global.AUDIT_APPROVED).First(&resource)

	return resource.ID != 0
}
