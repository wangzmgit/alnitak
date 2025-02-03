package config

type Log struct {
	FileName   string `mapstructure:"filename" json:"filename" yaml:"filename"`
	MaxAge     int    `mapstructure:"max-age" json:"max-age" yaml:"max-age"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	MaxSize    int    `mapstructure:"max-size" json:"max-size" yaml:"max-size"`
	Mode       string `mapstructure:"mode" json:"mode" yaml:"mode"`
}
