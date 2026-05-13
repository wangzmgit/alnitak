package dto

// Cid 兼容前端传入 shortId（如 "BV1xx"）或数字 ID；handler 统一转换为数字 ID 再进 service
type AddCommentReq struct {
	Cid           string   `json:"cid" form:"cid"`
	Content       string   `json:"content" form:"content"`
	At            []string `json:"at" form:"at"`
	ParentID      uint     `json:"parentID" form:"parentID"`
	ReplyUserID   uint     `json:"replyUserID" form:"replyUserID"`
	ReplyUserName string   `json:"replyUserName" form:"replyUserName"`
	ReplyContent  string   `json:"replyContent" form:"replyContent"`
}
