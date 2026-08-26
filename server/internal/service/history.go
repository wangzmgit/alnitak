package service

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func AddHistory(ctx *gin.Context, historyReq dto.HistoryReq) error {
	userId := ctx.GetUint("userId")
	if historyReq.Part == 0 {
		historyReq.Part = 1
	}

	// 转换 Vid 为字符串
	vidStr := fmt.Sprintf("%v", historyReq.Vid)

	// 解析视频ID（支持 shortId 或数字ID）
	videoId, err := ParseVideoID(vidStr)
	if err != nil {
		return errors.New("视频不存在")
	}

	// 获取视频的 short_id
	var video model.Video
	if err := global.Mysql.Where("id = ?", videoId).First(&video).Error; err != nil || video.ShortID == "" {
		return errors.New("视频不存在")
	}

	// 优先使用前端传入的 rid
	rid := historyReq.Rid
	if rid == "" {
		rid, _ = GetResourceShortIDByPart(videoId, historyReq.Part)
	}

	// 优先通过 resource_short_id 查找（精准匹配，不受排序影响）
	if rid != "" {
		var history model.History
		err := global.Mysql.Where("video_short_id = ? AND uid = ? AND resource_short_id = ?", video.ShortID, userId, rid).First(&history).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			utils.ErrorLog("保存历史记录失败", "history", err.Error())
			return errors.New("保存失败")
		}

		if history.ID == 0 {
			// 不存在，创建新记录
			if err := global.Mysql.Create(&model.History{
				VideoShortID:     video.ShortID,
				Uid:              userId,
				Time:             historyReq.Time,
				Part:             historyReq.Part,
				Duration:         float64(historyReq.Duration),
				ResourceShortID:  rid,
			}).Error; err != nil {
				utils.ErrorLog("保存历史记录失败", "history", err.Error())
				return errors.New("保存失败")
			}
		} else {
			// 已存在，更新记录
			history.Time = historyReq.Time
			history.Part = historyReq.Part
			history.Duration = float64(historyReq.Duration)
			history.ResourceShortID = rid
			if err := global.Mysql.Save(&history).Error; err != nil {
				utils.ErrorLog("保存历史记录失败", "history", err.Error())
				return errors.New("保存失败")
			}
		}
		return nil
	}

	// 回退到旧的按 video_short_id + uid 查询逻辑
	var history model.History
	err = global.Mysql.Where("video_short_id = ? AND uid = ?", video.ShortID, userId).First(&history).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		utils.ErrorLog("保存历史记录失败", "history", err.Error())
		return errors.New("保存失败")
	}

	if history.ID == 0 {
		if err := global.Mysql.Create(&model.History{
			VideoShortID:     video.ShortID,
			Uid:              userId,
			Time:             historyReq.Time,
			Part:             historyReq.Part,
			Duration:         float64(historyReq.Duration),
			ResourceShortID:  rid,
		}).Error; err != nil {
			utils.ErrorLog("保存历史记录失败", "history", err.Error())
			return errors.New("保存失败")
		}
	} else {
		history.Time = historyReq.Time
		history.Part = historyReq.Part
		history.Duration = float64(historyReq.Duration)
		if rid != "" {
			history.ResourceShortID = rid
		}
		if err := global.Mysql.Save(&history).Error; err != nil {
			utils.ErrorLog("保存历史记录失败", "history", err.Error())
			return errors.New("保存失败")
		}
	}

	return nil
}

func GetHistoryList(ctx *gin.Context, page, pageSize int) (videos []vo.HistoryVideoResp, err error) {
	userId := ctx.GetUint("userId")
	subQuery := global.Mysql.Model(&model.History{}).Where("uid = ?", userId).
		Select("video_short_id, MAX(updated_at) as latest_updated_at, MAX(id) as latest_id").
		Group("video_short_id")
	if err := global.Mysql.Model(&model.History{}).Select(vo.HISTORY_VIDEO_FIELD).
		Joins("LEFT JOIN `video` ON `video`.short_id = `history`.video_short_id").
		Joins("INNER JOIN (?) latest on `history`.video_short_id = latest.video_short_id and `history`.updated_at = latest.latest_updated_at and `history`.id = latest.latest_id", subQuery).
		Where("`history`.uid = ? and video.deleted_at is null and video.`status` = ?", userId, global.AUDIT_APPROVED).
		Order("`history`.`updated_at` desc").Limit(pageSize).Offset((page - 1) * pageSize).
		Find(&videos).Error; err != nil {
		utils.ErrorLog("获取历史记录视频失败", "history", err.Error())
		return videos, errors.New("获取失败")
	}

	enrichHistoryPGCMeta(videos)
	return
}

