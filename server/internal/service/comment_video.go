package service

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

// cid 为已解析好的视频数字 ID
func AddVideoComment(ctx *gin.Context, addCommentReq dto.AddCommentReq, cid uint) (vo.AddCommentResp, error) {
    userId := ctx.GetUint("userId")

    // 处理@的用户
    atUserIds := FindUserIdsByName(addCommentReq.At)

    // 保存到数据库
    comment := model.Comment{
        Cid:           cid,
        Uid:           userId,
        Content:       addCommentReq.Content,
        AtUsernames:   strings.Join(addCommentReq.At, ","),
        AtUserIds:     utils.UintJoin(atUserIds, ","),
        ParentId:      addCommentReq.ParentID,
        ReplyUserID:   addCommentReq.ReplyUserID,
        ReplyUserName: addCommentReq.ReplyUserName,
        Type:          global.CONTENT_TYPE_VIDEO,
    }

    if err := global.Mysql.Create(&comment).Error; err != nil {
        utils.ErrorLog("创建评论失败", "comment", err.Error())
        return vo.AddCommentResp{}, errors.New("评论失败")
    }

    // 发送回复通知
    InsertReplyMessage(addCommentReq, cid, comment.ID, userId, global.CONTENT_TYPE_VIDEO)

    // 当有@用户时才创建AT通知
    if len(atUserIds) > 0 {
        newAtMsg := make([]model.AtMessage, len(atUserIds))
        for i := range atUserIds {
            newAtMsg[i].Uid = atUserIds[i]
            newAtMsg[i].Cid = cid
            newAtMsg[i].Sid = userId
            newAtMsg[i].Type = global.CONTENT_TYPE_VIDEO
        }
        if err := global.Mysql.Create(&newAtMsg).Error; err != nil {
            utils.ErrorLog("创建AT信息失败", "comment", err.Error())
        }
    }

    return vo.CommentToCommentResp(comment), nil
}


// 获取视频评论
func GetVideoComment(ctx *gin.Context, vid uint, page, pageSize int) ([]vo.CommentResp, int64, error) {
	var total int64
	var comments []vo.CommentResp

	global.Mysql.Model(&model.Comment{}).Where("cid = ? and `type` = ?", vid, global.CONTENT_TYPE_VIDEO).Count(&total)
	err := global.Mysql.Model(&model.Comment{}).Select(vo.COMMENT_FIELD).
		Joins("LEFT JOIN `comment` AS reply ON `comment`.id = `reply`.parent_id").
		Where("`comment`.parent_id = 0 and `comment`.deleted_at is null and `comment`.cid = ? and `comment`.`type` = ?", vid, global.CONTENT_TYPE_VIDEO).
		Group("`comment`.id").Limit(pageSize).Offset((page - 1) * pageSize).
		Find(&comments).Error
	if err != nil {
		utils.ErrorLog("获取评论失败", "comment", err.Error())
		return comments, total, errors.New("获取失败")
	}

	userId := ctx.GetUint("userId")
	ids := make([]uint, len(comments))
	for i := range comments {
		ids[i] = comments[i].ID
	}
	likedSet := GetCommentLikedSet(userId, ids)
	dislikedSet := GetCommentDislikedSet(userId, ids)

	for i := 0; i < len(comments); i++ {
		comments[i].Author = GetUserBaseInfo(comments[i].Uid)
		comments[i].Reply, _ = FindVideoReplyList(comments[i].ID, 1, 3, userId)
		comments[i].Liked = likedSet[comments[i].ID]
		comments[i].Disliked = dislikedSet[comments[i].ID]
	}

	return comments, total, nil
}

// 获取视频回复
func GetVideoReply(ctx *gin.Context, commentId uint, page, pageSize int) ([]vo.ReplyResp, error) {
	userId := ctx.GetUint("userId")
	return FindVideoReplyList(commentId, page, pageSize, userId)
}

// 删除评论回复
func DeleteVideoComment(ctx *gin.Context, id uint) error {
	comment, err := FindCommentById(id)
	if err != nil {
		utils.ErrorLog("查询评论失败", "comment", err.Error())
		return errors.New("无法获取评论信息")
	}

	video, err := FindVideoById(comment.Cid)
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

	// 收集要删除的评论ID（含子评论），用于级联清理点赞
	var childIds []uint
	global.Mysql.Model(&model.Comment{}).Where("parent_id = ?", id).Pluck("id", &childIds)
	allIds := append(childIds, id)

	// 删除子评论
	if err := global.Mysql.Where("parent_id = ?", id).Delete(&model.Comment{}).Error; err != nil {
		utils.ErrorLog("删除子评论失败", "comment", err.Error())
	}

	if err := global.Mysql.Where("id = ?", id).Delete(&model.Comment{}).Error; err != nil {
		utils.ErrorLog("删除评论失败", "comment", err.Error())
		return errors.New("删除失败")
	}

	// 级联清理点赞记录与点赞消息（孤儿数据）
	CleanupCommentLikes(allIds)

	return nil
}

