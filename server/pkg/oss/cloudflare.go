package oss

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	//"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"interastral-peace.com/alnitak/utils"
)

type MinIOStorage struct {
	config *Config
	client *s3.Client
	bucket string
}

func newR2Storage(config Config) (Storage, error) {
	minio := MinIOStorage{}
	err := minio.init(config)
	if err != nil {
		return nil, err
	}
	return &minio, nil
}

func (m *MinIOStorage) init(config Config) error {
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

	m.bucket = config.Bucket
	return nil
}

func (m *MinIOStorage) initMinIOClient(config Config) (*s3.Client, error) {
	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               config.Endpoint,
			HostnameImmutable: true,
			SigningRegion:     "us-east-1",
		}, nil
	})

	cfg := aws.Config{
		Credentials:                 credentials.NewStaticCredentialsProvider(config.KeyID, config.KeySecret, ""),
		Region:                      "us-east-1",
		EndpointResolverWithOptions: resolver,
	}

	// Force path style for MinIO compatibility
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	return client, nil
}

// 获取文件
func (m *MinIOStorage) GetObjectToFile(objectKey, filePath string) error {
	output, err := m.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return err
	}
	defer output.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, output.Body)
	return err
}

// 删除文件
func (m *MinIOStorage) DeleteObject(objectKey string) error {
	_, err := m.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(objectKey),
	})
	return err
}

// 上传文件
func (m *MinIOStorage) PutObject(objectKey string, reader io.Reader) error {
	uploader := manager.NewUploader(m.client)
	_, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(objectKey),
		Body:   reader,
	})
	return err
}

// 上传文件（从文件）
func (m *MinIOStorage) PutObjectFromFile(objectKey, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	uploader := manager.NewUploader(m.client)
	_, err = uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	return err
}

// 检查文件是否存在
func (m *MinIOStorage) IsExists(objectKey string) (bool, error) {
	_, err := m.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// 获取访问URL
func (m *MinIOStorage) GetObjectUrl(objectKey string) string {
	signer := s3.NewPresignClient(m.client)
	presignedURL, err := signer.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(time.Hour*24))
	if err != nil {
		utils.ErrorLog("MinIO 生成文件URL失败", "transcoding", fmt.Sprintf("Error: %v, ObjectKey: %s", err, objectKey))
		return ""
	}
	return presignedURL.URL
}
