package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 通过视频ID获取点赞收藏数据
func GetArchiveStat(ctx *gin.Context) {
	vid := utils.StringToUint(ctx.Query("vid"))

	stat, err := service.GetArchiveStat(ctx, vid)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"stat": stat})
}
