package service

import (
	"errors"

	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

var errShortIDExhausted = errors.New("shortId 分配失败，请重试")

const shortIDAllocAttempts = 24

// AllocateUniqueVideoShortID 为 video 表分配全局唯一的随机 shortId。
func AllocateUniqueVideoShortID() (string, error) {
	for range shortIDAllocAttempts {
		id, err := utils.NewOpaqueShortID()
		if err != nil {
			return "", err
		}
		var n int64
		global.Mysql.Model(&model.Video{}).Where("short_id = ?", id).Count(&n)
		if n == 0 {
			return id, nil
		}
	}
	return "", errShortIDExhausted
}

// AllocateUniqueResourceShortID 为 resource 表分配全局唯一的随机 shortId。
func AllocateUniqueResourceShortID() (string, error) {
	for range shortIDAllocAttempts {
		id, err := utils.NewOpaqueShortID()
		if err != nil {
			return "", err
		}
		var n int64
		global.Mysql.Model(&model.Resource{}).Where("short_id = ?", id).Count(&n)
		if n == 0 {
			return id, nil
		}
	}
	return "", errShortIDExhausted
}

// AllocateUniqueArticleShortID 为 article 表分配全局唯一的随机 shortId。
func AllocateUniqueArticleShortID() (string, error) {
	for range shortIDAllocAttempts {
		id, err := utils.NewOpaqueShortID()
		if err != nil {
			return "", err
		}
		var n int64
		global.Mysql.Model(&model.Article{}).Where("short_id = ?", id).Count(&n)
		if n == 0 {
			return id, nil
		}
	}
	return "", errShortIDExhausted
}

// AllocateUniqueSubtitleShortID 为 subtitle_track 表分配全局唯一的随机 shortId。
func AllocateUniqueSubtitleShortID() (string, error) {
	for range shortIDAllocAttempts {
		id, err := utils.NewOpaqueShortID()
		if err != nil {
			return "", err
		}
		var n int64
		global.Mysql.Model(&model.SubtitleTrack{}).Where("short_id = ?", id).Count(&n)
		if n == 0 {
			return id, nil
		}
	}
	return "", errShortIDExhausted
}
