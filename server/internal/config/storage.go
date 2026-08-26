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
	UseSSL        bool   `mapstructure:"use_ssl" json:"use_ssl" yaml:"use_ssl"` // OSS是否使用HTTPS
	UploadMp4File bool   `mapstructure:"upload_mp4_file" json:"upload_mp4_file" yaml:"upload_mp4_file"`
	UploadTimeout int    `mapstructure:"upload_timeout" json:"upload_timeout" yaml:"upload_timeout"` // 上传超时（秒），0=默认30s
	Backup        *StorageBackup `mapstructure:"backup" json:"backup" yaml:"backup"` // 备用OSS（多源容灾）
}

// StorageBackup 备用存储配置，用于多 OSS 源容灾。
// Private/UseSSL/UploadTimeout 继承主配置。
type StorageBackup struct {
	OssType   string `mapstructure:"oss_type" json:"oss_type" yaml:"oss_type"`
	Endpoint  string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	Bucket    string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	KeyId     string `mapstructure:"key_id" json:"key_id" yaml:"key_id"`
	KeySecret string `mapstructure:"key_secret" json:"key_secret" yaml:"key_secret"`
	AppId     string `mapstructure:"app_id" json:"app_id" yaml:"app_id"`
	Region    string `mapstructure:"region" json:"region" yaml:"region"`
	Domain    string `mapstructure:"domain" json:"domain" yaml:"domain"`
}

// ToStorageConfig 将备用配置转为主 Storage 格式，继承主配置的 Private/UseSSL/UploadTimeout。
func (b *StorageBackup) ToStorageConfig(primary Storage) Storage {
	return Storage{
		OssType:       b.OssType,
		Bucket:        b.Bucket,
		Endpoint:      b.Endpoint,
		KeyId:         b.KeyId,
		KeySecret:     b.KeySecret,
		AppId:         b.AppId,
		Region:        b.Region,
		Domain:        b.Domain,
		Private:       primary.Private,
		UseSSL:        primary.UseSSL,
		UploadTimeout: primary.UploadTimeout,
	}
}
