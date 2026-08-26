package ffmpeg

// StreamProbe 对应 ffprobe 输出的一个 stream
type StreamProbe struct {
	CodecType string     `json:"codec_type"`
	CodecName string     `json:"codec_name"`
	Width     int        `json:"width,omitempty"`
	Height    int        `json:"height,omitempty"`
	PixFmt    string     `json:"pix_fmt,omitempty"`
	Duration  string     `json:"duration"`
	RFrameRate string    `json:"r_frame_rate,omitempty"`
	AvgFrameRate string  `json:"avg_frame_rate,omitempty"`
	SampleRate string    `json:"sample_rate,omitempty"`
	Channels  int        `json:"channels,omitempty"`
	BitRate   string     `json:"bit_rate,omitempty"`
	Tags      StreamTags `json:"tags,omitempty"`
}

// StreamTags ffprobe stream.tags
type StreamTags struct {
	Language string `json:"language,omitempty"`
	Title    string `json:"title,omitempty"`
}

// FormatProbe 对应 ffprobe 输出的 format
type FormatProbe struct {
	BitRate  string `json:"bit_rate"`
	Duration string `json:"duration"`
}

// ProbeResult ffprobe 完整输出
type ProbeResult struct {
	Streams []StreamProbe `json:"streams"`
	Format  FormatProbe   `json:"format"`
}
