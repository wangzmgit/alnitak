package vo

// SubtitleTrackItem 字幕列表项（给播放端 / 创作端）
type SubtitleTrackItem struct {
	ID        uint   `json:"id"`
	ShortID   string `json:"shortId"`
	Lang      string `json:"lang"`
	Label     string `json:"label"`
	URL       string `json:"url"`
	IsDefault bool   `json:"isDefault"`
}
