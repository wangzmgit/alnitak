package dto

// AudioStreamProbe ffprobe 探测到的音频流信息
type AudioStreamProbe struct {
	StreamIndex int    // ffprobe stream index
	Language    string // 语言代码 (ISO 639-2, 如 jpn/eng/und)
	SampleRate  int    // 采样率
	Channels    int    // 声道数
	BitRate     int    // 码率 (bps)
}

type TranscodingInfo struct {
	Width        int     // 视频宽度
	Height       int     // 视频高度
	Duration     float64 // 视频时长
	DirName      string  // 目录名称
	OutputDir    string  // 输出位置
	InputFile    string  // 输入文件
	ResourceID   uint    // 资源ID
	VideoID      uint    // 视频ID
	CodecName    string  // 视频编码名称
	FPS          string  // 视频帧率
	FPS30        string  // 30帧实际帧率
	FPS60        string  // 60帧实际帧率
	VideoBitRate int     // 视频码率 (bps)，优先取视频流码率
	Suffix       string  // 原始文件后缀（如 .mkv）

	// 音频源文件参数（旧单音轨模式，向后兼容）
	AudioBitRate    int // 音频码率 (bps)，如 320000、192000
	AudioSampleRate int // 音频采样率 (Hz)，如 48000、44100
	AudioChannels   int // 音频声道数，如 2（立体声）、6（5.1）

	// 多音轨支持：ffprobe 探测到的所有音频流
	AudioStreams []AudioStreamProbe

	// 重新转码时记录原始视频状态，转码完成后恢复（-1表示普通上传转码，无需特殊处理）
	OriginalVideoStatus int

	// 转码编码参数（由主服务下发，Worker 取此为准而非本机配置，保证前后端一致）
	UseGpu            bool `json:"useGpu"`
	UseH265           bool `json:"useH265"`
	UseAv1            bool `json:"useAv1"`
	Generate1080p60   bool `json:"generate1080p60"`
}
