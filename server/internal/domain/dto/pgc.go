package dto

type CreatePGCReq struct {
	PGCType   int     `json:"pgc_type" binding:"required"`
	Title     string  `json:"title" binding:"required"`
	Cover     string  `json:"cover" binding:"required"`
	Desc      string  `json:"desc"`
	Year      int     `json:"year"`
	Area      string  `json:"area"`
	Rating    float64 `json:"rating"`
	IsOngoing bool    `json:"is_ongoing"`
	// Episodes 可选：先创建 season（内容主体），剧集在后续接口逐集添加，对齐常见 PGC 流程
	Episodes []EpisodeReq `json:"episodes"`
	// 无剧集创建时可选填计划总集数（仅写入 total_episodes，current_episodes 为 0）
	PlannedTotalEpisodes *int `json:"total_episodes"`
}

type EpisodeReq struct {
	EpisodeNumber int    `json:"episode_number" binding:"required"`
	Title         string `json:"title"`
	// VID 可选：0 表示仅创建剧集占位，后续通过 bind 接口绑定视频
	VID         uint    `json:"vid"`
	Duration    float64 `json:"duration"`
	PublishTime string  `json:"publish_time"`
}

// BindPGCEpisodeVideoReq 为占位剧集绑定已存在的视频
type BindPGCEpisodeVideoReq struct {
	VID         uint     `json:"vid" binding:"required"`
	Duration    *float64 `json:"duration"`
	PublishTime string   `json:"publish_time"`
}

type UpdatePGCReq struct {
	// PGCID 由 Snowflake 生成（uint64），前端 JS Number 无法安全表示，所以允许用 string 传递
	PGCID         uint64   `json:"pgc_id,string" binding:"required"`
	Title         string   `json:"title"`
	Cover         string   `json:"cover"`
	Desc          string   `json:"desc"`
	Year          int      `json:"year"`
	Area          string   `json:"area"`
	Rating        *float64 `json:"rating"`
	IsOngoing     *bool    `json:"is_ongoing"`
	TotalEpisodes *int     `json:"total_episodes"`
}

type UpdatePGCStatusReq struct {
	PGCID  uint64 `json:"pgc_id,string" binding:"required"`
	Status int    `json:"status" binding:"required"`
}

type UpdatePGCEpisodeStatusReq struct {
	Status int `json:"status" binding:"required"`
}

type UpdatePGCEpisodeReq struct {
	Title string `json:"title"`
}

type PGCListReq struct {
	Page      int    `form:"page" binding:"required,min=1"`
	PageSize  int    `form:"page_size" binding:"required,min=1,max=100"`
	PGCType   int    `form:"pgc_type"`
	Status    *int   `form:"status"`
	Keyword   string `form:"keyword"`
	Year      int    `form:"year"`
	Area      string `form:"area"`
	IsOngoing *bool  `form:"is_ongoing"`
}

type EpisodeListReq struct {
	// PGCID 来自 path 参数 `/:pgc_id/episodes`，不从 query 绑定，也不参与校验（避免 bind 阶段 required 失败）
	PGCID    uint `form:"-"`
	Page     int  `form:"page" binding:"required,min=1"`
	PageSize int  `form:"page_size" binding:"required,min=1,max=100"`
}

// PGCManageListReq 后台 PGC 内容管理列表（POST JSON）
type PGCManageListReq struct {
	Page     int    `json:"page" binding:"required,min=1"`
	PageSize int    `json:"pageSize" binding:"required,min=1,max=100"`
	PGCType  int    `json:"pgc_type"`
	Status   *int   `json:"status"`
	Keyword  string `json:"keyword"`
}

// PGCReviewListReq 后台 PGC 待审列表
type PGCReviewListReq struct {
	Page     int `json:"page" binding:"required,min=1"`
	PageSize int `json:"pageSize" binding:"required,min=1,max=100"`
}

// PGCReviewActionReq 后台通过 / 驳回
type PGCReviewActionReq struct {
	PGCID uint64 `json:"pgc_id,string" binding:"required"`
}
