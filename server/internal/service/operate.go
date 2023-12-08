package service

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func AddOkOperation(ctx *gin.Context, position, msg string, param gin.H) {
	operate := model.Operate{
		Position: position,
		Success:  true,
		Msg:      msg,
		Ip:       ctx.ClientIP(),
	}

	// 处理参数
	if param != nil {
		operate.Param = utils.MapToJson(param)
	}

	// 用户ID
	userId := ctx.GetUint("userId")
	if userId != 0 {
		operate.UserID = userId
	}

	AddOperate(&operate)
}

func AddFailOperation(ctx *gin.Context, position, msg string, param gin.H, err error) {
	operate := model.Operate{
		Position: position,
		Success:  false,
		Msg:      msg,
		Ip:       ctx.ClientIP(),
	}

	// 处理参数
	if param != nil {
		operate.Param = utils.MapToJson(param)
	}

	// 用户ID
	userId := ctx.GetUint("userId")
	if userId != 0 {
		operate.UserID = userId
	}

	// 错误信息
	if err != nil {
		operate.Error = err.Error()
	}

	AddOperate(&operate)
}

func AddOperate(operate *model.Operate) error {
	return global.Mysql.Create(&operate).Error
}
