package service

import (
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
)

func FindCommentById(id uint) (comment model.Comment, err error) {
	err = global.Mysql.First(&comment, id).Error
	return
}
