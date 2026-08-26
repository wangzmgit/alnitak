package vo

import (
	"time"
	"interastral-peace.com/alnitak/internal/domain/types"
)

const (
	// 注意：包含 short_id；tags 由多对多 tag / video_tag 组装，不从此处直接 Scan
	VIDEO_FIELD        = "`id`,`short_id`,`uid`,`title`,`cover`,`desc`,`created_at`,`copyright`,`clicks`,`duration`,`partition_id`"
	VIDEO_STATUS_FIELD = "`id`,`short_id`,`title`,`cover`,`desc`,`copyright`,`status`,`partition_id`,`clicks`,`duration`"
	UPLOAD_VIDEO_FIELD = "`id`,`short_id`,`title`,`cover`,`desc`,`copyright`,`status`,`created_at`,`clicks`,`duration`"
)

type VideoResp struct {
	ID           uint           `json:"vid,omitempty"`
	ShortID      string         `json:"shortId"`
	Uid          uint           `json:"uid"`
	Title        string         `json:"title"`
	Cover        string         `json:"cover"`
	Desc         string         `json:"desc"`
	CreatedAt    time.Time      `json:"createdAt"`
	Copyright    types.CopyrightType `json:"copyright"`
	Tags         []string            `json:"tags" gorm:"-"` // 非表字段，避免 Scan 时 GORM 将 []string 误判为关联
	Duration     int                 `json:"duration"` // 秒，整数
	Clicks       int64               `json:"clicks"`
	PartitionId  uint                `json:"partitionId"`
	DanmakuCount int64               `json:"danmakuCount" gorm:"-"` // 弹幕数量
	Fans         int64               `json:"fans" gorm:"-"`         // 作者粉丝数（与 author 解耦，避免在 author 内塞统计）
	Author       AuthorPublicResp    `json:"author" gorm:"-"`
	Resources    []ResourceResp      `json:"resources" gorm:"-"`
}

type VideoStatusResp struct {
	ID          uint           `json:"vid"`
	ShortID     string         `json:"shortId"`
	Title       string         `json:"title"`
	Cover       string         `json:"cover"`
	Desc        string         `json:"desc"`
	Copyright   types.CopyrightType `json:"copyright"`
	Status      int                 `json:"status"`
	PartitionId uint                `json:"partitionId"`
	Tags        []string            `json:"tags" gorm:"-"`
	Clicks      int64               `json:"clicks"`
	Duration    int                 `json:"duration"`
	Resources   []ResourceResp      `json:"resources" gorm:"-"`
}

type UploadVideoResp struct {
	ID                  uint                      `json:"vid"`
	ShortID             string                    `json:"shortId"`
	Title               string                    `json:"title"`
	Cover               string                    `json:"cover"`
	Desc                string                    `json:"desc"`
	Status              int                       `json:"status"`
	Copyright           types.CopyrightType        `json:"copyright"`
	CreatedAt           time.Time                  `json:"createdAt"`
	Clicks              int64                      `json:"clicks"`
	Duration            int                        `json:"duration"`
	TranscodingProgress float64                    `json:"transcodingProgress,omitempty" gorm:"-"`
	TranscodingDetails  []TranscodingProgressItem `json:"transcodingDetails,omitempty" gorm:"-"`
	UploadProgress      *UploadProgressInfo       `json:"uploadProgress,omitempty" gorm:"-"`
}

type AllVideoResp struct {
	ID      uint   `json:"-"`
	ShortID string `json:"shortId"`
	Title   string `json:"title"`
	Cover   string `json:"cover"`
}

type VideoInfoManageResp struct {
	ID                  uint                      `json:"vid"`
	ShortID             string                    `json:"shortId"`
	Uid                 uint                      `json:"uid"`
	Title               string                    `json:"title"`
	Cover               string                    `json:"cover"`
	Desc                string                    `json:"desc"`
	CreatedAt           time.Time                 `json:"createdAt"`
	Copyright           types.CopyrightType        `json:"copyright"`
	Tags                []string                   `json:"tags" gorm:"-"`
	Clicks              int64                      `json:"clicks"`
	Duration            int                        `json:"duration"`
	PartitionId         uint                       `json:"partitionId"`
	Author              UserInfoResp               `json:"author" gorm:"-"`
	TranscodingProgress float64                    `json:"transcodingProgress,omitempty" gorm:"-"`
	TranscodingDetails  []TranscodingProgressItem `json:"transcodingDetails,omitempty" gorm:"-"`
	UploadProgress      *UploadProgressInfo       `json:"uploadProgress,omitempty" gorm:"-"`
}

type TranscodingProgressItem struct {
	ResourceID    uint               `json:"resourceId"`
	ResourceTitle string             `json:"resourceTitle"`
	Quality       string             `json:"quality"`
	Progress      float64            `json:"progress"`
	Status        string             `json:"status"` // processing/success/fail
	Upload        *UploadProgressInfo `json:"upload,omitempty"`
}

type UploadProgressInfo struct {
	OssType  string  `json:"ossType"`  // aliyun/minio/cloudflare/local
	Progress float64 `json:"progress"` // 0-100
	Status   string  `json:"status"`   // uploading/success/fail/local
}

type ReviewListResp struct {
	ID          uint         `json:"vid"`
	ShortID     string       `json:"shortId"`
	Uid         uint         `json:"uid"`
	Title       string       `json:"title"`
	Cover       string       `json:"cover"`
	Desc        string       `json:"desc"`
	CreatedAt   time.Time    `json:"createdAt"`
	Copyright   types.CopyrightType `json:"copyright"`
	Tags        []string            `json:"tags" gorm:"-"`
	Duration    int                 `json:"duration"`
	PartitionId uint                `json:"partitionId"`
	Author      UserInfoResp        `json:"author" gorm:"-"`
}
