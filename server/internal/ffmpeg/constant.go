package ffmpeg

// 音视频编解码常量
const (
	AudioMaxBitrateBps = 320000
	DefaultAudioSample = 48000
	DefaultAudioChan   = 2
	FragDurationUs     = "2000000" // 2秒切片
	TimeBase30fps      = "30000/1000"
	TimeBase60fps      = "60000/1001"
	DefaultVideoCodec  = "avc1.640028"
	DefaultHEVCCodec   = "hvc1.1.6.L150.B0"
	DefaultAV1Codec    = "av01.0.08M.08"
	DefaultAudioCodec  = "mp4a.40.2"
)

// QualityPreset 分辨率 - 码率对照
type QualityPreset struct {
	LongSide   int
	ShortSide  int
	BitrateH   string
	BitrateV   string
	Bitrate60H string
	Bitrate60V string
}

var QualityPresets = []QualityPreset{
	{1920, 1080, "8000k", "5000k", "12000k", "8000k"},
	{1280, 720, "5000k", "3000k", "7500k", "5000k"},
	{854, 480, "2500k", "1500k", "", ""},
	{640, 360, "1000k", "700k", "", ""},
}
