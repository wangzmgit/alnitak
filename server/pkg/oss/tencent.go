package oss

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
	"interastral-peace.com/alnitak/utils"
)

type TencentCOS struct {
	config *Config
	client *cos.Client
}

func newTencentCOS(config Config) (Storage, error) {
	tencentCOS := TencentCOS{}
	err := tencentCOS.init(config)
	if err != nil {
		return nil, err
	}
	return &tencentCOS, nil
}

func (t *TencentCOS) init(config Config) error {
	if t.config == nil {
		t.config = &config
	}

	// 检查配置项是否正确
	if config.Bucket == "" || config.KeyID == "" || config.KeySecret == "" || config.Region == "" {
		return errors.New("配置不正确")
	}

	// 使用 BaseURL 和授权信息初始化 COS 客户端
	u, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", config.Bucket, config.Region))
	if err != nil {
		return err
	}

	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.KeyID,
			SecretKey: config.KeySecret,
		},
	})
	t.client = client

	return nil
}

func (t *TencentCOS) GetObjectToFile(objectKey, filePath string) error {
	// 获取腾讯云 COS 存储中的文件
	resp, err := t.client.Object.Get(context.Background(), objectKey, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建本地文件准备写入
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将腾讯云 COS 的内容复制到本地文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (t *TencentCOS) DeleteObject(objectKey string) error {
	// 删除腾讯云 COS 存储中的文件
	_, err := t.client.Object.Delete(context.Background(), objectKey)
	if err != nil {
		return err
	}
	return nil
}

func (t *TencentCOS) PutObject(objectKey string, reader io.Reader) error {
	// 将文件上传到腾讯云 COS 存储
	_, err := t.client.Object.Put(context.Background(), objectKey, reader, nil)
	if err != nil {
		return err
	}
	return nil
}

func (t *TencentCOS) PutObjectFromFile(objectKey, filePath string) error {
	// 从文件直接上传到腾讯云 COS 存储
	_, err := t.client.Object.PutFromFile(context.Background(), objectKey, filePath, nil)
	if err != nil {
		return err
	}
	return nil
}

func (t *TencentCOS) IsExists(objectKey string) (bool, error) {
	// 检查文件是否存在
	_, err := t.client.Object.Get(context.Background(), objectKey, nil)
	if err != nil {
		// 如果返回文件不存在错误，则返回 false
		if cosError, ok := err.(*cos.ErrorResponse); ok && cosError.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (t *TencentCOS) GetObjectUrl(objectKey string) string {
	// 获取腾讯云 COS 的文件签名 URL，指定有效期为 24 小时
	presignedURL, err := t.client.Object.GetPresignedURL(
		context.Background(), // 上下文
		"GET",                // 请求方法
		objectKey,            // 文件对象的 key
		t.config.KeyID,       // SecretID
		t.config.KeySecret,   // SecretKey
		time.Hour*5,          // 签名有效期 5 小时
		nil,                  // 可选的额外参数，可以为 nil
	)
	if err != nil {
		// 记录错误并显示更多详细信息
		utils.ErrorLog("腾讯云COS生成文件URL失败", "transcoding", fmt.Sprintf("Error: %v, ObjectKey: %s", err, objectKey))
		return ""
	}
	// GetPresignedURL 返回 *url.URL，直接 .String()
	url := presignedURL.String()
	if url == "" {
		utils.ErrorLog("腾讯云COS生成文件URL失败", "transcoding", "empty url")
	}
	return url
}

// 获取对象读取器（用于代理流式传输）
func (t *TencentCOS) GetObjectReader(objectKey string) (io.ReadCloser, error) {
	resp, err := t.client.Object.Get(context.Background(), objectKey, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// ── 直传 OSS 预签名方法 ──

// PresignPutObject 生成单文件上传预签名 URL（用于图片直传）。
func (t *TencentCOS) PresignPutObject(objectKey string, expiry time.Duration) (string, error) {
	presignedURL, err := t.client.Object.GetPresignedURL(
		context.Background(),
		http.MethodPut,
		objectKey,
		t.config.KeyID,
		t.config.KeySecret,
		expiry,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("presign put object: %w", err)
	}
	return presignedURL.String(), nil
}

// InitiateMultipartUpload 发起分片上传，返回 uploadID。
func (t *TencentCOS) InitiateMultipartUpload(objectKey string) (string, error) {
	resp, _, err := t.client.Object.InitiateMultipartUpload(context.Background(), objectKey, nil)
	if err != nil {
		return "", fmt.Errorf("initiate multipart upload: %w", err)
	}
	return resp.UploadID, nil
}

// PresignUploadPart 为指定分片生成预签名 PUT URL（partNumber 从 1 开始）。
// COS 需要 uploadId 和 partNumber 作为 query 参数参与签名，
// 通过 cos-go-sdk-v5 的 GetPresignedURL 最后一个 url.Values 参数传入，
// 签名时会将这些 query 纳入签名计算。
func (t *TencentCOS) PresignUploadPart(uploadID, objectKey string, partNumber int, expiry time.Duration) (string, error) {
	query := &url.Values{}
	query.Set("uploadId", uploadID)
	query.Set("partNumber", fmt.Sprintf("%d", partNumber))

	// GetPresignedURL 最后一个可选参数如果是 *url.Values，SDK 会把它纳入签名
	presignedURL, err := t.client.Object.GetPresignedURL(
		context.Background(),
		http.MethodPut,
		objectKey,
		t.config.KeyID,
		t.config.KeySecret,
		expiry,
		query, // 作为 cos.ObjectPutHeaderOptions 的接口实现，实际是 url.Values
	)
	if err != nil {
		return "", fmt.Errorf("presign upload part %d: %w", partNumber, err)
	}
	return presignedURL.String(), nil
}

// CompleteMultipartUpload 通知 OSS 合并所有已上传分片。
func (t *TencentCOS) CompleteMultipartUpload(uploadID, objectKey string, parts []CompletePart) error {
	opt := &cos.CompleteMultipartUploadOptions{
		Parts: make([]cos.Object, 0, len(parts)),
	}
	for _, p := range parts {
		opt.Parts = append(opt.Parts, cos.Object{
			PartNumber: p.PartNumber,
			ETag:       p.ETag,
		})
	}

	_, _, err := t.client.Object.CompleteMultipartUpload(context.Background(), objectKey, uploadID, opt)
	if err != nil {
		return fmt.Errorf("complete multipart upload: %w", err)
	}
	return nil
}

