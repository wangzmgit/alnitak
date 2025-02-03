package config

type File struct {
	MaxImgSize   int `mapstructure:"max_img_size" json:"max_img_size" yaml:"max_img_size"`
	MaxVideoSize int `mapstructure:"max_video_size" json:"max_video_size" yaml:"max_video_size"`
}
