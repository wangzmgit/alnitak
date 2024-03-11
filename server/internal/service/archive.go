package service

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
)

func GetArchiveStat(ctx *gin.Context, vid uint) (vo.ArchiveStatResp, error) {
	var archive vo.ArchiveStatResp

	global.Mysql.Model(&model.Like{}).Where("vid = ?", vid).Count(&archive.Like)
	global.Mysql.Model(&model.Collect{}).Distinct("uid").Where("vid = ?", vid).Count(&archive.Collect)

	return archive, nil
}
