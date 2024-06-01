package vo

import (
	"time"
)

const (
	CAROUSEL_FIELD = "`id`,`img`,`url`,`title`,`color`,`use`,`partition_id`,`created_at`"
)

type CarouselResp struct {
	ID          uint      `json:"id"`
	Img         string    `json:"img"`
	Url         string    `json:"url"`
	Title       string    `json:"title"`
	Color       string    `json:"color"`
	Use         bool      `json:"use"`
	PartitionId uint      `json:"partitionId"`
	CreatedAt   time.Time `json:"createdAt"`
}
