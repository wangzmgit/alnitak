package config

type Storage struct {
	Bucket        string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Endpoint      string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	KeyId         string `mapstructure:"key_id" json:"key_id" yaml:"key_id"`
	AppId         string `mapstructure:"app_id" json:"app_id" yaml:"app_id"`
	KeySecret     string `mapstructure:"key_secret" json:"key_secret" yaml:"key_secret"`
	OssType       string `mapstructure:"oss_type" json:"oss_type" yaml:"oss_type"`
	Region        string `mapstructure:"region" json:"region" yaml:"region"`
	Domain        string `mapstructure:"domain" json:"domain" yaml:"domain"`
	Private       bool   `mapstructure:"private" json:"private" yaml:"private"`
	UploadMp4File bool   `mapstructure:"upload_mp4_file" json:"upload_mp4_file" yaml:"upload_mp4_file"`
}
