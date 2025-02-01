package vo

import "time"

const (
	PARTITION_FIELD = "`id`,`name`,`type`,`parent_id`,`created_at`"
)

// 分区
type PartitionResp struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Type      int       `json:"type"`
	ParentID  uint      `json:"parentId"`
	CreatedAt time.Time `json:"createdAt"`
}
