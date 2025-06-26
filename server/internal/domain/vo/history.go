package vo

import "time"

const (
	HISTORY_VIDEO_FIELD = "`video`.`id`,`video`.`uid`,`title`,`cover`,`desc`,`history`.`updated_at`,`time`"
)

type HistoryVideoResp struct {
	ID        uint      `json:"vid"`
	Uid       uint      `json:"uid"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Desc      string    `json:"desc"`
	Time      float64   `json:"time"`
	UpdatedAt time.Time `json:"updatedAt"`
	Part      uint      `json:"part"`      // 分P序号
	//可能不用了
	//PartTitle string    `json:"partTitle"` // 分P标题
}
