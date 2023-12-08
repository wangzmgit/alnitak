package vo

import "time"

type UserInfoResp struct {
	ID        uint      `json:"uid"`
	Username  string    `json:"name"`
	Sign      string    `json:"sign"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Avatar    string    `json:"avatar"`
	Gender    int       `json:"gender"`
	Status    uint      `json:"status"`
	Birthday  time.Time `json:"birthday"`
	CreatedAt time.Time `json:"createdAt"`
}
