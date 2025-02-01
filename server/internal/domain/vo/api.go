package vo

import "time"

type ApiResp struct {
	ID        uint      `json:"id"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Category  string    `json:"category"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"createdAt"`
}
