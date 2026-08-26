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
		UseSSL:   global.Config.Storage.UseSSL,

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
		UseSSL:        storageConfigReq.UseSSL,
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
	viper.Set("storage.use_ssl", storageConfigReq.UseSSL)
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
	// 重新初始化备用OSS（多源容灾）
	global.StorageBackup = oss.InitBackupStorage(global.Config.Storage)

	return nil
}

func GetOtherConfig() vo.OtherConfigResp {
	return vo.OtherConfigResp{
		AllowOrigin: global.Config.Cors.AllowOrigin,
		Prefix:      global.Config.User.Prefix,

		ServerPort:  global.Config.Server.Port,
		SslEnabled:  global.Config.Server.Ssl.Enabled,
		SslPort:     global.Config.Server.Ssl.Port,
		SslCertFile: global.Config.Server.Ssl.CertFile,
		SslKeyFile:  global.Config.Server.Ssl.KeyFile,
	}
}

func SetOtherConfig(otherConfigReq dto.OtherConfigReq) error {
	oldCorsConfig := global.Config.Cors
	oldUserConfig := global.Config.User
	oldServerConfig := global.Config.Server

	global.Config.Cors = config.Cors{
		AllowOrigin: otherConfigReq.AllowOrigin,
	}
	global.Config.User = config.User{
		Prefix: otherConfigReq.Prefix,
	}
	global.Config.Server = config.Server{
		Port: otherConfigReq.ServerPort,
		Ssl: config.Ssl{
			Enabled:  otherConfigReq.SslEnabled,
			Port:     otherConfigReq.SslPort,
			CertFile: otherConfigReq.SslCertFile,
			KeyFile:  otherConfigReq.SslKeyFile,
		},
	}

	viper.Set("cors.allow_origin", otherConfigReq.AllowOrigin)
	viper.Set("user.prefix", otherConfigReq.Prefix)
	viper.Set("server.port", otherConfigReq.ServerPort)
	viper.Set("server.ssl.enabled", otherConfigReq.SslEnabled)
	viper.Set("server.ssl.port", otherConfigReq.SslPort)
	viper.Set("server.ssl.cert_file", otherConfigReq.SslCertFile)
	viper.Set("server.ssl.key_file", otherConfigReq.SslKeyFile)

	if err := viper.WriteConfig(); err != nil {
		global.Config.Cors = oldCorsConfig
		global.Config.User = oldUserConfig
		global.Config.Server = oldServerConfig
		utils.ErrorLog("写入其他配置失败", "config", err.Error())
		return errors.New("更新失败")
	}

	return nil
}

// 获取转码配置
func GetTranscodingConfig() vo.TranscodingConfigResp {
	return vo.TranscodingConfigResp{
		Mode:                global.Config.Transcoding.Mode,
		UseGpu:              global.Config.Transcoding.UseGpu,
		UseH265:             global.Config.Transcoding.UseH265,
		UseAv1:              global.Config.Transcoding.UseAv1,
		Generate1080p60:     global.Config.Transcoding.Generate1080p60,
		MaxCpuConcurrency:   global.Config.Transcoding.MaxCpuConcurrency,
		MaxGpuConcurrency:   global.Config.Transcoding.MaxGpuConcurrency,
		WorkerConcurrency:   global.Config.Transcoding.WorkerConcurrency,
		EncodingConcurrency: global.Config.Transcoding.EncodingConcurrency,
		MaxQueueDepth:       global.Config.Transcoding.MaxQueueDepth,
		WorkDir:             global.Config.Transcoding.WorkDir,
	}
}

// 修改转码配置
func SetTranscodingConfig(req dto.TranscodingConfigReq) error {
	old := global.Config.Transcoding

	global.Config.Transcoding = config.Transcoding{
		Mode:                req.Mode,
		UseGpu:              req.UseGpu,
		UseH265:             req.UseH265,
		UseAv1:              req.UseAv1,
		Generate1080p60:     req.Generate1080p60,
		MaxCpuConcurrency:   req.MaxCpuConcurrency,
		MaxGpuConcurrency:   req.MaxGpuConcurrency,
		WorkerConcurrency:   req.WorkerConcurrency,
		EncodingConcurrency: req.EncodingConcurrency,
		MaxQueueDepth:       req.MaxQueueDepth,
		WorkDir:             req.WorkDir,
	}

	viper.Set("transcoding.mode", req.Mode)
	viper.Set("transcoding.use_gpu", req.UseGpu)
	viper.Set("transcoding.use_h265", req.UseH265)
	viper.Set("transcoding.use_av1", req.UseAv1)
	viper.Set("transcoding.generate_1080p60", req.Generate1080p60)
	viper.Set("transcoding.max_cpu_concurrency", req.MaxCpuConcurrency)
	viper.Set("transcoding.max_gpu_concurrency", req.MaxGpuConcurrency)
	viper.Set("transcoding.worker_concurrency", req.WorkerConcurrency)
	viper.Set("transcoding.encoding_concurrency", req.EncodingConcurrency)
	viper.Set("transcoding.max_queue_depth", req.MaxQueueDepth)
	viper.Set("transcoding.work_dir", req.WorkDir)

	if err := viper.WriteConfig(); err != nil {
		global.Config.Transcoding = old
		utils.ErrorLog("写入转码配置失败", "config", err.Error())
		return errors.New("更新失败")
	}

	return nil
}
