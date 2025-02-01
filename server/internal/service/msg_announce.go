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

func GetAnnounce(ctx *gin.Context, page, pageSize int) (total int64, announces []vo.AnnounceResp) {
	global.Mysql.Model(&model.Announce{}).Count(&total)
	global.Mysql.Model(&model.Announce{}).Limit(pageSize).Offset((page - 1) * pageSize).Scan(&announces)
	return
}

func AddAnnounce(ctx *gin.Context, addAnnounceReq dto.AddAnnounceReq) error {
	// 保存到数据库
	return global.Mysql.Create(&model.Announce{
		Title:   addAnnounceReq.Title,
		Content: addAnnounceReq.Content,
		Url:     addAnnounceReq.Url,
	}).Error
}

// 删除公告
func DeleteAnnounce(ctx *gin.Context, id uint) error {
	if err := global.Mysql.Where("id = ?", id).Delete(&model.Announce{}).Error; err != nil {
		utils.ErrorLog("删除公告失败", "announce", err.Error())
		return errors.New("删除角色失败")
	}

	return nil
}
