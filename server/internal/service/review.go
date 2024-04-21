package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func ReviewVideoApproved(ctx *gin.Context, reviewVideoReq dto.ReviewVideoReq) error {

	tx := global.Mysql.Begin()

	// 把所有待审核的资源改为审核通过
	if err := tx.Model(&model.Resource{}).Where("vid = ?", reviewVideoReq.Vid).Updates(
		map[string]interface{}{"status": global.AUDIT_APPROVED},
	).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("更新资源状态失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	// 更新视频状态为审核通过
	if err := tx.Model(&model.Video{}).Where("id = ?", reviewVideoReq.Vid).Updates(
		map[string]interface{}{"status": global.AUDIT_APPROVED},
	).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("更新视频状态失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	tx.Commit()
	return nil
}

func ReviewVideoFailed(ctx *gin.Context, reviewVideoReq dto.ReviewVideoReq) error {
	tx := global.Mysql.Begin()
	// 更新视频状态为审核不通过
	if err := tx.Model(&model.Video{}).Where("id = ?", reviewVideoReq.Vid).Updates(
		map[string]interface{}{"status": global.REVIEW_FAILED},
	).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("更新视频状态失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	// 添加审核记录
	if err := tx.Create(&model.Review{
		Vid:    reviewVideoReq.Vid,
		Status: reviewVideoReq.Status,
		Remark: reviewVideoReq.Remark,
	}).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("更新视频状态失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	tx.Commit()
	return nil
}

func GetReviewRecord(ctx *gin.Context, videoId uint) (vo.ReviewResp, error) {
	userId := ctx.GetUint("userId")
	video, _ := FindVideoById(videoId)
	if video.Uid != userId {
		return vo.ReviewResp{}, errors.New("视频不存在")
	}

	var review vo.ReviewResp
	global.Mysql.Model(&model.Review{}).Where("vid = ?", videoId).
		Select(vo.REVIEW_FIELD).Last(&review)

	return review, nil
}
