package oss

import (
	"errors"
	"fmt"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
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
		c := config
		client, err := oss.New(c.Endpoint, c.KeyID, c.KeySecret)
		if err != nil {
			return err
		}

		bucket, err := client.Bucket(c.Bucket)
		if err != nil {
			return err
		}
		a.bucket = bucket
	}

	return nil
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
	fmt.Println("GetObjectUrl")

	fmt.Println(objectKey)

	if a.config.Domain == "" {
		url, err := a.bucket.SignURL(objectKey, oss.HTTPGet, 1800)
		if err != nil {
			return fmt.Sprintf("https://%s.%s/%s",
				a.config.Bucket,
				a.config.Endpoint,
				objectKey,
			)
		}
		return url
	}
	return fmt.Sprintf("https://%s/%s", a.config.Domain, objectKey)
}
