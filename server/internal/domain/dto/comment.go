package dto

type AddCommentReq struct {
	Vid           uint
	Content       string
	At            []string
	ParentID      uint
	ReplyUserID   uint
	ReplyUserName string
}
