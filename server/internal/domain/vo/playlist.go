package vo

import "time"

const (
	PLAYLIST_FIELD      = "`id`,`uid`,`title`,`cover`,`desc`,`is_open`,`status`,`video_count`,`views`,`favorites`,`created_at`,`updated_at`"
	PLAYLIST_LIST_FIELD = "`id`,`title`,`cover`,`desc`,`is_open`,`status`,`video_count`,`views`,`favorites`,`created_at`"
)

// 合集详情
type PlaylistResp struct {
	ID         uint         `json:"id"`
	Uid        uint         `json:"uid"`
	Title      string       `json:"title"`
	Cover      string       `json:"cover"`
	Desc       string       `json:"desc"`
	IsOpen     bool         `json:"isOpen"`
	Status     int          `json:"status"`
	VideoCount int          `json:"videoCount"`
	Views      int64        `json:"views"`
	Favorites  int64        `json:"favorites"`
	CreatedAt  time.Time    `json:"createdAt"`
	UpdatedAt  time.Time    `json:"updatedAt"`
	Author     UserInfoResp `json:"author" gorm:"-"`
}

// 合集列表项
type PlaylistListResp struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Cover      string    `json:"cover"`
	Desc       string    `json:"desc"`
	IsOpen     bool      `json:"isOpen"`
	Status     int       `json:"status"`
	VideoCount int       `json:"videoCount"`
	Views      int64     `json:"views"`
	Favorites  int64     `json:"favorites"`
	CreatedAt  time.Time `json:"createdAt"`
}

// 合集中的视频项
type PlaylistVideoResp struct {
	Vid       uint      `json:"-"`
	ShortID   string    `json:"shortId"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Duration  int       `json:"duration"`
	Clicks    int64     `json:"clicks"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"createdAt"`
}

// 视频所属的合集项（用于视频页面显示）
type VideoPlaylistResp struct {
	ID        uint            `json:"id"`
	Uid       uint            `json:"uid"`
	Title     string          `json:"title"`
	Cover     string          `json:"cover"`
	Desc      string          `json:"desc"`
	IsOpen    bool            `json:"isOpen"`
	CreatedAt time.Time       `json:"createdAt"`
	Author    UserInfoResp    `json:"author" gorm:"-"`
}
