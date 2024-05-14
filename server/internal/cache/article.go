package cache

import (
	"interastral-peace.com/alnitak/internal/global"
)

func SetArticleId(articleId uint) {
	global.Redis.SAdd(ALL_ARTICLE_KEY, articleId)
}

func DelArticleId(articleId uint) {
	global.Redis.SRem(ALL_ARTICLE_KEY, articleId)
}

func DelAllArticleId() {
	global.Redis.Del(ALL_ARTICLE_KEY)
}

func GetRandomArticleIds(count int64) []string {
	return global.Redis.SRandMemberN(ALL_ARTICLE_KEY, count)
}
