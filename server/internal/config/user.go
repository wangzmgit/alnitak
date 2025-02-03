package config

type Cors struct {
	AllowOrigin string `mapstructure:"allow_origin" json:"allow_origin" yaml:"allow_origin"`
}
