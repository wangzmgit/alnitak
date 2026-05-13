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
	global.Mysql.Model(&model.LikeMessage{}).Where("uid = ?", userId).Order("id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Scan(&msg)
	for i := 0; i < len(msg); i++ {
		msg[i].User = GetUserInfo(msg[i].Sid)
		switch msg[i].Type {
		case global.CONTENT_TYPE_VIDEO:
			msg[i].Video = GetVideoInfo(msg[i].Cid)
		case global.CONTENT_TYPE_ARTICLE:
			msg[i].Article = GetArticleInfo(msg[i].Cid)
		case global.CONTENT_TYPE_COMMENT:
			msg[i].Comment = GetCommentInfo(msg[i].CommentId)
		}
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

func InsertCommentLikeMessage(senderId, commentId, cid, targetId uint, parentType int) error {
	return global.Mysql.Create(&model.LikeMessage{
		Sid:        senderId,
		Cid:        cid,
		CommentId:  commentId,
		Uid:        targetId,
		Type:       global.CONTENT_TYPE_COMMENT,
		ParentType: parentType,
	}).Error
}

func RemoveLikeMessage(videoId, senderId uint, contentType int) error {
	return global.Mysql.Where("cid = ? and sid = ? and `type` = ?", videoId, senderId, contentType).
		Delete(&model.LikeMessage{}).Error
}

func RemoveCommentLikeMessage(commentId, senderId uint) error {
	return global.Mysql.Where("comment_id = ? and sid = ? and `type` = ?", commentId, senderId, global.CONTENT_TYPE_COMMENT).
		Delete(&model.LikeMessage{}).Error
}
