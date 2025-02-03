package service

import (
	"errors"

	"github.com/spf13/viper"
	"interastral-peace.com/alnitak/internal/config"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/pkg/oss"
	"interastral-peace.com/alnitak/utils"
)

// 获取邮箱配置信息
func GetEmailConfig() vo.EmailConfigResp {
	return vo.EmailConfigResp{
		User:      global.Config.Mail.User,
		Host:      global.Config.Mail.Host,
		Port:      global.Config.Mail.Port,
		Addresser: global.Config.Mail.Addresser,
	}
}

func SetEmailConfig(emailConfigReq dto.EmailConfigReq) error {
	oldConfig := global.Config.Mail

	global.Config.Mail = config.Mail{
		User:      emailConfigReq.User,
		Host:      emailConfigReq.Host,
		Port:      emailConfigReq.Port,
		Addresser: emailConfigReq.Addresser,
		Pass:      oldConfig.Pass,
		Debug:     oldConfig.Debug,
	}

	viper.Set("mail.user", emailConfigReq.User)
	viper.Set("mail.host", emailConfigReq.Host)
	viper.Set("mail.port", emailConfigReq.Port)
	viper.Set("mail.addresser", emailConfigReq.Addresser)

	if len(emailConfigReq.Pass) != 0 {
		global.Config.Mail.Pass = emailConfigReq.Pass
		viper.Set("mail.pass", emailConfigReq.Pass)
	}

	if err := viper.WriteConfig(); err != nil {
		global.Config.Mail = oldConfig
		utils.ErrorLog("写入邮箱配置失败", "config", err.Error())
		return errors.New("更新失败")
	}

	return nil
}

func GetStorageConfig() vo.StorageConfigResp {
	return vo.StorageConfigResp{
		MaxImgSize:   global.Config.File.MaxImgSize,
		MaxVideoSize: global.Config.File.MaxVideoSize,

		Type:     global.Config.Storage.OssType,
		KeyID:    global.Config.Storage.KeyId,
		Bucket:   global.Config.Storage.Bucket,
		Endpoint: global.Config.Storage.Endpoint,
		AppID:    global.Config.Storage.AppId,
		Region:   global.Config.Storage.Region,
		Domain:   global.Config.Storage.Domain,
		Private:  global.Config.Storage.Private,

		UploadMp4File: global.Config.Storage.UploadMp4File,
	}
}

func SetStorageConfig(storageConfigReq dto.StorageConfigReq) error {
	oldFileConfig := global.Config.File
	oldStorageConfig := global.Config.Storage

	global.Config.File = config.File{
		MaxImgSize:   storageConfigReq.MaxImgSize,
		MaxVideoSize: storageConfigReq.MaxVideoSize,
	}

	global.Config.Storage = config.Storage{
		OssType:       storageConfigReq.Type,
		KeyId:         storageConfigReq.KeyID,
		Bucket:        storageConfigReq.Bucket,
		Endpoint:      storageConfigReq.Endpoint,
		AppId:         storageConfigReq.AppID,
		Region:        storageConfigReq.Region,
		Domain:        storageConfigReq.Domain,
		Private:       storageConfigReq.Private,
		UploadMp4File: storageConfigReq.UploadMp4File,
	}

	viper.Set("file.max_img_size", storageConfigReq.MaxImgSize)
	viper.Set("file.max_video_size", storageConfigReq.MaxVideoSize)
	viper.Set("storage.oss_type", storageConfigReq.Type)
	viper.Set("storage.key_id", storageConfigReq.KeyID)
	viper.Set("storage.bucket", storageConfigReq.Bucket)
	viper.Set("storage.endpoint", storageConfigReq.Endpoint)
	viper.Set("storage.app_id", storageConfigReq.AppID)
	viper.Set("storage.region", storageConfigReq.Region)
	viper.Set("storage.domain", storageConfigReq.Domain)
	viper.Set("storage.private", storageConfigReq.Private)
	viper.Set("storage.upload_mp4_file", storageConfigReq.UploadMp4File)

	if len(storageConfigReq.KeySecret) != 0 {
		global.Config.Storage.KeySecret = storageConfigReq.KeySecret
		viper.Set("storage.key_secret", storageConfigReq.KeySecret)
	}

	if err := viper.WriteConfig(); err != nil {
		global.Config.File = oldFileConfig
		global.Config.Storage = oldStorageConfig
		utils.ErrorLog("写入存储配置失败", "config", err.Error())
		return errors.New("更新失败")
	}

	// 重新初始化OSS
	if storageConfigReq.Type != "local" {
		global.Storage = oss.InitStorage(global.Config.Storage)
	}

	return nil
}

func GetOtherConfig() vo.OtherConfigResp {
	return vo.OtherConfigResp{
		AllowOrigin: viper.GetString("cors.allow_origin"),
		Prefix:      viper.GetString("user.prefix"),
	}
}

func SetOtherConfig(otherConfigReq dto.OtherConfigReq) error {
	oldCorsConfig := global.Config.Cors
	oldUserConfig := global.Config.User

	global.Config.Cors = config.Cors{
		AllowOrigin: otherConfigReq.AllowOrigin,
	}
	global.Config.User = config.User{
		Prefix: otherConfigReq.Prefix,
	}

	viper.Set("cors.allow_origin", otherConfigReq.AllowOrigin)
	viper.Set("user.prefix", otherConfigReq.Prefix)

	if err := viper.WriteConfig(); err != nil {
		global.Config.Cors = oldCorsConfig
		global.Config.User = oldUserConfig
		utils.ErrorLog("写入存其他置失败", "config", err.Error())
		return errors.New("更新失败")
	}

	return nil
}
