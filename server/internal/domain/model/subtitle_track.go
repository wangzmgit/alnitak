package model

import "gorm.io/gorm"

// 字幕来源（预留 PGC / 自动识别）
const (
	SubtitleOriginUser = 1
)

const (
	SubtitleStatusActive = 1
)

// 上传时的原始格式（存储统一为 WebVTT 文件）
const (
	SubtitleFormatSRT = 1
	SubtitleFormatVTT = 2
)

// SubtitleTrack 分 P 字幕轨（元数据在库，文件在 subtitle/ 前缀下）
type SubtitleTrack struct {
	gorm.Model
	ShortID          string `gorm:"type:varchar(16);uniqueIndex;comment:短ID"`
	ResourceShortID  string `gorm:"type:varchar(16);not null;index;uniqueIndex:uk_subtitle_short_id_lang,priority:1;comment:资源短ID"`
	Vid              uint   `gorm:"not null;index;comment:冗余稿件 ID"`
	Lang             string `gorm:"type:varchar(20);not null;uniqueIndex:uk_subtitle_short_id_lang,priority:2;comment:BCP-47 语言"`
	Label      string `gorm:"type:varchar(64);comment:展示名"`
	SourceFmt  int    `gorm:"column:source_fmt;not null;comment:原始格式 1=SRT 2=VTT"`
	ObjectKey  string `gorm:"type:varchar(255);not null;comment:对象键 subtitle/{id}.vtt"`
	IsDefault  bool   `gorm:"column:is_default;not null;default:0;comment:是否默认轨"`
	Origin     int    `gorm:"not null;default:1;comment:来源"`
	Status     int    `gorm:"not null;default:1;comment:状态 1=可用"`
	CreatedBy  uint   `gorm:"column:created_by;not null;index;comment:上传者"`
}

func (SubtitleTrack) TableName() string {
	return "subtitle_track"
}
