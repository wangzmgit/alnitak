package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

type UserApi struct{}

func (u *UserApi) GetUserInfo(ctx *gin.Context) {
	userId := ctx.GetUint("userId")
	user := service.GetUserInfo(userId)

	resp.OkWithData(ctx, gin.H{"userInfo": user})
}

// 获取用户基本信息
func (u *UserApi) GetUserBaseInfo(ctx *gin.Context) {
	userId := utils.StringToUint(ctx.Query("userId"))

	user := service.GetUserBaseInfo(userId)

	resp.OkWithData(ctx, gin.H{"userInfo": user})
}
