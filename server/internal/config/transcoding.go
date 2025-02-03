package config

type Transcoding struct {
	Gpu        bool `mapstructure:"gpu" json:"gpu" yaml:"gpu"`
	Res1080p60 bool `mapstructure:"1080p60" json:"1080p60" yaml:"1080p60"`
}
