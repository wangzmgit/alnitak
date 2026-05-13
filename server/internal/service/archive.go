package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"interastral-peace.com/alnitak/internal/cache"
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

	var video model.Video
	if err := global.Mysql.Select("shares").First(&video, vid).Error; err == nil {
		archive.Share = video.Shares
	}

	return archive
}

func FindArticleArchiveData(aid uint) vo.ArchiveStatResp {
	var archive vo.ArchiveStatResp

	global.Mysql.Model(&model.LikeArticle{}).Where("aid = ? and is_like = 1", aid).Count(&archive.Like)
	global.Mysql.Model(&model.CollectArticle{}).Where("aid = ? and is_collect = 1", aid).Count(&archive.Collect)

	var article model.Article
	if err := global.Mysql.Select("shares").First(&article, aid).Error; err == nil {
		archive.Share = article.Shares
	}

	return archive
}

// ShareVideo 视频分享计数+1
// 使用 Redis SetNX 做 (userId, videoId) 维度的冷却，防脚本/外挂刷分享数。
// 命中冷却时静默跳过计数（返回 nil），不向前端暴露频率提示，避免影响 UX。
func ShareVideo(ctx *gin.Context, vid string) error {
	videoId, err := ParseVideoID(vid)
	if err != nil {
		return errors.New("视频不存在")
	}
	userId := ctx.GetUint("userId")
	if userId == 0 {
		return errors.New("请先登录")
	}
	if !cache.TryAcquireShareVideoLimit(userId, videoId) {
		return nil
	}
	result := global.Mysql.Model(&model.Video{}).Where("id = ?", videoId).
		UpdateColumn("shares", gorm.Expr("shares + 1"))
	if result.Error != nil {
		return errors.New("分享失败")
	}
	if result.RowsAffected == 0 {
		return errors.New("视频不存在")
	}
	return nil
}

// ShareArticle 文章分享计数+1
// 使用 Redis SetNX 做 (userId, articleId) 维度的冷却，防脚本/外挂刷分享数。
// 命中冷却时静默跳过计数（返回 nil），不向前端暴露频率提示。
func ShareArticle(ctx *gin.Context, aid uint) error {
	userId := ctx.GetUint("userId")
	if userId == 0 {
		return errors.New("请先登录")
	}
	if !cache.TryAcquireShareArticleLimit(userId, aid) {
		return nil
	}
	result := global.Mysql.Model(&model.Article{}).Where("id = ?", aid).
		UpdateColumn("shares", gorm.Expr("shares + 1"))
	if result.Error != nil {
		return errors.New("分享失败")
	}
	if result.RowsAffected == 0 {
		return errors.New("文章不存在")
	}
	return nil
}
