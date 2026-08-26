package oss

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscrd "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	miniosdk "github.com/minio/minio-go/v7"
	minicred "github.com/minio/minio-go/v7/pkg/credentials"
	"interastral-peace.com/alnitak/utils"
)

type MinIO struct {
	config   *Config
	client   *miniosdk.Client
	s3Client *s3.Client // 用于预签名
	bucket   *string
}

func newMinio(config Config) (Storage, error) {
	minio := MinIO{}
	err := minio.init(config)
	if err != nil {
		return nil, err
	}
	return &minio, nil
}

func (m *MinIO) init(config Config) error {
	if m.config == nil {
		m.config = &config
	}

	if config.Endpoint == "" {
		return errors.New("configuration not correct")
	}

	if m.client == nil {
		client, err := m.initMinIOClient(config)
		if err != nil {
			return err
		}
		m.client = client

		// 同时创建 S3 客户端用于预签名
		m.s3Client = m.initS3Client(config)
	}

	if m.bucket == nil {
		exists, err := m.client.BucketExists(context.Background(), config.Bucket)
		if err != nil {
			return err
		}
		if !exists {
			err = m.client.MakeBucket(context.Background(), config.Bucket, miniosdk.MakeBucketOptions{})
			if err != nil {
				return err
			}
		}
		m.bucket = &config.Bucket
	}

	return nil
}

func (m *MinIO) initMinIOClient(config Config) (*miniosdk.Client, error) {
	options := miniosdk.Options{
		Creds:  minicred.NewStaticV4(config.KeyID, config.KeySecret, ""),
		Secure: config.UseSSL,
	}
	if config.Timeout > 0 {
		options.Transport = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   120 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: time.Duration(config.Timeout) * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
	}
	client, err := miniosdk.New(config.Endpoint, &options)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// initS3Client 创建 AWS S3 兼容客户端，用于预签名操作。
func (m *MinIO) initS3Client(config Config) *s3.Client {
	scheme := "https"
	if !config.UseSSL {
		scheme = "http"
	}
	endpoint := fmt.Sprintf("%s://%s", scheme, config.Endpoint)

	cfg := aws.Config{
		Credentials: awscrd.NewStaticCredentialsProvider(config.KeyID, config.KeySecret, ""),
		Region:      "us-east-1",
	}
	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})
}

// 获取文件
func (m *MinIO) GetObjectToFile(objectKey, filePath string) error {
	object, err := m.client.GetObject(context.Background(), m.config.Bucket, objectKey, miniosdk.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer object.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, object)
	return err
}

// 删除文件
func (m *MinIO) DeleteObject(objectKey string) error {
	err := m.client.RemoveObject(context.Background(), m.config.Bucket, objectKey, miniosdk.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// 上传文件
func (m *MinIO) PutObject(objectKey string, reader io.Reader) error {
	_, err := m.client.PutObject(context.Background(), m.config.Bucket, objectKey, reader, -1, miniosdk.PutObjectOptions{})
	return err
}

// 上传文件（从文件）
func (m *MinIO) PutObjectFromFile(objectKey, filePath string) error {
	_, err := m.client.FPutObject(context.Background(), m.config.Bucket, objectKey, filePath, miniosdk.PutObjectOptions{})
	return err
}

// 检查文件是否存在
func (m *MinIO) IsExists(objectKey string) (bool, error) {
	_, err := m.client.StatObject(context.Background(), m.config.Bucket, objectKey, miniosdk.StatObjectOptions{})
	if err != nil {
		if errResp, ok := err.(miniosdk.ErrorResponse); ok && errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// 获取访问URL
func (m *MinIO) GetObjectUrl(objectKey string) string {
	presignedURL, err := m.client.PresignedGetObject(context.Background(), m.config.Bucket, objectKey, 24*time.Hour, nil)
	if err != nil {
		utils.ErrorLog("MinIO生成文件URL失败", "transcoding", fmt.Sprintf("Error: %v, ObjectKey: %s", err, objectKey))
		return ""
	}
	return presignedURL.String()
}

// 获取对象读取器（用于代理流式传输）
func (m *MinIO) GetObjectReader(objectKey string) (io.ReadCloser, error) {
	object, err := m.client.GetObject(context.Background(), m.config.Bucket, objectKey, miniosdk.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

// ── 直传 OSS 预签名方法（使用 S3 兼容客户端） ──

// PresignPutObject 生成单文件上传预签名 URL（用于图片直传）。
func (m *MinIO) PresignPutObject(objectKey string, expiry time.Duration) (string, error) {
	signer := s3.NewPresignClient(m.s3Client)
	resp, err := signer.PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(m.config.Bucket),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", fmt.Errorf("presign put object: %w", err)
	}
	return resp.URL, nil
}

// InitiateMultipartUpload 发起分片上传，返回 uploadID。
func (m *MinIO) InitiateMultipartUpload(objectKey string) (string, error) {
	resp, err := m.s3Client.CreateMultipartUpload(context.Background(), &s3.CreateMultipartUploadInput{
		Bucket: aws.String(m.config.Bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return "", fmt.Errorf("initiate multipart upload: %w", err)
	}
	if resp.UploadId == nil {
		return "", fmt.Errorf("initiate multipart upload: nil uploadId")
	}
	return *resp.UploadId, nil
}

// PresignUploadPart 为指定分片生成预签名 PUT URL（partNumber 从 1 开始）。
func (m *MinIO) PresignUploadPart(uploadID, objectKey string, partNumber int, expiry time.Duration) (string, error) {
	signer := s3.NewPresignClient(m.s3Client)
	resp, err := signer.PresignUploadPart(context.Background(), &s3.UploadPartInput{
		Bucket:     aws.String(m.config.Bucket),
		Key:        aws.String(objectKey),
		PartNumber: aws.Int32(int32(partNumber)),
		UploadId:   aws.String(uploadID),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", fmt.Errorf("presign upload part %d: %w", partNumber, err)
	}
	return resp.URL, nil
}

// CompleteMultipartUpload 通知 OSS 合并所有已上传分片。
func (m *MinIO) CompleteMultipartUpload(uploadID, objectKey string, parts []CompletePart) error {
	completedParts := make([]types.CompletedPart, 0, len(parts))
	for _, p := range parts {
		completedParts = append(completedParts, types.CompletedPart{
			PartNumber: aws.Int32(int32(p.PartNumber)),
			ETag:       aws.String(p.ETag),
		})
	}

	_, err := m.s3Client.CompleteMultipartUpload(context.Background(), &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(m.config.Bucket),
		Key:      aws.String(objectKey),
		UploadId: aws.String(uploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	if err != nil {
		return fmt.Errorf("complete multipart upload: %w", err)
	}
	return nil
}
