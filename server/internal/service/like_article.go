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

func LikeArticle(ctx *gin.Context, likeReq dto.LikeArticleReq) error {
	var err error
	articleId := likeReq.Aid
	userId := ctx.GetUint("userId")

	like, err := FindLikeArticleByUid(articleId, userId)
	if err != nil && err != gorm.ErrRecordNotFound {
		utils.ErrorLog("获取点赞信息失败", "like", err.Error())
		return errors.New("点赞失败")
	}

	if like.IsLike {
		return errors.New("已经点过赞了")
	}

	// 更新数据库
	if like.ID == 0 {
		err = global.Mysql.Create(&model.LikeArticle{Uid: userId, Aid: articleId, IsLike: true}).Error
	} else {
		err = global.Mysql.Model(&model.LikeArticle{}).Where("uid = ? and aid = ?", userId, articleId).Updates(
			map[string]interface{}{"is_like": true},
		).Error
	}
	if err != nil {
		utils.ErrorLog("点赞失败", "like", err.Error())
		return errors.New("点赞失败")
	}

	// 查询视频作者并添加点赞通知
	article := GetArticleInfo(articleId)
	InsertLikeMessage(userId, article.ID, article.Uid, global.CONTENT_TYPE_ARTICLE)

	return nil
}

func CancelLikeArticle(ctx *gin.Context, likeReq dto.LikeArticleReq) error {
	articleId := likeReq.Aid
	userId := ctx.GetUint("userId")

	like, err := FindLikeArticleByUid(articleId, userId)
	if err != nil && err != gorm.ErrRecordNotFound {
		utils.ErrorLog("获取点赞信息失败", "like", err.Error())
		return errors.New("点赞失败")
	}

	if !like.IsLike {
		return errors.New("未点赞")
	}

	if err := global.Mysql.Model(&model.LikeArticle{}).Where("uid = ? and aid = ?", userId, articleId).Updates(
		map[string]interface{}{"is_like": false},
	).Error; err != nil {
		utils.ErrorLog("点赞失败", "like", err.Error())
		return errors.New("取消点赞失败")
	}

	// 查询视频作者并删除点赞通知
	RemoveLikeMessage(articleId, userId, global.CONTENT_TYPE_ARTICLE)

	return nil
}

func HasLikeArticle(ctx *gin.Context, articleId uint) (bool, error) {
	userId := ctx.GetUint("userId")
	like, err := FindLikeArticleByUid(articleId, userId)
	if err != nil {
		utils.ErrorLog("获取点赞信息失败", "like", err.Error())
		return false, errors.New("获取点赞信息失败")
	}

	return like.IsLike, nil
}

func FindLikeArticleByUid(articleId, userId uint) (like model.LikeArticle, err error) {
	err = global.Mysql.Where("`uid` = ? and aid = ?", userId, articleId).First(&like).Error

	return
}
