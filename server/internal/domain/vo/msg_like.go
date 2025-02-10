package vo

import (
	"time"
)

type LikeMessageResp struct {
	ID        uint         `json:"id"`
	Cid       uint         `json:"cid"`
	Sid       uint         `json:"sid"`
	CreatedAt time.Time    `json:"createdAt"`
	User      UserInfoResp `json:"user" gorm:"-"`
	Video     VideoResp    `json:"video" gorm:"-"`
	Article   ArticleResp  `json:"article" gorm:"-"`
	Type      int          `json:"type"`
}
