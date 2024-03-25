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

	history, err := FindHistory(historyReq.Vid, userId, historyReq.Part)
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
		if err := global.Mysql.Save(&history).Error; err != nil {
			utils.ErrorLog("保存历史记录失败", "history", err.Error())
			return errors.New("保存失败")
		}
	}

	return nil
}

func GetHistoryList(ctx *gin.Context, page, pageSize int) (videos []vo.VideoResp, err error) {
	userId := ctx.GetUint("userId")
	videoIds := global.Mysql.Model(&model.History{}).Select("vid").Where("userId = ?", userId).
		Limit(pageSize).Offset((page - 1) * pageSize)

	if err := global.Mysql.Model(&model.Video{}).Select(vo.VIDEO_FIELD).
		Where("id in (?)", videoIds).Scan(&videos).Error; err != nil {
		utils.ErrorLog("获取历史记录视频失败", "history", err.Error())
		return videos, errors.New("获取失败")
	}

	return
}

func GetHistoryProgress(ctx *gin.Context, videoId, part uint) (progress float64, err error) {
	userId := ctx.GetUint("userId")
	history, err := FindHistory(videoId, userId, part)
	if err != nil {
		utils.ErrorLog("获取历史记录进度失败", "history", err.Error())
		return 0, errors.New("获取失败")
	}

	return history.Time, nil
}

func FindHistory(videoId, userId, part uint) (history model.History, err error) {
	if err = global.Mysql.Where("vid = ? and uid= ? and part = ?", videoId, userId, part).
		First(&history).Error; err != nil {
		return
	}

	return
}
