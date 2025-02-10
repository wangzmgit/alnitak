package config

type Config struct {
	Cors        Cors        `mapstructure:"cors" json:"cors" yaml:"cors"`
	File        File        `mapstructure:"file" json:"file" yaml:"file"`
	Log         Log         `mapstructure:"log" json:"log" yaml:"log"`
	Mail        Mail        `mapstructure:"mail" json:"mail" yaml:"mail"`
	Mysql       Mysql       `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis       Redis       `mapstructure:"redis" json:"redis" yaml:"redis"`
	Security    Security    `mapstructure:"security" json:"security" yaml:"security"`
	Storage     Storage     `mapstructure:"storage" json:"storage" yaml:"storage"`
	Transcoding Transcoding `mapstructure:"transcoding" json:"transcoding" yaml:"transcoding"`
	User        User        `mapstructure:"user" json:"user" yaml:"user"`
}
