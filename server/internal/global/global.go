package global

import (
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
	"interastral-peace.com/alnitak/internal/config"
	"interastral-peace.com/alnitak/pkg/casbin"
	"interastral-peace.com/alnitak/pkg/oss"
	"interastral-peace.com/alnitak/pkg/redis"
)

var (
	Config            *config.Config
	Mysql             *gorm.DB
	Redis             *redis.Redis
	Casbin            *casbin.Casbin
	Storage           oss.Storage
	StorageBackup     oss.Storage // 备用OSS实例（多源容灾），可为nil
	SnowflakeNode     *snowflake.Node
	VideoPartitionMap map[uint]uint
)

// GetOssUrl 获取主 OSS 访问 URL
// 公开bucket: 直接拼接公开URL（短、不过期）
// 私有bucket: 通过SDK生成预签名URL（带认证参数）
func GetOssUrl(objectKey string) string {
	return buildOssURL(objectKey, Config.Storage)
}

// GetBackupOssUrl 生成备用 OSS 的访问 URL（用于播放层容灾）。
// 公开 bucket: 直接拼接 URL
// 私有 bucket: 通过 StorageBackup SDK 生成预签名 URL
// 无备用 OSS 时返回空。
func GetBackupOssUrl(objectKey string) string {
	if StorageBackup == nil {
		return ""
	}
	backupCfg := Config.Storage.Backup.ToStorageConfig(Config.Storage)
	if backupCfg.Private {
		// 私有 bucket：通过备份 OSS SDK 实例生成预签名 URL
		return StorageBackup.GetObjectUrl(objectKey)
	}
	// 公开 bucket：直接拼接
	return buildOssURL(objectKey, backupCfg)
}

func buildOssURL(objectKey string, cfg config.Storage) string {
	if cfg.Private {
		if cfg.Bucket == Config.Storage.Bucket && cfg.Endpoint == Config.Storage.Endpoint {
			return Storage.GetObjectUrl(objectKey)
		}
		return ""
	}
	scheme := "http://"
	if cfg.UseSSL {
		scheme = "https://"
	}
	domain := cfg.Domain
	if domain == "" {
		domain = cfg.Endpoint
	}
	return scheme + domain + "/" + cfg.Bucket + "/" + objectKey
}
