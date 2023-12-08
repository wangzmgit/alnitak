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
