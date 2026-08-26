package vo

import (
	"time"
)

type LikeMessageResp struct {
	ID         uint         `json:"id"`
	Cid        uint         `json:"cid"`
	CommentId  uint         `json:"commentId"`
	Sid        uint         `json:"sid"`
	CreatedAt  time.Time    `json:"createdAt"`
	User       UserInfoResp `json:"user" gorm:"-"`
	Video      VideoResp    `json:"video" gorm:"-"`
	Article    ArticleResp  `json:"article" gorm:"-"`
	Comment    CommentInfoResp `json:"comment" gorm:"-"`
	Type       int          `json:"type"`
	ParentType int          `json:"parentType"`
}
