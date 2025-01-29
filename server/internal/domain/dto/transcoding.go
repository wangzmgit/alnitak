package dto

type TranscodingInfo struct {
	Width      int     // 视频宽度
	Height     int     // 视频高度
	Duration   float64 // 视频时长
	DirName    string  // 目录名称
	OutputDir  string  // 输出位置
	InputFile  string  // 输入文件
	ResourceID uint    // 资源ID
	VideoID    uint    // 视频ID
	FPS        string  // 视频帧率
}
