package vo

import (
	"time"
)

type LikeMessageResp struct {
	ID        uint         `json:"id"`
	Vid       uint         `json:"vid"`
	Sid       uint         `json:"sid"`
	CreatedAt time.Time    `json:"created_at"`
	User      UserInfoResp `json:"user" gorm:"-"`
	Video     VideoResp    `json:"video" gorm:"-"`
}
