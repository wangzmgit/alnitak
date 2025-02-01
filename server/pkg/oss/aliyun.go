package oss

import (
	"errors"
	"io"

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
	if config.Domain == "" {
		return oss.New(config.Endpoint, config.KeyID, config.KeySecret)
	} else {
		return oss.New(config.Domain, config.KeyID, config.KeySecret, oss.UseCname(true))
	}
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
	url, err := a.bucket.SignURL(objectKey, oss.HTTPGet, 1800)
	if err != nil {
		utils.ErrorLog("OSS生成文件URL失败", "transcoding", err.Error())
		return ""
	}

	return url
}
