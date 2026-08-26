package vo

import "time"

const (
	HISTORY_VIDEO_FIELD    = "`video`.`id`,`video`.`uid`,`video`.`title`,`video`.`cover`,`video`.`desc`,`history`.`updated_at`,`history`.`time`,`history`.`duration`,`history`.`part`,`video`.`pgc_attached`,`video`.`short_id`,`history`.`resource_short_id` as `rid`"
	HISTORY_SUBQUERY_FIELD = "video_short_id, MAX(updated_at) as latest_updated_at"
)

type HistoryVideoResp struct {
	ID        uint      `json:"vid"`
	Uid       uint      `json:"uid"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Desc      string    `json:"desc"`
	Time      float64   `json:"time"`     // 播放进度
	Duration  float64   `json:"duration"` // 分P总时长
	UpdatedAt time.Time `json:"updatedAt"`
	Part      uint      `json:"part"` // 分P序号

	ShortID string `json:"shortId,omitempty"` // 用于 UGC 跳转

	// 资源短ID，用于精准跳转（不受分P排序影响）
	Rid string `json:"rid,omitempty"`

	// PGC：存储仍是 vid+part；列表展示与续播跳转按剧集维度区分
	PGCAttached   bool   `json:"pgcAttached"`
	PGCTitle      string `json:"pgcTitle,omitempty"`
	EpisodeTitle  string `json:"episodeTitle,omitempty"`
	EpisodeNumber int    `json:"episodeNumber,omitempty"`
	EpID          uint   `json:"epId,omitempty"`
}
