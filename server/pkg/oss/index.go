package oss

import (
	"errors"
	"io"

	"interastral-peace.com/alnitak/internal/config"
	"interastral-peace.com/alnitak/utils"
)

const (
	ALIYUN = "aliyun"
)

type Storage interface {
	GetObjectToFile(objectKey, downloadedFileName string) error
	DeleteObject(objectKey string) error
	PutObject(objectKey string, reader io.Reader) error
	PutObjectFromFile(objectKey, filePath string) error
	IsExists(objectKey string) (bool, error)
	GetObjectUrl(objectKey string) string
}

func InitStorage(c config.Storage) Storage {
	config := Config{
		KeyID:     c.KeyId,
		KeySecret: c.KeySecret,
		Bucket:    c.Bucket,
		Endpoint:  c.Endpoint,
		AppID:     c.AppId,
		Region:    c.Region,
		Domain:    c.Domain,
		Private:   c.Private,
	}

	s, err := initOss(c.OssType, config)
	if err != nil {
		utils.ErrorLog("oss初始化失败", "oss", err.Error())
		panic(err)
	}

	return s
}

func initOss(ossName string, config Config) (Storage, error) {
	switch ossName {
	case ALIYUN:
		return newAliyun(config)
	default:
		return nil, errors.New("driver not exists")
	}
}
