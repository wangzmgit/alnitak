package vo

// 分区
type PartitionResp struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	ParentID uint   `json:"parentId"`
}
