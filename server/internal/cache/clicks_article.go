package cache

import (
	"errors"
	"strconv"

	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func GetArticleClicksLimit(articleId uint, ip string) string {
	s := global.Redis.Get(ARTICLE_CLICKS_LIMIT_KEY + utils.UintToString(articleId) + ":" + ip)

	return s
}

func SetArticleClicksLimit(articleId uint, ip string) {
	global.Redis.Set(ARTICLE_CLICKS_LIMIT_KEY+utils.UintToString(articleId)+":"+ip, 1,
		ARTICLE_CLICKS_LIMIT_EXPRIRATION_TIME)
}

func GetArticleClicks(articleId uint) (int64, error) {
	s := global.Redis.Get(ARTICLE_CLICKS_KEY + utils.UintToString(articleId))
	if s == "" {
		return 0, errors.New("no record")
	}

	count, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		utils.ErrorLog("阅读量转换为int64类型失败", "cache", err.Error())
		return 0, err
	}

	return count, nil
}

func SetArticleClicks(articleId uint, count int64) {
	global.Redis.Set(ARTICLE_CLICKS_KEY+utils.UintToString(articleId), count, ARTICLE_CLICKS_EXPRIRATION_TIME)
}

func GetArticleClicksKeys() []string {
	return global.Redis.Keys(ARTICLE_CLICKS_KEY + "*")
}

func AddArticleClicks(articleId uint) {
	global.Redis.Incr(ARTICLE_CLICKS_KEY + utils.UintToString(articleId))
}
