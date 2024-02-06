package vo

const (
	VIDEO_STATUS_FIELD = "`id`,`title`,`cover`,`desc`,`copyright`,`status`,`partition_id`,`tags`"
)

type VideoStatusResp struct {
	ID          uint   `json:"vid"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	Desc        string `json:"desc"`
	Copyright   bool   `json:"copyright"`
	Status      int    `json:"status"`
	PartitionId uint   `json:"partitionId"`
	Tags        string `json:"tags"`
}
