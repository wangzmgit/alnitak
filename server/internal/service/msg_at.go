package service

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
)

func GetAtMessage(ctx *gin.Context, page, pageSize int) (total int64, msg []vo.AtMessageResp) {
	userId := ctx.GetUint("userId")

	global.Mysql.Model(&model.AtMessage{}).Where("uid = ?", userId).Count(&total)
	global.Mysql.Model(&model.AtMessage{}).Where("uid = ?", userId).Limit(pageSize).Offset((page - 1) * pageSize).Scan(&msg)
	for i := 0; i < len(msg); i++ {
		msg[i].User = GetUserInfo(msg[i].Sid)
		if msg[i].Type == global.CONTENT_TYPE_VIDEO {
			msg[i].Video = GetVideoInfo(msg[i].Cid)
		} else {
			msg[i].Article = GetArticleInfo(msg[i].Cid)
		}
	}

	return
}

func InsertAtMessage(senderId, videoId, targetId uint, contentType int) error {
	return global.Mysql.Create(&model.AtMessage{
		Sid:  senderId,
		Cid:  videoId,
		Uid:  targetId,
		Type: contentType,
	}).Error
}

func RemoveAtMessage(videoId, senderId uint) error {
	return global.Mysql.Where("vid = ? and sid = ?", videoId, senderId).Delete(&model.AtMessage{}).Error
}
