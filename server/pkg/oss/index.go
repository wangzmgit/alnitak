package oss

import (
	"errors"
	"io"
	"time"

	"interastral-peace.com/alnitak/internal/config"
	"interastral-peace.com/alnitak/utils"
)

const (
	ALIYUN  = "aliyun"
	MINIO   = "minio"
	TENCENT = "tencent"
	//QINIU    = "qiniu"
	R2TORAGE = "cloudflare"
)

// CompletePart 表示一个已上传的分片，用于 CompleteMultipartUpload。
type CompletePart struct {
	PartNumber int    `json:"partNumber"`
	ETag       string `json:"etag"`
}

type Storage interface {
	GetObjectToFile(objectKey, downloadedFileName string) error
	DeleteObject(objectKey string) error
	PutObject(objectKey string, reader io.Reader) error
	PutObjectFromFile(objectKey, filePath string) error
	IsExists(objectKey string) (bool, error)
	GetObjectUrl(objectKey string) string
	GetObjectReader(objectKey string) (io.ReadCloser, error)

	// ── 直传 OSS 预签名方法 ──

	// PresignPutObject 生成单文件上传预签名 URL（用于图片直传）。
	// expiry: 签名有效期。
	PresignPutObject(objectKey string, expiry time.Duration) (string, error)

	// InitiateMultipartUpload 发起分片上传，返回 uploadID。
	// 用于视频直传场景：客户端拿到 uploadID 后逐片预签名并直传。
	InitiateMultipartUpload(objectKey string) (string, error)

	// PresignUploadPart 为指定分片生成预签名 PUT URL。
	// partNumber 从 1 开始。
	PresignUploadPart(uploadID, objectKey string, partNumber int, expiry time.Duration) (string, error)

	// CompleteMultipartUpload 通知 OSS 合并所有已上传分片。
	CompleteMultipartUpload(uploadID, objectKey string, parts []CompletePart) error
}

func InitStorage(c config.Storage) Storage {
	clientConfig := buildOssConfig(c)

	s, err := initOss(c.OssType, clientConfig)
	if err != nil {
		utils.ErrorLog("oss初始化失败", "oss", err.Error())
		panic(err)
	}

	return s
}

// InitBackupStorage 初始化备用 OSS 实例（多源容灾），忽略 OssType=local 或配置为空。
// 继承主配置的 Private/UseSSL/UploadTimeout。
func InitBackupStorage(primary config.Storage) Storage {
	if primary.Backup == nil || primary.Backup.OssType == "" || primary.Backup.OssType == "local" {
		return nil
	}

	merged := primary.Backup.ToStorageConfig(primary)
	return InitStorage(merged)
}

func buildOssConfig(c config.Storage) Config {
	return Config{
		KeyID:     c.KeyId,
		KeySecret: c.KeySecret,
		Bucket:    c.Bucket,
		Endpoint:  c.Endpoint,
		AppID:     c.AppId,
		Region:    c.Region,
		Domain:    c.Domain,
		Private:   c.Private,
		UseSSL:    c.UseSSL,
		Timeout:   c.UploadTimeout,
	}
}

func initOss(ossName string, config Config) (Storage, error) {
	switch ossName {
	case ALIYUN:
		return newAliyun(config)
	case MINIO:
		return newMinio(config)
	case TENCENT:
		return newTencentCOS(config)
	//case QINIU:
	//	return newQiniuOSS(config)
	case R2TORAGE:
		return newR2Storage(config)
	default:
		return nil, errors.New("driver not exists")
	}
}
