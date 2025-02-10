package config

type Mail struct {
	Addresser string `mapstructure:"addresser" json:"addresser" yaml:"addresser"`
	Host      string `mapstructure:"host" json:"host" yaml:"host"`
	Port      int    `mapstructure:"port" json:"port" yaml:"port"`
	User      string `mapstructure:"user" json:"user" yaml:"user"`
	Pass      string `mapstructure:"pass" json:"pass" yaml:"pass"`
	Debug     bool   `mapstructure:"debug" json:"debug" yaml:"debug"`
}
