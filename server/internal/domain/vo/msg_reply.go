package vo

import (
	"time"
)

type ReplyMessageResp struct {
	ID                 uint         `json:"id"`
	Cid                uint         `json:"cid"`
	Sid                uint         `json:"sid"`
	CreatedAt          time.Time    `json:"createdAt"`
	Content            string       `json:"content"`
	TargetReplyContent string       `json:"targetReplyContent"`
	RootContent        string       `json:"rootContent"`
	CommentId          string       `json:"commentId"`
	User               UserInfoResp `json:"user" gorm:"-"`
	Video              VideoResp    `json:"video" gorm:"-"`
	Article            ArticleResp  `json:"article" gorm:"-"`
	Type               int          `json:"type"`
}
