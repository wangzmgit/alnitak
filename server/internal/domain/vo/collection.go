package vo

import (
	"time"
)

const (
	COLLECTION_LIST_FIELD = "`id`,`cover`,`name`,`desc`,`open`,`created_at`"
	COLLECTION_FIELD      = "`id`,`uid`,`cover`,`name`,`desc`,`open`,`created_at`"
)

// 收藏夹
type CollectionResp struct {
	ID        uint         `json:"id"`
	Uid       uint         `json:"uid"`
	Cover     string       `json:"cover"`
	Name      string       `json:"name"` //收藏夹名称
	Desc      string       `json:"desc"` //简介
	Open      bool         `json:"open"` //是否公开
	CreatedAt time.Time    `json:"createdAt"`
	Author    UserInfoResp `json:"author" gorm:"-"`
}

type CollectionListResp struct {
	ID        uint      `json:"id"`
	Cover     string    `json:"cover"`
	Name      string    `json:"name"` //收藏夹名称
	Desc      string    `json:"desc"` //简介
	Open      bool      `json:"open"` //是否公开
	CreatedAt time.Time `json:"createdAt"`
}
