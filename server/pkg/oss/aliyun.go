package oss

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"interastral-peace.com/alnitak/utils"
)

type Aliyun struct {
	config *Config
	bucket *oss.Bucket
}

func newAliyun(config Config) (*Aliyun, error) {
	aliyun := Aliyun{}
	err := aliyun.init(config)
	if err != nil {
		return nil, err
	}
	return &aliyun, nil
}

func (a *Aliyun) init(config Config) error {
	if a.config == nil {
		a.config = &config
	}

	if config.Endpoint == "" {
		return errors.New("configuration not correct")
	}

	if a.bucket == nil {
		client, err := a.initOssClinet(config)
		if err != nil {
			return err
		}

		bucket, err := client.Bucket(config.Bucket)
		if err != nil {
			return err
		}
		a.bucket = bucket
	}

	return nil
}

func (a *Aliyun) initOssClinet(config Config) (*oss.Client, error) {
	opts := make([]oss.ClientOption, 0)
	if config.Domain != "" {
		opts = append(opts, oss.UseCname(true))
	}
	if config.Timeout > 0 {
		opts = append(opts, oss.Timeout(int64(config.Timeout), int64(config.Timeout)))
	}

	endpoint := config.Endpoint
	if config.Domain != "" {
		endpoint = config.Domain
	}

	return oss.New(endpoint, config.KeyID, config.KeySecret, opts...)
}

// 获取文件
func (a *Aliyun) GetObjectToFile(objectKey, filePath string) error {
	return a.bucket.GetObjectToFile(objectKey, filePath)
}

// 删除文件
func (a *Aliyun) DeleteObject(objectKey string) error {
	return a.bucket.DeleteObject(objectKey)
}

func (a *Aliyun) PutObject(objectKey string, reader io.Reader) error {
	return a.bucket.PutObject(objectKey, reader)
}

func (a *Aliyun) PutObjectFromFile(objectKey, filePath string) error {
	return a.bucket.PutObjectFromFile(objectKey, filePath)
}

func (a *Aliyun) IsExists(objectKey string) (bool, error) {
	return a.bucket.IsObjectExist(objectKey)
}

// 获取访问URL
func (a *Aliyun) GetObjectUrl(objectKey string) string {
	url, err := a.bucket.SignURL(objectKey, oss.HTTPGet, 18000) // 5小时
	if err != nil {
		utils.ErrorLog("OSS生成文件URL失败", "transcoding", err.Error())
		return ""
	}

	return url
}

// 获取对象读取器（用于代理流式传输）
func (a *Aliyun) GetObjectReader(objectKey string) (io.ReadCloser, error) {
	return a.bucket.GetObject(objectKey)
}

// ── 直传 OSS 预签名方法 ──

// PresignPutObject 生成单文件上传预签名 URL（用于图片直传）。
func (a *Aliyun) PresignPutObject(objectKey string, expiry time.Duration) (string, error) {
	url, err := a.bucket.SignURL(objectKey, oss.HTTPPut, int64(expiry.Seconds()))
	if err != nil {
		return "", fmt.Errorf("presign put object: %w", err)
	}
	return url, nil
}

// InitiateMultipartUpload 发起分片上传，返回 uploadID。
func (a *Aliyun) InitiateMultipartUpload(objectKey string) (string, error) {
	imur, err := a.bucket.InitiateMultipartUpload(objectKey)
	if err != nil {
		return "", fmt.Errorf("initiate multipart upload: %w", err)
	}
	return imur.UploadID, nil
}

// PresignUploadPart 为指定分片生成预签名 PUT URL（partNumber 从 1 开始）。
func (a *Aliyun) PresignUploadPart(uploadID, objectKey string, partNumber int, expiry time.Duration) (string, error) {
	url, err := a.bucket.SignURL(objectKey, oss.HTTPPut, int64(expiry.Seconds()),
		oss.AddParam("uploadId", uploadID),
		oss.AddParam("partNumber", strconv.Itoa(partNumber)),
	)
	if err != nil {
		return "", fmt.Errorf("presign upload part %d: %w", partNumber, err)
	}
	return url, nil
}

// CompleteMultipartUpload 通知 OSS 合并所有已上传分片。
func (a *Aliyun) CompleteMultipartUpload(uploadID, objectKey string, parts []CompletePart) error {
	ossParts := make([]oss.UploadPart, 0, len(parts))
	for _, p := range parts {
		ossParts = append(ossParts, oss.UploadPart{
			PartNumber: p.PartNumber,
			ETag:       p.ETag,
		})
	}

	imur := oss.InitiateMultipartUploadResult{
		Bucket:   a.bucket.BucketName,
		Key:      objectKey,
		UploadID: uploadID,
	}
	_, err := a.bucket.CompleteMultipartUpload(imur, ossParts)
	return err
}
