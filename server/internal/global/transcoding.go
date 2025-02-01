package global

type VideoInfo struct {
	Stream []Streams `json:"streams"`
	Format Format    `json:"format"`
}

type Streams struct {
	CodecName    string `json:"codec_name"`
	Width        int    `json:"width,omitempty"`
	Height       int    `json:"height,omitempty"`
	PixFmt       string `json:"pix_fmt,omitempty"`
	Duration     string `json:"duration"`
	RFrameRate   string `json:"r_frame_rate,omitempty"`  // 原始帧率 (如 "30/1")
	AvgFrameRate string `json:"avg_frame_rate,omitempty"` // 平均帧率 (如 "30/1")
}

type Format struct {
	BitRate string `json:"bit_rate"`
}
