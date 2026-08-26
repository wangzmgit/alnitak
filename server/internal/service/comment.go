package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func FindCommentById(id uint) (comment model.Comment, err error) {
	err = global.Mysql.First(&comment, id).Error
	return
}

func GetCommentInfo(commentId uint) (commentResp vo.CommentInfoResp) {
	comment, err := FindCommentById(commentId)
	if err != nil {
		return
	}
	commentResp = vo.CommentInfoResp{
		ID:        comment.ID,
		Content:   comment.Content,
		Type:      comment.Type,
		CreatedAt: comment.CreatedAt,
	}
	return
}

// CleanupCommentLikes 清理一组评论的点赞/点踩记录与点赞消息（用于删除评论时级联）
func CleanupCommentLikes(commentIds []uint) {
	if len(commentIds) == 0 {
		return
	}
	if err := global.Mysql.Where("comment_id IN ?", commentIds).Delete(&model.CommentLike{}).Error; err != nil {
		utils.ErrorLog("清理评论点赞记录失败", "comment_like", err.Error())
	}
	if err := global.Mysql.Where("comment_id IN ?", commentIds).Delete(&model.CommentDislike{}).Error; err != nil {
		utils.ErrorLog("清理评论点踩记录失败", "comment_dislike", err.Error())
	}
	if err := global.Mysql.
		Where("comment_id IN ? AND `type` = ?", commentIds, global.CONTENT_TYPE_COMMENT).
		Delete(&model.LikeMessage{}).Error; err != nil {
		utils.ErrorLog("清理评论点赞消息失败", "msg_like", err.Error())
	}
}

// GetCommentDislikedSet 批量查询当前用户对一组评论的点踩状态
func GetCommentDislikedSet(userId uint, commentIds []uint) map[uint]bool {
	result := make(map[uint]bool, len(commentIds))
	if userId == 0 || len(commentIds) == 0 {
		return result
	}
	var dislikedIds []uint
	global.Mysql.Model(&model.CommentDislike{}).
		Where("uid = ? AND comment_id IN ?", userId, commentIds).
		Pluck("comment_id", &dislikedIds)
	for _, id := range dislikedIds {
		result[id] = true
	}
	return result
}

// 点踩评论（与点赞互斥）
func DislikeComment(ctx *gin.Context, commentId uint) error {
	userId := ctx.GetUint("userId")

	if _, err := FindCommentById(commentId); err != nil {
		return errors.New("评论不存在")
	}

	err := global.Mysql.Transaction(func(tx *gorm.DB) error {
		dislike := model.CommentDislike{CommentId: commentId, Uid: userId}
		res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&dislike)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("已点踩")
		}
		if err := tx.Model(&model.Comment{}).Where("id = ?", commentId).
			Update("dislikes", gorm.Expr("dislikes + ?", 1)).Error; err != nil {
			return err
		}
		// 互斥：若已点赞，清除点赞并原子递减点赞数 + 删除点赞消息
		delRes := tx.Where("comment_id = ? AND uid = ?", commentId, userId).Delete(&model.CommentLike{})
		if delRes.Error != nil {
			return delRes.Error
		}
		if delRes.RowsAffected > 0 {
			if err := tx.Model(&model.Comment{}).
				Where("id = ? AND likes > 0", commentId).
				Update("likes", gorm.Expr("likes - ?", 1)).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		if err.Error() == "已点踩" {
			return err
		}
		utils.ErrorLog("点踩评论失败", "comment_dislike", err.Error())
		return errors.New("点踩失败")
	}

	// 互斥后删除历史点赞通知（异步 + panic 守卫）
	go func() {
		defer func() {
			if r := recover(); r != nil {
				utils.ErrorLog("互斥删除点赞通知 panic", "msg_like", "")
			}
		}()
		if err := RemoveCommentLikeMessage(commentId, userId); err != nil {
			utils.ErrorLog("删除评论点赞消息失败", "msg_like", err.Error())
		}
	}()

	return nil
}

// 取消点踩评论
func UndislikeComment(ctx *gin.Context, commentId uint) error {
	userId := ctx.GetUint("userId")

	if _, err := FindCommentById(commentId); err != nil {
		return errors.New("评论不存在")
	}

	err := global.Mysql.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("comment_id = ? and uid = ?", commentId, userId).Delete(&model.CommentDislike{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("未点踩")
		}
		return tx.Model(&model.Comment{}).
			Where("id = ? AND dislikes > 0", commentId).
			Update("dislikes", gorm.Expr("dislikes - ?", 1)).Error
	})
	if err != nil {
		if err.Error() == "未点踩" {
			return err
		}
		utils.ErrorLog("取消点踩失败", "comment_dislike", err.Error())
		return errors.New("取消点踩失败")
	}

	return nil
}

// CleanupContentComments 删除内容（视频/文章）时级联清理所有评论相关数据
// contentType: global.CONTENT_TYPE_VIDEO 或 global.CONTENT_TYPE_ARTICLE
func CleanupContentComments(contentId uint, contentType int) {
	// 先把所有评论ID收集起来，用于清理点赞与点赞消息
	var commentIds []uint
	global.Mysql.Model(&model.Comment{}).
		Where("cid = ? AND `type` = ?", contentId, contentType).
		Pluck("id", &commentIds)

	// 删评论（软删除）
	if err := global.Mysql.
		Where("cid = ? AND `type` = ?", contentId, contentType).
		Delete(&model.Comment{}).Error; err != nil {
		utils.ErrorLog("删除内容关联评论失败", "comment", err.Error())
	}

	// 级联清理：comment_like + msg_like(type=3)
	CleanupCommentLikes(commentIds)

	// 清理内容自身的点赞消息（msg_like，type 对应内容类型）
	if err := global.Mysql.
		Where("cid = ? AND `type` = ?", contentId, contentType).
		Delete(&model.LikeMessage{}).Error; err != nil {
		utils.ErrorLog("清理内容点赞消息失败", "msg_like", err.Error())
	}

	// 清理回复消息
	if err := global.Mysql.
		Where("cid = ? AND `type` = ?", contentId, contentType).
		Delete(&model.ReplyMessage{}).Error; err != nil {
		utils.ErrorLog("清理回复消息失败", "msg_reply", err.Error())
	}
}
