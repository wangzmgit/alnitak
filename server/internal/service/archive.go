package service

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
)

func GetVideoArchiveStat(ctx *gin.Context, vid uint) (vo.ArchiveStatResp, error) {
	archive := FindVideoArchiveData(vid)

	return archive, nil
}

func GetArticleArchiveStat(ctx *gin.Context, vid uint) (vo.ArchiveStatResp, error) {
	archive := FindArticleArchiveData(vid)

	return archive, nil
}

func FindVideoArchiveData(vid uint) vo.ArchiveStatResp {
	var archive vo.ArchiveStatResp

	global.Mysql.Model(&model.LikeVideo{}).Where("vid = ? and is_like = 1", vid).Count(&archive.Like)
	global.Mysql.Model(&model.CollectVideo{}).Distinct("uid").Where("vid = ?", vid).Count(&archive.Collect)

	return archive
}

func FindArticleArchiveData(aid uint) vo.ArchiveStatResp {
	var archive vo.ArchiveStatResp

	global.Mysql.Model(&model.LikeArticle{}).Where("aid = ? and is_like = 1", aid).Count(&archive.Like)
	global.Mysql.Model(&model.CollectArticle{}).Where("aid = ? and is_collect = 1", aid).Count(&archive.Collect)

	return archive
}
