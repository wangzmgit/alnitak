package vo

import (
	"time"

	"interastral-peace.com/alnitak/internal/domain/model"
)

const (
	COMMENT_LIST_FIELD = "`id`,`vid`,`uid`,`content`,`at_usernames`,`at_user_ids`,`created_at`,`parent_id`"
	COMMENT_FIELD      = "`comment`.`id`,`comment`.`uid`,`comment`.`content`,`comment`.`at_usernames`,`comment`.`at_user_ids`,`comment`.`created_at`,`comment`.`parent_id`,COUNT(`reply`.id) AS reply_count,`comment`.`likes`,`comment`.`dislikes`"
	REPLY_FIELD        = "`id`,`uid`,`content`,`at_usernames`,`at_user_ids`,`created_at`,`parent_id`,`reply_user_id`,`reply_user_name`,`likes`,`dislikes`"
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
	Likes       int          `json:"likes"`
	Liked       bool         `json:"liked" gorm:"-"`
	Dislikes    int          `json:"dislikes"`
	Disliked    bool         `json:"disliked" gorm:"-"`
}

type AddCommentResp struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	Uid           uint      `json:"uid"`
	Content       string    `json:"content"`
	AtUsernames   string    `json:"atUsernames"`
	AtUserIds     string    `json:"atUserIds"`
	ParentId      uint      `json:"parentId"`
	ReplyUserID   uint      `json:"replyUserId"`
	ReplyUserName string    `json:"replyUserName"`
	Likes         int       `json:"likes"`
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
	Likes         int          `json:"likes"`
	Liked         bool         `json:"liked" gorm:"-"`
	Dislikes      int          `json:"dislikes"`
	Disliked      bool         `json:"disliked" gorm:"-"`
}

type CommentListResp struct {
	ID                 uint         `json:"id"`
	Cid                uint         `json:"cid"`
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
	Article            ArticleResp  `json:"article" gorm:"-"`
}

func CommentToCommentResp(comment model.Comment) AddCommentResp {
	return AddCommentResp{
		ID:            comment.ID,
		CreatedAt:     comment.CreatedAt,
		Uid:           comment.Uid,
		Content:       comment.Content,
		AtUsernames:   comment.AtUsernames,
		AtUserIds:     comment.AtUserIds,
		ParentId:      comment.ParentId,
		ReplyUserID:   comment.ReplyUserID,
		ReplyUserName: comment.ReplyUserName,
		Likes:         comment.Likes,
	}
}

type CommentInfoResp struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	Type      int       `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
}
