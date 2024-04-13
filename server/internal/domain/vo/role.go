package vo

import "time"

type RoleResp struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Desc      string    `json:"desc"`
	HomePage  string    `json:"homePage"`
	CreatedAt time.Time `json:"createdAt"`
}

type AllRoleResp struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
