package global

type VideoInfo struct {
	Stream []Streams `json:"streams"`
	Format Format    `json:"format"`
}

// StreamTags ffprobe stream.tags
type StreamTags struct {
	Language string `json:"language,omitempty"`
	Title    string `json:"title,omitempty"`
}

type Streams struct {
	CodecType    string     `json:"codec_type"`               // 流类型: "video" 或 "audio"
	CodecName    string     `json:"codec_name"`
	Width        int        `json:"width,omitempty"`
	Height       int        `json:"height,omitempty"`
	PixFmt       string     `json:"pix_fmt,omitempty"`
	Duration     string     `json:"duration"`
	RFrameRate   string     `json:"r_frame_rate,omitempty"`   // 原始帧率 (如 "30/1")
	AvgFrameRate string     `json:"avg_frame_rate,omitempty"` // 平均帧率 (如 "30/1")
	SampleRate   string     `json:"sample_rate,omitempty"`    // 音频采样率 (如 "48000")
	Channels     int        `json:"channels,omitempty"`       // 音频声道数
	BitRate      string     `json:"bit_rate,omitempty"`       // 流码率 (如 "320000")
	BitsPerSample string    `json:"bits_per_raw_sample,omitempty"` // 音频位深
	Tags         StreamTags `json:"tags,omitempty"`           // 元数据标签(含语言)
}

type Format struct {
	BitRate  string `json:"bit_rate"`
	Duration string `json:"duration"`
}
