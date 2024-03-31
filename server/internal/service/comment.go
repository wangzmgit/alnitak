package service

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func AddComment(ctx *gin.Context, addCommentReq dto.AddCommentReq) (vo.CommentResp, error) {
	userId := ctx.GetUint("userId")

	// 处理@的用户
	atUserIds := FindUserIdsByName(addCommentReq.At)

	// 保存到数据库
	comment := model.Comment{
		Vid:           addCommentReq.Vid,
		Uid:           userId,
		Content:       addCommentReq.Content,
		AtUsernames:   strings.Join(addCommentReq.At, ","),
		AtUserIds:     utils.UintJoin(atUserIds, ","),
		ParentId:      addCommentReq.ParentID,
		ReplyUserID:   addCommentReq.ReplyUserID,
		ReplyUserName: addCommentReq.ReplyUserName,
	}
	if err := global.Mysql.Create(&comment).Error; err != nil {
		utils.ErrorLog("创建评论失败", "comment", err.Error())
		return vo.CommentResp{}, errors.New("评论失败")
	}

	// 发送回复通知
	InsertReplyMessage(addCommentReq, comment.ID, userId)

	// 处理@通知
	length := len(atUserIds)
	newAtMsg := make([]model.AtMessage, length)
	for i := 0; i < length; i++ {
		newAtMsg[i].Uid = atUserIds[i]
		newAtMsg[i].Vid = addCommentReq.Vid
		newAtMsg[i].Sid = userId
	}
	if err := global.Mysql.Create(&newAtMsg).Error; err != nil {
		utils.ErrorLog("创建AT信息失败失败", "comment", err.Error())
	}

	return vo.CommentToCommentResp(comment), nil
}

// 获取评论
func GetComment(ctx *gin.Context, vid uint, page, pageSize int) ([]vo.CommentResp, int64, error) {
	var total int64
	var comments []vo.CommentResp

	global.Mysql.Model(&model.Comment{}).Where("vid = ?", vid).Count(&total)
	err := global.Mysql.Model(&model.Comment{}).Select(vo.COMMENT_FIELD).
		Joins("LEFT JOIN `comment` AS reply ON `comment`.id = `reply`.parent_id").
		Where("`comment`.parent_id = 0 and `comment`.deleted_at is null").
		Group("`comment`.id").Limit(pageSize).Offset((page - 1) * pageSize).
		Find(&comments).Error
	if err != nil {
		utils.ErrorLog("获取评论失败", "comment", err.Error())
		return comments, total, errors.New("获取失败")
	}

	for i := 0; i < len(comments); i++ {
		comments[i].Author = GetUserBaseInfo(comments[i].Uid)
		comments[i].Reply, _ = FindReplyList(comments[i].ID, 1, 3)
	}

	return comments, total, nil
}

// 获取回复
func GetReply(ctx *gin.Context, commentId uint, page, pageSize int) ([]vo.ReplyResp, error) {
	return FindReplyList(commentId, page, pageSize)
}

// 删除评论回复
func DeleteComment(ctx *gin.Context, id uint) error {
	comment, err := FindCommentById(id)
	if err != nil {
		utils.ErrorLog("查询评论失败", "comment", err.Error())
		return errors.New("无法获取评论信息")
	}

	video, err := FindVideoById(comment.Vid)
	if err != nil {
		utils.ErrorLog("查询视频失败", "comment", err.Error())
		return errors.New("无法获取视频信息")
	}

	//uid为评论作者或视频作者
	userId := ctx.GetUint("userId")
	if comment.Uid != userId && userId != video.Uid {
		return errors.New("评论或回复不存在")
	}

	// 移除评论回复通知
	RemoveReplyMessage(comment.ID)

	if err := global.Mysql.Where("id = ?", id).Delete(&model.Comment{}).Error; err != nil {
		utils.ErrorLog("删除评论失败", "comment", err.Error())
		return errors.New("删除失败")
	}

	return nil
}

func FindCommentById(id uint) (comment model.Comment, err error) {
	err = global.Mysql.First(&comment, id).Error
	return
}

func FindReplyList(commentId uint, page, pageSize int) ([]vo.ReplyResp, error) {
	var replies []vo.ReplyResp
	err := global.Mysql.Model(&model.Comment{}).Select(vo.REPLY_FIELD).
		Where("parent_id = ?", commentId).
		Limit(pageSize).Offset((page - 1) * pageSize).Scan(&replies).Error
	if err != nil {
		utils.ErrorLog("获取回复失败", "comment", err.Error())
		return replies, errors.New("获取回复失败")
	}

	for i := 0; i < len(replies); i++ {
		replies[i].Author = GetUserBaseInfo(replies[i].Uid)
	}

	return replies, nil
}
