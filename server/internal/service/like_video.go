package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func LikeVideo(ctx *gin.Context, likeReq dto.LikeVideoReq) error {
	var err error
	videoId := likeReq.Vid
	userId := ctx.GetUint("userId")

	like, err := FindLikeVideoByUid(videoId, userId)
	if err != nil && err != gorm.ErrRecordNotFound {
		utils.ErrorLog("获取点赞信息失败", "like", err.Error())
		return errors.New("点赞失败")
	}

	if like.IsLike {
		return errors.New("已经点过赞了")
	}

	// 更新数据库
	if like.ID == 0 {
		err = global.Mysql.Create(&model.LikeVideo{Uid: userId, Vid: videoId, IsLike: true}).Error
	} else {
		err = global.Mysql.Model(&model.LikeVideo{}).Where("uid = ? and vid = ?", userId, videoId).Updates(
			map[string]interface{}{"is_like": true},
		).Error
	}
	if err != nil {
		utils.ErrorLog("点赞失败", "like", err.Error())
		return errors.New("点赞失败")
	}

	// 查询视频作者并添加点赞通知
	video := GetVideoInfo(videoId)
	if userId != video.Uid {
		InsertLikeMessage(userId, video.ID, video.Uid, global.CONTENT_TYPE_VIDEO)
	}

	return nil
}

func CancelLikeVideo(ctx *gin.Context, likeReq dto.LikeVideoReq) error {
	videoId := likeReq.Vid
	userId := ctx.GetUint("userId")

	like, err := FindLikeVideoByUid(videoId, userId)
	if err != nil && err != gorm.ErrRecordNotFound {
		utils.ErrorLog("获取点赞信息失败", "like", err.Error())
		return errors.New("点赞失败")
	}

	if !like.IsLike {
		return errors.New("未点赞")
	}

	if err := global.Mysql.Model(&model.LikeVideo{}).Where("uid = ? and vid = ?", userId, videoId).Updates(
		map[string]interface{}{"is_like": false},
	).Error; err != nil {
		utils.ErrorLog("点赞失败", "like", err.Error())
		return errors.New("取消点赞失败")
	}

	// 查询视频作者并删除点赞通知
	RemoveLikeMessage(videoId, userId, global.CONTENT_TYPE_VIDEO)

	return nil
}

func HasLikeVideo(ctx *gin.Context, videoId uint) (bool, error) {
	userId := ctx.GetUint("userId")
	like, err := FindLikeVideoByUid(videoId, userId)
	if err != nil {
		utils.ErrorLog("获取点赞信息失败", "like", err.Error())
		return false, errors.New("获取点赞信息失败")
	}

	return like.IsLike, nil
}

func FindLikeVideoByUid(videoId, userId uint) (like model.LikeVideo, err error) {
	err = global.Mysql.Where("`uid` = ? and vid = ?", userId, videoId).First(&like).Error

	return
}
