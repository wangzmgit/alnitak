package service

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
)

func GetReplyMessage(ctx *gin.Context, page, pageSize int) (total int64, msg []vo.ReplyMessageResp) {
	userId := ctx.GetUint("userId")

	global.Mysql.Model(&model.ReplyMessage{}).Where("uid = ? and sid != ?", userId, userId).Count(&total)
	global.Mysql.Model(&model.ReplyMessage{}).Where("uid = ? and sid != ?", userId, userId).Limit(pageSize).Offset((page - 1) * pageSize).Scan(&msg)
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

// 添加回复消息
func InsertReplyMessage(addCommentReq dto.AddCommentReq, commentId, userId uint, contentType int) error {
	msg := model.ReplyMessage{
		Cid:                addCommentReq.Cid,
		Sid:                userId,
		CommentId:          commentId,
		Content:            addCommentReq.Content,
		TargetReplyContent: addCommentReq.ReplyContent,
		Type:               contentType,
	}

	if addCommentReq.ParentID != 0 {
		rootComment, _ := FindCommentById(addCommentReq.ParentID)
		msg.CommentId = rootComment.ID
		msg.RootContent = rootComment.Content
		// 通知给评论作者
		msg.Uid = rootComment.Uid
	}

	if addCommentReq.ParentID == 0 {
		if contentType == global.CONTENT_TYPE_VIDEO {
			// 通知给视频作者
			video := GetVideoInfo(addCommentReq.Cid)
			msg.Uid = video.Uid
		} else {
			// 通知给文章作者
			article := GetArticleInfo(addCommentReq.Cid)
			msg.Uid = article.Uid
		}

	}

	if addCommentReq.ReplyUserID != 0 {
		// 通知给回复目标
		msg.Uid = addCommentReq.ReplyUserID
	}

	return global.Mysql.Create(&msg).Error
}

func RemoveReplyMessage(commentId uint) error {
	return global.Mysql.Where("comment_id = ?", commentId).Delete(&model.ReplyMessage{}).Error
}