// 获取视频评论列表
func GetVideoCommentList(ctx *gin.Context, cid uint, page, pageSize int) ([]vo.CommentListResp, int64, error) {
	userId := ctx.GetUint("userId")

	db := global.Mysql.Model(&model.ReplyMessage{})
	if cid != 0 {
		if video, _ := FindVideoById(cid); video.Uid != userId {
			return nil, 0, errors.New("内容不存在")
		}
		db = db.Where("cid = ? and `type` = ?", cid, global.CONTENT_TYPE_VIDEO)
	} else {
		db = db.Where("cid in (?) and `type` = ?", global.Mysql.Model(&model.Video{}).
			Select("id").Where("uid = ?", userId), global.CONTENT_TYPE_VIDEO)
	}

	var total int64
	var comments []vo.CommentListResp
	db.Count(&total)
	db.Limit(pageSize).Offset((page - 1) * pageSize).Scan(&comments)
	for i := 0; i < len(comments); i++ {
		comments[i].Author = GetUserInfo(comments[i].Sid)
		comments[i].TargetUser = GetUserInfo(comments[i].Uid)
		comments[i].Video = GetVideoInfo(comments[i].Cid)
		comments[i].Video.ID = 0
	}

	return comments, total, nil
}

func FindVideoReplyList(commentId uint, page, pageSize int, userId uint) ([]vo.ReplyResp, error) {
	var replies []vo.ReplyResp
	err := global.Mysql.Model(&model.Comment{}).Select(vo.REPLY_FIELD).
		Where("parent_id = ?", commentId).
		Limit(pageSize).Offset((page - 1) * pageSize).Scan(&replies).Error
	if err != nil {
		utils.ErrorLog("获取回复失败", "comment", err.Error())
		return replies, errors.New("获取回复失败")
	}

	ids := make([]uint, len(replies))
	for i := range replies {
		ids[i] = replies[i].ID
	}
	likedSet := GetCommentLikedSet(userId, ids)
	dislikedSet := GetCommentDislikedSet(userId, ids)

	for i := 0; i < len(replies); i++ {
		replies[i].Author = GetUserBaseInfo(replies[i].Uid)
		replies[i].Liked = likedSet[replies[i].ID]
		replies[i].Disliked = dislikedSet[replies[i].ID]
	}

	return replies, nil
}

// 点赞评论
func LikeComment(ctx *gin.Context, commentId uint) error {
	userId := ctx.GetUint("userId")

	comment, err := FindCommentById(commentId)
	if err != nil {
		return errors.New("评论不存在")
	}

	err = global.Mysql.Transaction(func(tx *gorm.DB) error {
		like := model.CommentLike{CommentId: commentId, Uid: userId}
		// 利用唯一索引兜底并发；冲突即"已点赞"
		res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&like)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("已点赞")
		}
		// 原子自增点赞数，避免 read-then-write 丢失更新
		if err := tx.Model(&model.Comment{}).Where("id = ?", commentId).
			Update("likes", gorm.Expr("likes + ?", 1)).Error; err != nil {
			return err
		}
		// 互斥：若已点踩，清除点踩并原子递减点踩数
		delRes := tx.Where("comment_id = ? AND uid = ?", commentId, userId).Delete(&model.CommentDislike{})
		if delRes.Error != nil {
			return delRes.Error
		}
		if delRes.RowsAffected > 0 {
			return tx.Model(&model.Comment{}).
				Where("id = ? AND dislikes > 0", commentId).
				Update("dislikes", gorm.Expr("dislikes - ?", 1)).Error
		}
		return nil
	})
	if err != nil {
		if err.Error() == "已点赞" {
			return err
		}
		utils.ErrorLog("点赞评论失败", "comment_like", err.Error())
		return errors.New("点赞失败")
	}

	// 不给自己发通知；异步且加 panic 守卫
	if comment.Uid != userId {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					utils.ErrorLog("评论点赞通知 panic", "msg_like", "")
				}
			}()
			if err := InsertCommentLikeMessage(userId, commentId, comment.Cid, comment.Uid, comment.Type); err != nil {
				utils.ErrorLog("发送评论点赞消息失败", "msg_like", err.Error())
			}
		}()
	}

	return nil
}

// 取消点赞评论
func UnlikeComment(ctx *gin.Context, commentId uint) error {
	userId := ctx.GetUint("userId")

	if _, err := FindCommentById(commentId); err != nil {
		return errors.New("评论不存在")
	}

	err := global.Mysql.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("comment_id = ? and uid = ?", commentId, userId).Delete(&model.CommentLike{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("未点赞")
		}
		// 原子递减，且通过 likes > 0 防止下溢
		return tx.Model(&model.Comment{}).
			Where("id = ? AND likes > 0", commentId).
			Update("likes", gorm.Expr("likes - ?", 1)).Error
	})
	if err != nil {
		if err.Error() == "未点赞" {
			return err
		}
		utils.ErrorLog("取消点赞失败", "comment_like", err.Error())
		return errors.New("取消点赞失败")
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				utils.ErrorLog("评论点赞通知删除 panic", "msg_like", "")
			}
		}()
		if err := RemoveCommentLikeMessage(commentId, userId); err != nil {
			utils.ErrorLog("删除评论点赞消息失败", "msg_like", err.Error())
		}
	}()

	return nil
}

// 批量查询当前用户对一组评论的点赞状态
func GetCommentLikedSet(userId uint, commentIds []uint) map[uint]bool {
	result := make(map[uint]bool, len(commentIds))
	if userId == 0 || len(commentIds) == 0 {
		return result
	}
	var likedIds []uint
	global.Mysql.Model(&model.CommentLike{}).
		Where("uid = ? AND comment_id IN ?", userId, commentIds).
		Pluck("comment_id", &likedIds)
	for _, id := range likedIds {
		result[id] = true
	}
	return result
}
