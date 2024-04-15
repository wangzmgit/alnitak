package vo

import "time"

const (
	VIDEO_FIELD        = "`id`,`uid`,`title`,`cover`,`desc`,`created_at`,`copyright`,`tags`"
	VIDEO_STATUS_FIELD = "`id`,`title`,`cover`,`desc`,`copyright`,`status`,`partition_id`,`tags`"
	UPLOAD_VIDEO_FIELD = "`id`,`title`,`cover`,`desc`,`copyright`,`created_at`"
)

type VideoResp struct {
	ID        uint           `json:"vid"`
	Uid       uint           `json:"uid"`
	Title     string         `json:"title"`
	Cover     string         `json:"cover"`
	Desc      string         `json:"desc"`
	CreatedAt time.Time      `json:"createdAt"`
	Copyright bool           `json:"copyright"`
	Tags      string         `json:"tags"`
	Clicks    int64          `json:"clicks" gorm:"-"`
	Author    UserInfoResp   `json:"author" gorm:"-"`
	Resources []ResourceResp `json:"resources" gorm:"-"`
}

type VideoStatusResp struct {
	ID          uint           `json:"vid"`
	Title       string         `json:"title"`
	Cover       string         `json:"cover"`
	Desc        string         `json:"desc"`
	Copyright   bool           `json:"copyright"`
	Status      int            `json:"status"`
	PartitionId uint           `json:"partitionId"`
	Tags        string         `json:"tags"`
	Resources   []ResourceResp `json:"resources" gorm:"-"`
}

type UploadVideoResp struct {
	ID        uint      `json:"vid"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Desc      string    `json:"desc"`
	Copyright bool      `json:"copyright"`
	CreatedAt time.Time `json:"createdAt"`
	Clicks    int64     `json:"clicks" gorm:"-"`
}

type AllVideoResp struct {
	ID    uint   `json:"vid"`
	Title string `json:"title"`
}

type VideoInfoManageResp struct {
	ID        uint         `json:"vid"`
	Uid       uint         `json:"uid"`
	Title     string       `json:"title"`
	Cover     string       `json:"cover"`
	Desc      string       `json:"desc"`
	CreatedAt time.Time    `json:"createdAt"`
	Copyright bool         `json:"copyright"`
	Tags      string       `json:"tags"`
	Clicks    int64        `json:"clicks" gorm:"-"`
	Author    UserInfoResp `json:"author" gorm:"-"`
}

type ReviewListResp struct {
	ID        uint         `json:"vid"`
	Uid       uint         `json:"uid"`
	Title     string       `json:"title"`
	Cover     string       `json:"cover"`
	Desc      string       `json:"desc"`
	CreatedAt time.Time    `json:"createdAt"`
	Copyright bool         `json:"copyright"`
	Tags      string       `json:"tags"`
	Author    UserInfoResp `json:"author" gorm:"-"`
}
