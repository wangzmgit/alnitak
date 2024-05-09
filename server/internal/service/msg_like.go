package service

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
)

func GetLikeMessage(ctx *gin.Context, page, pageSize int) (total int64, msg []vo.LikeMessageResp) {
	userId := ctx.GetUint("userId")

	global.Mysql.Model(&model.LikeMessage{}).Where("uid = ?", userId).Count(&total)
	global.Mysql.Model(&model.LikeMessage{}).Where("uid = ?", userId).Limit(pageSize).Offset((page - 1) * pageSize).Scan(&msg)
	for i := 0; i < len(msg); i++ {
		msg[i].User = GetUserInfo(msg[i].Sid)
		msg[i].Video = GetVideoInfo(msg[i].Vid)
	}

	return
}

func InsertLikeMessage(senderId, cid, targetId uint, contentType int) error {
	return global.Mysql.Create(&model.LikeMessage{
		Sid:  senderId,
		Cid:  cid,
		Uid:  targetId,
		Type: contentType,
	}).Error
}

func RemoveLikeMessage(videoId, senderId uint, contentType int) error {
	return global.Mysql.Where("vid = ? and sid = ? and `type` = ?", videoId, senderId, contentType).
		Delete(&model.LikeMessage{}).Error
}
