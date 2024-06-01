package dto

type AddCommentReq struct {
	Cid           uint
	Content       string
	At            []string
	ParentID      uint
	ReplyUserID   uint
	ReplyUserName string
	ReplyContent  string
}
