package vo

import (
	"time"

	"interastral-peace.com/alnitak/internal/domain/model"
)

const (
	COMMENT_LIST_FIELD = "`id`,`vid`,`uid`,`content`,`at_usernames`,`at_user_ids`,`created_at`,`parent_id`"
	COMMENT_FIELD      = "`comment`.`id`,`comment`.`uid`,`comment`.`content`,`comment`.`at_usernames`,`comment`.`at_user_ids`,`comment`.`created_at`,`comment`.`parent_id`,COUNT(`reply`.id) AS reply_count"
	REPLY_FIELD        = "`id`,`uid`,`content`,`at_usernames`,`at_user_ids`,`created_at`,`parent_id`,`reply_user_id`,`reply_user_name`"
)

type CommentResp struct {
	ID          uint         `json:"id"`
	CreatedAt   time.Time    `json:"createdAt"`
	Uid         uint         `json:"uid"`
	Content     string       `json:"content"`
	AtUsernames string       `json:"atUsernames"`
	AtUserIds   string       `json:"atUserIds"`
	ParentId    uint         `json:"parentId"`
	Author      UserInfoResp `json:"author" gorm:"-"`
	Reply       []ReplyResp  `json:"reply" gorm:"-"`
	ReplyCount  int64        `json:"replyCount"`
}

type ReplyResp struct {
	ID            uint         `json:"id"`
	CreatedAt     time.Time    `json:"createdAt"`
	Uid           uint         `json:"uid"`
	Content       string       `json:"content"`
	AtUsernames   string       `json:"atUsernames"`
	AtUserIds     string       `json:"atUserIds"`
	ReplyUserID   uint         `json:"replyUserId"`
	ReplyUserName string       `json:"replyUserName"`
	ParentId      uint         `json:"parentId"`
	Author        UserInfoResp `json:"author" gorm:"-"`
}

type CommentListResp struct {
	ID                 uint         `json:"id"`
	Vid                uint         `json:"vid"`
	Sid                uint         `json:"sid"`
	Uid                uint         `json:"uid"`
	CreatedAt          time.Time    `json:"createdAt"`
	Content            string       `json:"content"`
	TargetReplyContent string       `json:"targetReplyContent"`
	RootContent        string       `json:"rootContent"`
	CommentId          string       `json:"commentId"`
	Author             UserInfoResp `json:"author" gorm:"-"` // 作者
	TargetUser         UserInfoResp `json:"target" gorm:"-"` // 回复目标
	Video              VideoResp    `json:"video" gorm:"-"`
}

func CommentToCommentResp(comment model.Comment) CommentResp {
	return CommentResp{
		ID:          comment.ID,
		CreatedAt:   comment.CreatedAt,
		Uid:         comment.Uid,
		Content:     comment.Content,
		AtUsernames: comment.AtUsernames,
		AtUserIds:   comment.AtUserIds,
		ParentId:    comment.ParentId,
	}
}
