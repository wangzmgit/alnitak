package vo

import "time"

type ResourceResp struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Vid       uint      `json:"vid"`
	Title     string    `json:"title"`
	Duration  float64   `json:"duration"`
	Status    int       `json:"status"`
}
