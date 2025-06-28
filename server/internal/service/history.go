package service

import (
	"errors"
	"fmt"
	"sort"

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

	history, err := FindHistoryByPart(historyReq.Vid, userId, historyReq.Part)
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
	var all []vo.HistoryVideoResp
	if err := global.Mysql.Model(&model.History{}).Select(vo.HISTORY_VIDEO_FIELD).
		Joins("LEFT JOIN `video` ON `video`.id = `history`.vid").Where("`history`.uid = ?", userId).
		Order("`history`.`updated_at` desc").Find(&all).Error; err != nil {
		utils.ErrorLog("获取历史记录视频失败", "history", err.Error())
		return videos, errors.New("获取失败")
	}
	// 只保留每个视频id下分集中updated_at最大的那条
	vidMap := make(map[uint]vo.HistoryVideoResp)
	for _, v := range all {
		if old, ok := vidMap[v.ID]; !ok || v.UpdatedAt.After(old.UpdatedAt) {
			vidMap[v.ID] = v
		}
	}
	// 分页
	var result []vo.HistoryVideoResp
	for _, v := range vidMap {
		result = append(result, v)
	}
	// 按updated_at倒序
	sort.Slice(result, func(i, j int) bool {
		return result[i].UpdatedAt.After(result[j].UpdatedAt)
	})
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(result) {
		return []vo.HistoryVideoResp{}, nil
	}
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], nil
}

func GetHistoryProgress(ctx *gin.Context, videoId, part uint) (progress float64, realPart uint, err error) {
	userId := ctx.GetUint("userId")
	var history model.History
	if part == 0 {
		err = global.Mysql.Where("vid = ? and uid= ?", videoId, userId).Order("updated_at desc").First(&history).Error
		if err != nil {
			detail := err.Error() + " | videoId=" + fmt.Sprint(videoId) + ", userId=" + fmt.Sprint(userId) + ", part=" + fmt.Sprint(part)
			utils.ErrorLog("获取历史记录进度失败", "history", detail)
			return 0, 0, errors.New("获取失败")
		}
	} else {
		history, err = FindHistoryByPart(videoId, userId, part)
		if err != nil {
			detail := err.Error() + " | videoId=" + fmt.Sprint(videoId) + ", userId=" + fmt.Sprint(userId) + ", part=" + fmt.Sprint(part)
			utils.ErrorLog("获取历史记录进度失败", "history", detail)
			return 0, 0, errors.New("获取失败")
		}
	}
	return history.Time, history.Part, nil
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
