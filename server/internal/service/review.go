package service

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
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

	// 校验视频状态：只有 WAITING_REVIEW（首次审核）或 AUDIT_APPROVED（新增分P）才允许通过
	// 使用 FOR UPDATE 行锁防止并发审核导致状态不一致
	var curVideo model.Video
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", reviewVideoReq.Vid).First(&curVideo).Error; err != nil {
		tx.Rollback()
		return errors.New("视频不存在")
	}
	if curVideo.Status != global.WAITING_REVIEW && curVideo.Status != global.AUDIT_APPROVED {
		tx.Rollback()
		return errors.New("当前状态不允许审核操作")
	}

	// 只把待审核的资源改为审核通过（不影响转码中、已通过等其他状态的资源）
	// 同时设为对外可见（处理中的分P仍为隐藏，等转码完成后再改为可见）
	if err := tx.Model(&model.Resource{}).
		Where("vid = ? AND status = ?", reviewVideoReq.Vid, global.WAITING_REVIEW).
		Updates(map[string]interface{}{
			"status":         global.AUDIT_APPROVED,
			"visible_status": global.VISIBLE_SHOWN,
		}).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("更新资源状态失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	// 计算总时长时排除即将被替换的旧资源，避免 duration 在替换清理后虚增
	var totalSec int64
	totalSecQuery := tx.Model(&model.Resource{}).Where("vid = ? AND status = ?", reviewVideoReq.Vid, global.AUDIT_APPROVED)

	// 收集即将被清理的旧资源 ID（这些资源被替换后会被硬删除）
	var replacedIDs []uint
	tx.Model(&model.Resource{}).Where("vid = ? AND replace_id > 0 AND status = ?", reviewVideoReq.Vid, global.AUDIT_APPROVED).
		Pluck("replace_id", &replacedIDs)
	if len(replacedIDs) > 0 {
		totalSecQuery = totalSecQuery.Where("id NOT IN ?", replacedIDs)
	}
	totalSecQuery.Select("COALESCE(SUM(duration), 0)").Scan(&totalSec)

	// 更新视频状态为审核通过
	if err := tx.Model(&model.Video{}).Where("id = ?", reviewVideoReq.Vid).Updates(map[string]interface{}{
		"status":   global.AUDIT_APPROVED,
		"duration": int(totalSec),
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

	// 处理替换资源：刚审核通过的替换资源，将旧资源的 ShortID 转移到新资源
	// （审核期间旧资源保留 ShortID 保证弹幕/字幕/历史正常，现在转移给新资源）
	type replaceCleanupTask struct {
		ReplaceID uint
		FileID    uint
		Uid       uint
		DirName   string
	}
	var cleanupTasks []replaceCleanupTask

	var replacementResources []model.Resource
	tx.Model(&model.Resource{}).Where("vid = ? AND replace_id > 0 AND status = ?", reviewVideoReq.Vid, global.AUDIT_APPROVED).Find(&replacementResources)
	for _, r := range replacementResources {
		// 查出旧资源，获取其 ShortID
		var oldRes model.Resource
		if err := tx.Where("id = ?", r.ReplaceID).First(&oldRes).Error; err != nil {
			continue
		}
		oldShortID := oldRes.ShortID

		// 先清空旧资源的 ShortID（释放 unique 约束）
		if oldShortID != "" {
			tx.Model(&model.Resource{}).Where("id = ?", r.ReplaceID).Update("short_id", "")
			// 将新资源设为旧资源的 ShortID（弹幕/字幕/历史绑定跟随 ShortID 自动继承）
			tx.Model(&model.Resource{}).Where("id = ?", r.ID).Update("short_id", oldShortID)
		} else {
			// oldShortID 为空有两种情况：
			// 1. 旧资源没有 ShortID（不会发生，创建时一定分配）
			// 2. 已有另一个替换资源先处理了此旧资源，清空了 ShortID
			//    此时新资源保留自己原有的 ShortID，不做覆盖
			utils.InfoLog(fmt.Sprintf("【审核通过-替换】旧资源ShortID已为空，新ResourceID=%d 保留自身ShortID",
				r.ID), "review")
		}

		// 隐藏旧资源（后续提交事务后硬删除并清理文件）
		tx.Model(&model.Resource{}).Where("id = ?", r.ReplaceID).Update("visible_status", global.VISIBLE_HIDDEN)
		utils.InfoLog(fmt.Sprintf("【审核通过-替换】新ResourceID=%d→ShortID=%s, 旧ResourceID=%d 隐藏",
			r.ID, oldShortID, r.ReplaceID), "review")

		// 收集旧资源信息用于事务提交后的硬删除和文件清理
		var oldIndexFile model.VideoIndexFile
		tx.Where("resource_id = ?", r.ReplaceID).First(&oldIndexFile)
		cleanupTasks = append(cleanupTasks, replaceCleanupTask{
			ReplaceID: r.ReplaceID,
			FileID:    oldRes.FileID,
			Uid:       oldRes.Uid,
			DirName:   oldIndexFile.DirName,
		})
	}

	// 先提交事务
	if err := tx.Commit().Error; err != nil {
		utils.ErrorLog("提交事务失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	// 事务成功后再操作缓存
	video, _ := FindVideoById(reviewVideoReq.Vid)
	// 先加入分区列表，再清除详情缓存，避免两操作间的空隙视频丢失
	if !video.PGCAttached {
		cache.SetVideoId(global.VideoPartitionMap[video.PartitionId], video.ID)
	}
	cache.DelVideoInfo(reviewVideoReq.Vid)

	// 硬删除旧替换资源并清理文件（事务已提交，ShortID 已转移，安全清理）
	for _, task := range cleanupTasks {
		// 先减引用再删记录：decreaseVideoFileRefCount 失败时资源记录还在，可重试
		decreaseVideoFileRefCount(task.FileID, task.Uid, task.ReplaceID, task.DirName)

		// 删除旧资源的 VideoIndexFile
		if err := global.Mysql.Where("resource_id = ?", task.ReplaceID).Delete(&model.VideoIndexFile{}).Error; err != nil {
			utils.ErrorLog(fmt.Sprintf("【替换清理】删除VideoIndexFile失败 resourceID=%d", task.ReplaceID), "review", err.Error())
		}

		// 硬删除旧资源记录
		if err := global.Mysql.Unscoped().Delete(&model.Resource{}, task.ReplaceID).Error; err != nil {
			utils.ErrorLog(fmt.Sprintf("【替换清理】硬删除旧资源失败 resourceID=%d", task.ReplaceID), "review", err.Error())
		}
	}

	return nil
}

// 审核不通过(视频)
func ReviewVideoFailed(ctx *gin.Context, reviewVideoReq dto.ReviewVideoReq) error {
	tx := global.Mysql.Begin()

	// 校验视频状态
	// 使用 FOR UPDATE 行锁防止并发审核导致状态不一致
	var curVideo model.Video
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", reviewVideoReq.Vid).First(&curVideo).Error; err != nil {
		tx.Rollback()
		return errors.New("视频不存在")
	}
	if curVideo.Status != global.WAITING_REVIEW && curVideo.Status != global.AUDIT_APPROVED {
		tx.Rollback()
		return errors.New("当前状态不允许审核操作")
	}

	// 判断是否已有已审核通过的分P（替换资源驳回场景）
	var approvedCount int64
	tx.Model(&model.Resource{}).Where("vid = ? AND status = ?", reviewVideoReq.Vid, global.AUDIT_APPROVED).Count(&approvedCount)

	if approvedCount > 0 {
		// 替换资源驳回：视频已有通过的分P，只驳回待审核的分P，不动视频整体状态
		// （如替换资源 replace_id > 0，但也有可能新增分P一并被驳回）
		if err := tx.Model(&model.Resource{}).
			Where("vid = ? AND status = ?", reviewVideoReq.Vid, global.WAITING_REVIEW).
			Update("status", global.REVIEW_FAILED).Error; err != nil {
			tx.Rollback()
			utils.ErrorLog("更新资源状态失败", "review", err.Error())
			return errors.New("更新状态失败")
		}

		// 如果指定了冲突稿件，只更新被驳回的资源
		if reviewVideoReq.ConflictResourceID != 0 {
			if err := tx.Model(&model.Resource{}).
				Where("vid = ? AND status = ?", reviewVideoReq.Vid, global.REVIEW_FAILED).
				Updates(map[string]interface{}{
					"conflict_resource_id": reviewVideoReq.ConflictResourceID,
					"conflict_reason":      reviewVideoReq.ConflictReason,
				}).Error; err != nil {
				tx.Rollback()
				utils.ErrorLog("更新冲突关联失败", "review", err.Error())
				return errors.New("更新冲突关联失败")
			}
			utils.InfoLog(fmt.Sprintf("【审核驳回-替换】仅驳回待审核分P vid=%d -> conflictResourceId=%d", reviewVideoReq.Vid, reviewVideoReq.ConflictResourceID), "review")
		}
	} else {
		// 首次上传驳回：视频无任何已通过的分P，驳回整个视频
		if err := tx.Model(&model.Video{}).Where("id = ?", reviewVideoReq.Vid).Updates(
			map[string]interface{}{"status": global.REVIEW_FAILED},
		).Error; err != nil {
			tx.Rollback()
			utils.ErrorLog("更新视频状态失败", "review", err.Error())
			return errors.New("更新状态失败")
		}

		// 同时把所有待审核的分P设为驳回（避免资源残留 WAITING_REVIEW 成为孤儿数据）
		if err := tx.Model(&model.Resource{}).
			Where("vid = ? AND status = ?", reviewVideoReq.Vid, global.WAITING_REVIEW).
			Update("status", global.REVIEW_FAILED).Error; err != nil {
			tx.Rollback()
			utils.ErrorLog("更新资源状态失败", "review", err.Error())
			return errors.New("更新状态失败")
		}

		// 如果指定了冲突稿件，更新该视频下所有资源的冲突关联
		if reviewVideoReq.ConflictResourceID != 0 {
			if err := tx.Model(&model.Resource{}).Where("vid = ?", reviewVideoReq.Vid).Updates(
				map[string]interface{}{
					"conflict_resource_id": reviewVideoReq.ConflictResourceID,
					"conflict_reason":      reviewVideoReq.ConflictReason,
				},
			).Error; err != nil {
				tx.Rollback()
				utils.ErrorLog("更新冲突关联失败", "review", err.Error())
				return errors.New("更新冲突关联失败")
			}
			utils.InfoLog(fmt.Sprintf("【审核驳回-首次】设置冲突关联 vid=%d -> conflictResourceId=%d", reviewVideoReq.Vid, reviewVideoReq.ConflictResourceID), "review")
		}
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

	if err := tx.Commit().Error; err != nil {
		utils.ErrorLog("提交事务失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	// 清除视频信息缓存和视频ID缓存
	video, _ := FindVideoById(reviewVideoReq.Vid)
	if video.ID != 0 {
		cache.DelVideoInfo(reviewVideoReq.Vid)
		cache.DelVideoId(video.PartitionId, reviewVideoReq.Vid)
	}

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

	if err := tx.Commit().Error; err != nil {
		utils.ErrorLog("提交事务失败", "review", err.Error())
		return errors.New("更新状态失败")
	}
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

	if err := tx.Commit().Error; err != nil {
		utils.ErrorLog("提交事务失败", "review", err.Error())
		return errors.New("更新状态失败")
	}
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

// 获取待审核合集列表
func GetReviewPlaylistList(page, pageSize int) (total int64, list []vo.PlaylistResp) {
	global.Mysql.Model(&model.Playlist{}).Where("status = ?", global.WAITING_REVIEW).Count(&total)
	global.Mysql.Model(&model.Playlist{}).Select(vo.PLAYLIST_FIELD).
		Where("status = ?", global.WAITING_REVIEW).
		Order("created_at asc").
		Limit(pageSize).Offset((page - 1) * pageSize).Scan(&list)

	// 填充作者信息
	for i := range list {
		list[i].Author = GetUserBaseInfo(list[i].Uid)
	}
	return
}

// 审核通过(合集)
func ReviewPlaylistApproved(ctx *gin.Context, req dto.ReviewPlaylistReq) error {
	tx := global.Mysql.Begin()
	if err := tx.Model(&model.Playlist{}).Where("id = ?", req.ID).Updates(
		map[string]interface{}{"status": global.AUDIT_APPROVED},
	).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("更新合集状态失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	if err := tx.Create(&model.Review{
		Cid:    req.ID,
		Status: global.AUDIT_APPROVED,
		Remark: "",
		Uid:    ctx.GetUint("userId"),
		Type:   global.CONTENT_TYPE_PLAYLIST,
	}).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("创建审核记录失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	if err := tx.Commit().Error; err != nil {
		utils.ErrorLog("提交事务失败", "review", err.Error())
		return errors.New("更新状态失败")
	}
	return nil
}

// 审核不通过(合集)
func ReviewPlaylistFailed(ctx *gin.Context, req dto.ReviewPlaylistReq) error {
	tx := global.Mysql.Begin()
	if err := tx.Model(&model.Playlist{}).Where("id = ?", req.ID).Updates(
		map[string]interface{}{"status": global.REVIEW_FAILED},
	).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("更新合集状态失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	if err := tx.Create(&model.Review{
		Cid:    req.ID,
		Status: req.Status,
		Remark: req.Remark,
		Uid:    ctx.GetUint("userId"),
		Type:   global.CONTENT_TYPE_PLAYLIST,
	}).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("创建审核记录失败", "review", err.Error())
		return errors.New("更新状态失败")
	}

	if err := tx.Commit().Error; err != nil {
		utils.ErrorLog("提交事务失败", "review", err.Error())
		return errors.New("更新状态失败")
	}
	return nil
}

// 获取审核记录(合集)
func GetPlaylistReviewRecord(ctx *gin.Context, playlistId uint) (vo.ReviewResp, error) {
	userId := ctx.GetUint("userId")
	var playlist model.Playlist
	if err := global.Mysql.First(&playlist, playlistId).Error; err != nil {
		return vo.ReviewResp{}, errors.New("合集不存在")
	}
	if playlist.Uid != userId {
		return vo.ReviewResp{}, errors.New("合集不存在")
	}

	var review vo.ReviewResp
	global.Mysql.Model(&model.Review{}).Select(vo.REVIEW_FIELD).
		Where("cid = ? and `type` = ?", playlistId, global.CONTENT_TYPE_PLAYLIST).Last(&review)

	return review, nil
}
