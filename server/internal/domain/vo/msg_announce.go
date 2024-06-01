package vo

import (
	"time"
)

type AnnounceResp struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
}
