package config

type Transcoding struct {
	UseGpu          bool `mapstructure:"use_gpu" json:"use_gpu" yaml:"use_gpu"`
	Generate1080p60 bool `mapstructure:"generate_1080p60" json:"generate_1080p60" yaml:"generate_1080p60"`
}
