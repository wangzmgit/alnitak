package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 通过视频ID获取点赞收藏数据
func GetVideoArchiveStat(ctx *gin.Context) {
	vid := utils.StringToUint(ctx.Query("vid"))

	stat, err := service.GetVideoArchiveStat(ctx, vid)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"stat": stat})
}

// 通过文章ID获取点赞收藏数据
func GetArticleArchiveStat(ctx *gin.Context) {
	aid := utils.StringToUint(ctx.Query("aid"))

	stat, err := service.GetArticleArchiveStat(ctx, aid)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"stat": stat})
}
