package oss

import (
	"errors"
	"io"

	"github.com/spf13/viper"
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

func InitStorage() Storage {
	ossName := viper.GetString("storage.oss_type")
	config := Config{
		KeyID:     viper.GetString("storage.key_id"),
		KeySecret: viper.GetString("storage.key_secret"),
		Bucket:    viper.GetString("storage.bucket"),
		Endpoint:  viper.GetString("storage.endpoint"),
		AppID:     viper.GetString("storage.app_id"),
		Region:    viper.GetString("storage.region"),
		Domain:    viper.GetString("storage.domain"),
		Private:   viper.GetBool("storage.private"),
	}

	s, err := initOss(ossName, config)
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
