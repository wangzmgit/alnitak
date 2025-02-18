package oss

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"interastral-peace.com/alnitak/utils"
)

type MinIO struct {
	config *Config
	client *minio.Client
	bucket *string
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
	}

	if m.bucket == nil {
		exists, err := m.client.BucketExists(context.Background(), config.Bucket)
		if err != nil {
			return err
		}
		if !exists {
			err = m.client.MakeBucket(context.Background(), config.Bucket, minio.MakeBucketOptions{})
			if err != nil {
				return err
			}
		}
		m.bucket = &config.Bucket
	}

	return nil
}

func (m *MinIO) initMinIOClient(config Config) (*minio.Client, error) {
	// Use the endpoint and access/secret key to initialize the MinIO client
	options := minio.Options{
		Creds:  credentials.NewStaticV4(config.KeyID, config.KeySecret, ""),
		Secure: config.Private,
	}
	client, err := minio.New(config.Endpoint, &options)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// 获取文件
func (m *MinIO) GetObjectToFile(objectKey, filePath string) error {
	object, err := m.client.GetObject(context.Background(), m.config.Bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer object.Close()

	// Create the file to write to
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the object to the file
	_, err = io.Copy(file, object)
	if err != nil {
		return err
	}

	return nil
}

// 删除文件
func (m *MinIO) DeleteObject(objectKey string) error {
	err := m.client.RemoveObject(context.Background(), m.config.Bucket, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// 上传文件
func (m *MinIO) PutObject(objectKey string, reader io.Reader) error {
	_, err := m.client.PutObject(context.Background(), m.config.Bucket, objectKey, reader, -1, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// 上传文件（从文件）
func (m *MinIO) PutObjectFromFile(objectKey, filePath string) error {
	_, err := m.client.FPutObject(context.Background(), m.config.Bucket, objectKey, filePath, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// 检查文件是否存在
func (m *MinIO) IsExists(objectKey string) (bool, error) {
	_, err := m.client.StatObject(context.Background(), m.config.Bucket, objectKey, minio.StatObjectOptions{})
	if err != nil {
		if errResp, ok := err.(minio.ErrorResponse); ok && errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// 获取访问URL
func (m *MinIO) GetObjectUrl(objectKey string) string {
	// Generating a presigned URL for the object
	presignedURL, err := m.client.PresignedGetObject(context.Background(), m.config.Bucket, objectKey, time.Hour, nil)
	if err != nil {
		utils.ErrorLog("MinIO生成文件URL失败", "transcoding", err.Error())
		return ""
	}
	return presignedURL.String()
}
