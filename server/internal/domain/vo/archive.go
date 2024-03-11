package vo

// 视频点赞收藏数据
type ArchiveStatResp struct {
	Collect int64 `json:"collect"`
	Like    int64 `json:"like"`
}
