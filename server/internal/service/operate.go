package service

import (
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
)

func AddOperate(operate *model.Operate) error {
	return global.Mysql.Create(&operate).Error
}
