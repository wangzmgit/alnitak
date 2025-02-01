package vo

import "time"

const (
	REVIEW_FIELD = "`status`,`remark`,`created_at`"
)

// 消息列表
type ReviewResp struct {
	Status    int       `json:"status"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"createdAt"`
}
