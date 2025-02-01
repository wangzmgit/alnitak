package vo

import (
	"time"

	"interastral-peace.com/alnitak/internal/domain/model"
)

type ResourceResp struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Vid       uint      `json:"vid"`
	Title     string    `json:"title"`
	Duration  float64   `json:"duration"`
	Status    int       `json:"status"`
}

func ResourceToResourceResp(resource model.Resource) ResourceResp {
	return ResourceResp{
		ID:        resource.ID,
		CreatedAt: resource.CreatedAt,
		Vid:       resource.Vid,
		Title:     resource.Title,
		Duration:  resource.Duration,
		Status:    resource.Status,
	}
}
