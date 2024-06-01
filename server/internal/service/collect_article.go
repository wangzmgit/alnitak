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

func CollectArticle(ctx *gin.Context, collectReq dto.CollectArticleReq) error {
	var err error
	articleId := collectReq.Aid
	userId := ctx.GetUint("userId")

	collect, err := FindCollectArticleByUid(articleId, userId)
	if err != nil && err != gorm.ErrRecordNotFound {
		utils.ErrorLog("获取收藏信息失败", "collect", err.Error())
		return errors.New("收藏失败")
	}

	if collect.IsCollect {
		return errors.New("已经收藏了")
	}

	// 更新数据库
	if collect.ID == 0 {
		err = global.Mysql.Create(&model.CollectArticle{Uid: userId, Aid: articleId, IsCollect: true}).Error
	} else {
		err = global.Mysql.Model(&model.CollectArticle{}).Where("uid = ? and aid = ?", userId, articleId).Updates(
			map[string]interface{}{"is_collect": true},
		).Error
	}
	if err != nil {
		utils.ErrorLog("收藏失败", "collect", err.Error())
		return errors.New("收藏失败")
	}

	return nil
}

func CancelCollectArticle(ctx *gin.Context, collectReq dto.CollectArticleReq) error {
	articleId := collectReq.Aid
	userId := ctx.GetUint("userId")

	collect, err := FindCollectArticleByUid(articleId, userId)
	if err != nil && err != gorm.ErrRecordNotFound {
		utils.ErrorLog("获取收藏信息失败", "collect", err.Error())
		return errors.New("收藏失败")
	}

	if !collect.IsCollect {
		return errors.New("未收藏")
	}

	if err := global.Mysql.Model(&model.CollectArticle{}).Where("uid = ? and aid = ?", userId, articleId).Updates(
		map[string]interface{}{"is_collect": false},
	).Error; err != nil {
		utils.ErrorLog("收藏失败", "collect", err.Error())
		return errors.New("取消收藏失败")
	}

	return nil
}

func HasCollectArticle(ctx *gin.Context, articleId uint) (bool, error) {
	userId := ctx.GetUint("userId")
	collect, err := FindCollectArticleByUid(articleId, userId)
	if err != nil {
		utils.ErrorLog("获取收藏信息失败", "collect", err.Error())
		return false, errors.New("获取失败")
	}

	return collect.IsCollect, nil
}

func FindCollectArticleByUid(articleId, userId uint) (collect model.CollectArticle, err error) {
	err = global.Mysql.Where("`uid` = ? and aid = ?", userId, articleId).First(&collect).Error

	return
}
