package vo

import "time"

const (
	USER_BASE_INFO_FIELD = "`id`,`username`,`sign`,`avatar`,`gender`"
)

type UserInfoResp struct {
	ID         uint      `json:"uid"`
	Username   string    `json:"name"`
	Sign       string    `json:"sign"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Avatar     string    `json:"avatar"`
	Gender     int       `json:"gender"`
	SpaceCover string    `json:"spaceCover"`
	Birthday   time.Time `json:"birthday"`
	CreatedAt  time.Time `json:"createdAt"`
}

type UserBaseInfoResp struct {
	ID       uint   `json:"uid"`
	Username string `json:"name"`
	Sign     string `json:"sign"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
}

type UserInfoManageResp struct {
	ID         uint      `json:"uid"`
	Username   string    `json:"name"`
	Sign       string    `json:"sign"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Avatar     string    `json:"avatar"`
	Gender     int       `json:"gender"`
	SpaceCover string    `json:"spaceCover"`
	Birthday   time.Time `json:"birthday"`
	CreatedAt  time.Time `json:"createdAt"`
	Role       string    `json:"role"`
}
