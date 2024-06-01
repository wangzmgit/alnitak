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

// 审核通过(视频)
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

	// 统计视频时长并存入数据库
	var duration float64
	tx.Model(&model.Resource{}).Where("vid = ?", reviewVideoReq.Vid).Pluck("SUM(duration) as duration", &duration)

	// 更新视频状态为审核通过
	if err := tx.Model(&model.Video{}).Where("id = ?", reviewVideoReq.Vid).Updates(map[string]interface{}{
		"status":   global.AUDIT_APPROVED,
		"duration": duration,
	}).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("更新视频状态失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	// 添加审核记录
	if err := tx.Create(&model.Review{
		Cid:    reviewVideoReq.Vid,
		Status: global.AUDIT_APPROVED,
		Remark: "",
		Uid:    ctx.GetUint("userId"),
		Type:   global.CONTENT_TYPE_VIDEO,
	}).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("更新视频状态失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	// 视频ID添加到redis中
	video, _ := FindVideoById(reviewVideoReq.Vid)
	cache.SetVideoId(global.VideoPartitionMap[video.PartitionId], video.ID)

	tx.Commit()

	return nil
}

// 审核不通过(视频)
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
		Cid:    reviewVideoReq.Vid,
		Status: reviewVideoReq.Status,
		Remark: reviewVideoReq.Remark,
		Uid:    ctx.GetUint("userId"),
		Type:   global.CONTENT_TYPE_VIDEO,
	}).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("更新视频状态失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	tx.Commit()
	return nil
}

// 获取审核记录(视频)
func GetVideoReviewRecord(ctx *gin.Context, videoId uint) (vo.ReviewResp, error) {
	userId := ctx.GetUint("userId")
	video, _ := FindVideoById(videoId)
	if video.Uid != userId {
		return vo.ReviewResp{}, errors.New("视频不存在")
	}

	var review vo.ReviewResp
	global.Mysql.Model(&model.Review{}).Select(vo.REVIEW_FIELD).
		Where("cid = ? and `type` = ?", videoId, global.CONTENT_TYPE_VIDEO).Last(&review)

	return review, nil
}

// 审核通过(文章)
func ReviewArticleApproved(ctx *gin.Context, reviewArticleReq dto.ReviewArticleReq) error {
	tx := global.Mysql.Begin()

	// 更新文章状态为审核通过
	if err := tx.Model(&model.Article{}).Where("id = ?", reviewArticleReq.Aid).Updates(
		map[string]interface{}{"status": global.AUDIT_APPROVED},
	).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("更新文章状态失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	// 添加审核记录
	if err := tx.Create(&model.Review{
		Cid:    reviewArticleReq.Aid,
		Status: global.AUDIT_APPROVED,
		Remark: "",
		Uid:    ctx.GetUint("userId"),
		Type:   global.CONTENT_TYPE_ARTICLE,
	}).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("更新视频状态失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	// 文章ID添加到redis中
	cache.SetArticleId(reviewArticleReq.Aid)

	tx.Commit()
	return nil
}

// 审核不通过(文章)
func ReviewArticleFailed(ctx *gin.Context, reviewArticleReq dto.ReviewArticleReq) error {
	tx := global.Mysql.Begin()
	// 更新文章状态为审核不通过
	if err := tx.Model(&model.Article{}).Where("id = ?", reviewArticleReq.Aid).Updates(
		map[string]interface{}{"status": global.REVIEW_FAILED},
	).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("更新文章状态失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	// 添加审核记录
	if err := tx.Create(&model.Review{
		Cid:    reviewArticleReq.Aid,
		Status: reviewArticleReq.Status,
		Remark: reviewArticleReq.Remark,
		Uid:    ctx.GetUint("userId"),
		Type:   global.CONTENT_TYPE_ARTICLE,
	}).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("更新文章状态失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	tx.Commit()
	return nil
}

// 获取审核记录(文章)
func GetArticleReviewRecord(ctx *gin.Context, articleId uint) (vo.ReviewResp, error) {
	userId := ctx.GetUint("userId")
	article, _ := FindArticleById(articleId)
	if article.Uid != userId {
		return vo.ReviewResp{}, errors.New("内容不存在")
	}

	var review vo.ReviewResp
	global.Mysql.Model(&model.Review{}).Select(vo.REVIEW_FIELD).
		Where("cid = ? and `type` = ?", articleId, global.CONTENT_TYPE_ARTICLE).Last(&review)

	return review, nil
}
