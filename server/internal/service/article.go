package service

import (
	"errors"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func UploadArticleInfo(ctx *gin.Context, uploadArticleReq dto.UploadArticleReq) error {
	userId := ctx.GetUint("userId")
	if uploadArticleReq.Cover != "" && cache.GetUploadImage(uploadArticleReq.Cover) != userId {
		return errors.New("文件链接无效")
	}

	if !IsSubpartition(uploadArticleReq.PartitionId, global.CONTENT_TYPE_ARTICLE) {
		return errors.New("分区不存在")
	}

	if err := global.Mysql.Model(&model.Article{}).Create(&model.Article{
		Title:       uploadArticleReq.Title,
		Cover:       uploadArticleReq.Cover,
		Copyright:   uploadArticleReq.Copyright,
		Tags:        uploadArticleReq.Tags,
		Uid:         userId,
		Content:     uploadArticleReq.Content,
		ContentDesc: getContentDsec(uploadArticleReq.Content),
		PartitionId: uploadArticleReq.PartitionId, //分区ID
		Status:      global.WAITING_REVIEW,
	}).Error; err != nil {
		utils.ErrorLog("添加文章失败", "article", err.Error())
		return errors.New("保存失败")
	}

	return nil
}

func EditArticleInfo(ctx *gin.Context, editArticleReq dto.EditArticleReq) error {
	userId := ctx.GetUint("userId")
	if editArticleReq.Cover != "" && cache.GetUploadImage(editArticleReq.Cover) != userId {
		// 查询是否与旧封面图一致
		if v, _ := FindArticleById(editArticleReq.Aid); v.Cover != editArticleReq.Cover {
			return errors.New("文件链接无效")
		}
	}

	if err := global.Mysql.Model(&model.Article{}).Where("id = ?", editArticleReq.Aid).Updates(
		map[string]interface{}{
			"title":        editArticleReq.Title,
			"cover":        editArticleReq.Cover,
			"tags":         editArticleReq.Tags,
			"content_desc": getContentDsec(editArticleReq.Content),
			"content":      editArticleReq.Content,
			"status":       global.WAITING_REVIEW,
		},
	).Error; err != nil {
		utils.ErrorLog("修改文章失败", "article", err.Error())
		return errors.New("修改失败")
	}

	// 删除缓存中的文章ID信息
	cache.DelArticleId(editArticleReq.Aid)

	return nil
}

// 获取自己上传的文章
func GetUploadArticleList(ctx *gin.Context, page, pageSize int) (total int64, articles []vo.UploadArticleResp) {
	userId := ctx.GetUint("userId")

	global.Mysql.Model(&model.Article{}).Where("uid = ?", userId).Count(&total)
	global.Mysql.Model(&model.Article{}).Select(vo.UPLOAD_ARTICLE_FIELD).
		Where("uid = ?", userId).Limit(pageSize).Offset((page - 1) * pageSize).Scan(&articles)

	// 更新点击量数据
	for i := 0; i < len(articles); i++ {
		articles[i].Clicks += GetArticleClicks(articles[i].ID)
	}

	return
}

func GetArticleStatus(ctx *gin.Context, aid uint) (article vo.ArticleStatusResp, err error) {
	userId := ctx.GetUint("userId")
	global.Mysql.Model(&model.Article{}).Select(vo.ARTICLE_STATUS_FIELD).Where("id = ? and uid = ?", aid, userId).Scan(&article)
	if article.ID == 0 {
		return article, errors.New("内容不存在")
	}

	return article, nil
}

// 获文章信息
func GetArticleById(ctx *gin.Context, articleId uint) (vo.ArticleResp, error) {
	article := GetArticleInfo(articleId)
	if article.ID == 0 {
		return article, errors.New("信息不存在")
	}

	// 增加播放量(一个ip在同一个文章下，每30分钟可重新增加1点击量)
	AddArticleClicks(articleId, ctx.ClientIP())
	article.Clicks += GetArticleClicks(article.ID)

	return article, nil
}

// 删除文章
func DeleteArticle(ctx *gin.Context, id uint) error {
	var article model.Article
	userId := ctx.GetUint("userId")
	global.Mysql.Model(&model.Article{}).Where("id = ? and uid = ?", id, userId).First(&article)
	if article.ID == 0 {
		return errors.New("内容不存在")
	}

	if err := global.Mysql.Where("id = ?", id).Delete(&model.Article{}).Error; err != nil {
		utils.ErrorLog("删除文章失败", "article", err.Error())
		return errors.New("删除失败")
	}

	// 删除缓存中的文章ID信息
	cache.DelArticleId(article.ID)

	return nil
}

// 获取所有的文章列表
func GetAllArticleList(ctx *gin.Context) (articles []vo.AllArticleResp) {
	userId := ctx.GetUint("userId")
	global.Mysql.Model(&model.Article{}).Select("`id`,`title`").Where("uid = ?", userId).Scan(&articles)

	return
}

// 获取用户文章
func GetArticleByUser(ctx *gin.Context, userId uint, page, pageSize int) (total int64, articles []vo.UploadArticleResp) {
	global.Mysql.Model(&model.Article{}).
		Where("uid = ? and status = ?", userId, global.AUDIT_APPROVED).Count(&total)
	global.Mysql.Model(&model.Article{}).Select(vo.UPLOAD_ARTICLE_FIELD).
		Where("uid = ? and status = ?", userId, global.AUDIT_APPROVED).
		Limit(pageSize).Offset((page - 1) * pageSize).Scan(&articles)

	// 更新点击量数据
	for i := 0; i < len(articles); i++ {
		articles[i].Clicks += GetVideoClicks(articles[i].ID)
	}

	return
}

// 获取随机文章
func GetRandomArticleList(ctx *gin.Context, size int) []vo.ArticleResp {
	articleIds := cache.GetRandomArticleIds(int64(size))

	len := len(articleIds)
	articles := make([]vo.ArticleResp, len)
	for i := 0; i < len; i++ {
		id := utils.StringToUint(articleIds[i])
		if id == 0 {
			continue
		}
		articles[i] = GetArticleItemInfo(id)
		// 同步播放量
		articles[i].Clicks += GetVideoClicks(id)
	}

	return articles
}

// 获取待审核文章列表
func GetReviewArticleList(reviewListReq dto.ReviewArticleListReq) (total int64, articles []vo.ReviewArticleListResp) {
	global.Mysql.Model(&model.Article{}).Where("status = ?", global.WAITING_REVIEW).Count(&total)
	global.Mysql.Model(&model.Article{}).Where("status = ?", global.WAITING_REVIEW).
		Limit(reviewListReq.PageSize).Offset((reviewListReq.Page - 1) * reviewListReq.PageSize).Scan(&articles)

	// 更新播放量和作者数据
	for i := 0; i < len(articles); i++ {
		articles[i].Author = GetUserBaseInfo(articles[i].Uid)
	}

	return
}

// 获取文章信息
func GetArticleInfo(articleId uint) (article vo.ArticleResp) {
	global.Mysql.Model(&model.Article{}).Select(vo.ARTICLE_FIELD).
		Where("id = ? and status = ?", articleId, global.AUDIT_APPROVED).Scan(&article)
	if article.ID == 0 {
		utils.ErrorLog("信息不存在", "article", "")
		return
	}

	// 获取作者信息
	article.Author = GetUserBaseInfo(article.Uid)

	return
}

// 获取文章信息
func GetArticleItemInfo(articleId uint) (article vo.ArticleResp) {
	global.Mysql.Model(&model.Article{}).Select(vo.ARTICLE_LIST_FIELD).
		Where("id = ? and status = ?", articleId, global.AUDIT_APPROVED).Scan(&article)
	if article.ID == 0 {
		utils.ErrorLog("信息不存在", "article", "")
		return
	}

	// 获取作者信息
	article.Author = GetUserBaseInfo(article.Uid)

	return
}

// 获取文章列表(后台管理)
func GetArticleListManage(articleListReq dto.ArticleListReq) (total int64, articles []vo.ArticleResp) {
	global.Mysql.Model(&model.Article{}).Where("status = ?", global.AUDIT_APPROVED).Count(&total)
	global.Mysql.Model(&model.Article{}).Where("status = ?", global.AUDIT_APPROVED).
		Limit(articleListReq.PageSize).Offset((articleListReq.Page - 1) * articleListReq.PageSize).Scan(&articles)

	// 更新播放量和作者数据
	for i := 0; i < len(articles); i++ {
		articles[i].Clicks += GetArticleClicks(articles[i].ID)
		articles[i].Author = GetUserBaseInfo(articles[i].Uid)
	}

	return
}

// 删除文章(后台管理)
func DeleteArticleManage(ctx *gin.Context, id uint) error {
	if err := global.Mysql.Where("id = ?", id).Delete(&model.Article{}).Error; err != nil {
		utils.ErrorLog("删除文章失败", "article", err.Error())
		return errors.New("删除失败")
	}

	return nil
}

// 通过文章ID查询文章
func FindArticleById(id uint) (article model.Article, err error) {
	err = global.Mysql.Where("`id` = ?", id).First(&article).Error
	return
}

func getContentDsec(c string) string {

	if utf8.RuneCountInString(c) > 200 {
		runes := []rune(c)
		if len(runes) > 200 {
			return string(runes[:200])
		}
	}

	return c
}
