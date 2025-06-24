package service

import (
	"errors"

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

	history, err := FindHistory(historyReq.Vid, userId)
	if err != nil && err != gorm.ErrRecordNotFound {
		utils.ErrorLog("保存历史记录失败", "history", err.Error())
		return errors.New("保存失败")
	}

	if history.ID == 0 {
		if err := global.Mysql.Create(&model.History{
			Uid:  userId,
			Vid:  historyReq.Vid,
			Time: historyReq.Time,
			Part: historyReq.Part,
		}).Error; err != nil {
			utils.ErrorLog("保存历史记录失败", "history", err.Error())
			return errors.New("保存失败")
		}
	} else {
		history.Time = historyReq.Time
		history.Part = historyReq.Part
		if err := global.Mysql.Save(&history).Error; err != nil {
			utils.ErrorLog("保存历史记录失败", "history", err.Error())
			return errors.New("保存失败")
		}
	}

	return nil
}

func GetHistoryList(ctx *gin.Context, page, pageSize int) (videos []vo.HistoryVideoResp, err error) {
	userId := ctx.GetUint("userId")
	if err := global.Mysql.Model(&model.History{}).
		Select("video.id as id, video.uid, video.title, video.cover, video.desc, history.updated_at, history.time, history.part").
		Joins("LEFT JOIN `video` ON `video`.id = `history`.vid").
		Where("`history`.uid = ?", userId).
		Order("`history`.`updated_at` desc").Limit(pageSize).Offset((page - 1) * pageSize).
		Find(&videos).Error; err != nil {
		utils.ErrorLog("获取历史记录视频失败", "history", err.Error())
		return videos, errors.New("获取失败")
	}
	// 批量查分P标题
	vidPartMap := make(map[uint][]model.Resource)
	for i := range videos {
		if _, ok := vidPartMap[videos[i].ID]; !ok {
			var resources []model.Resource
			global.Mysql.Where("vid = ?", videos[i].ID).Order("id asc").Find(&resources)
			vidPartMap[videos[i].ID] = resources
		}
		resources := vidPartMap[videos[i].ID]
		if int(videos[i].Part) > 0 && int(videos[i].Part) <= len(resources) {
			videos[i].PartTitle = resources[videos[i].Part-1].Title
		}
	}

	return
}

func GetHistoryProgress(ctx *gin.Context, videoId, part uint) (progress float64, err error) {
	userId := ctx.GetUint("userId")
	history, err := FindHistoryByPart(videoId, userId, part)
	if err != nil {
		utils.ErrorLog("获取历史记录进度失败", "history", err.Error())
		return 0, errors.New("获取失败")
	}

	return history.Time, nil
}

func FindHistory(videoId, userId uint) (history model.History, err error) {
	if err = global.Mysql.Where("vid = ? and uid= ?", videoId, userId).First(&history).Error; err != nil {
		return
	}

	return
}

func FindHistoryByPart(videoId, userId, part uint) (history model.History, err error) {
	if err = global.Mysql.Where("vid = ? and uid= ? and part = ?", videoId, userId, part).
		First(&history).Error; err != nil {
		return
	}

	return
}