type historyPGCRow struct {
	VideoShortID string `gorm:"column:video_short_id"`
	EpID         uint   `gorm:"column:ep_id"`
}

func enrichHistoryPGCMeta(videos []vo.HistoryVideoResp) {
	if len(videos) == 0 {
		return
	}

	var shortIDs []string
	for _, v := range videos {
		if v.ShortID != "" {
			shortIDs = append(shortIDs, v.ShortID)
		}
	}
	if len(shortIDs) == 0 {
		return
	}

	var pgcRows []historyPGCRow
	global.Mysql.Model(&model.Video{}).Where("short_id IN ? AND pgc_attached = ?", shortIDs, true).
		Pluck("short_id, ep_id", &pgcRows)

	pgcMap := make(map[string]uint, len(pgcRows))
	for _, r := range pgcRows {
		pgcMap[r.VideoShortID] = r.EpID
	}

	for i := range videos {
		if epID, ok := pgcMap[videos[i].ShortID]; ok && epID > 0 {
			videos[i].PGCAttached = true
			videos[i].EpID = epID
		}
	}
}

func GetHistoryProgress(ctx *gin.Context, videoId, part uint) (progress float64, realPart uint, err error) {
	userId := ctx.GetUint("userId")

	// 获取视频的 short_id
	var video model.Video
	if err := global.Mysql.Where("id = ?", videoId).First(&video).Error; err != nil || video.ShortID == "" {
		return 0, 0, errors.New("视频不存在")
	}

	// 优先通过ResourceShortID查找（不受排序影响）
	if part > 0 {
		resourceShortID, _ := GetResourceShortIDByPart(videoId, part)
		if resourceShortID != "" {
			var history model.History
			err = global.Mysql.Where("video_short_id = ? AND uid = ? AND resource_short_id = ?", video.ShortID, userId, resourceShortID).First(&history).Error
			if err == nil {
				return history.Time, history.Part, nil
			}
			if err != gorm.ErrRecordNotFound {
				utils.ErrorLog("通过ResourceShortID获取历史记录进度失败", "history", err.Error())
			}
		}
	}

	// 回退：按旧逻辑查找
	var history model.History
	if part == 0 {
		history, err = FindLatestHistory(video.ShortID, userId)
	} else {
		history, err = FindHistoryByVideoShortID(video.ShortID, userId, part)
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 无历史记录，返回进度0，不视为错误
			return 0, 0, nil
		}
		utils.ErrorLog("获取历史记录进度失败", "history", err.Error())
		return 0, 0, errors.New("获取失败")
	}
	return history.Time, history.Part, nil
}

func FindLatestHistory(videoShortID string, userId uint) (history model.History, err error) {
	if err = global.Mysql.Where("video_short_id = ? and uid = ?", videoShortID, userId).Order("updated_at desc").First(&history).Error; err != nil {
		return
	}

	return
}

func FindHistoryByVideoShortID(videoShortID string, userId uint, part uint) (history model.History, err error) {
	if err = global.Mysql.Where("video_short_id = ? and uid = ? and part = ?", videoShortID, userId, part).First(&history).Error; err != nil {
		return
	}

	return
}

func GetHistoryProgressByRid(ctx *gin.Context, videoId uint, rid string, part uint) (float64, error) {
	userId := ctx.GetUint("userId")

	// 获取视频的 short_id
	var video model.Video
	if err := global.Mysql.Where("id = ?", videoId).First(&video).Error; err != nil || video.ShortID == "" {
		return 0, errors.New("视频不存在")
	}

	// 优先通过 video_short_id + rid 精准匹配
	if rid != "" {
		var history model.History
		err := global.Mysql.Where("video_short_id = ? AND uid = ? AND resource_short_id = ?", video.ShortID, userId, rid).First(&history).Error
		if err == nil {
			return history.Time, nil
		}
		if err != gorm.ErrRecordNotFound {
			utils.ErrorLog("通过video_short_id+rid获取进度失败", "history", err.Error())
		}
	}

	// 回退到按 part 查询（兼容旧数据）
	if part > 0 {
		var historyByPart model.History
		err := global.Mysql.Where("video_short_id = ? AND uid = ? AND part = ?", video.ShortID, userId, part).First(&historyByPart).Error
		if err == nil {
			return historyByPart.Time, nil
		}
		if err != gorm.ErrRecordNotFound {
			utils.ErrorLog("回退到part获取进度失败", "history", err.Error())
		}
	}

	return 0, errors.New("no record")
}