package model

import "gorm.io/gorm"

// AudioTrack 音频轨道（支持多音轨：多语言、多音轨切换）
// 一个 Resource 可以有多个 AudioTrack，分别对应不同语言或音轨。
// 与 VideoIndexFile 解耦：音频与清晰度无关，所有清晰度共享同一套音轨文件。
type AudioTrack struct {
	gorm.Model
	ResourceID uint   `gorm:"index;comment:资源ID"`
	DirName    string `gorm:"type:varchar(50);comment:存储目录名"`

	// 音轨元信息
	Language    string `gorm:"type:varchar(10);comment:语言代码(ISO 639-2, 如 jpn/eng/und)"`
	Title       string `gorm:"type:varchar(50);comment:显示名称(日语/英语/未知)"`
	TrackIndex  int    `gorm:"comment:ffprobe stream index(用于编码映射)"`
	IsDefault   bool   `gorm:"default:false;comment:是否默认音轨"`
	Channels    int    `gorm:"default:2;comment:声道数"`

	// 文件信息
	AudioFile       string `gorm:"type:varchar(100);comment:音频文件名(audio_jpn.m4s)"`
	Codec           string `gorm:"type:varchar(50);comment:音频编解码器"`
	Bandwidth       int    `gorm:"comment:音频码率(bps)"`
	SampleRate      int    `gorm:"comment:采样率(Hz)"`
	InitRange       string `gorm:"type:varchar(30);comment:初始化字节范围"`
	IndexRange      string `gorm:"type:varchar(30);comment:索引字节范围"`
}

func (AudioTrack) TableName() string {
	return "audio_track"
}
