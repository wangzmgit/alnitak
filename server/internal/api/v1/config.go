package api

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/pkg/oss"
)

// 获取邮箱配置信息
func GetEmailConfig(ctx *gin.Context) {
	config := vo.EmailConfigResp{
		User: viper.GetString("mail.user"),
		// Pass:      viper.GetString("mail.pass"),
		Host:      viper.GetString("mail.host"),
		Port:      viper.GetInt("mail.port"),
		Addresser: viper.GetString("mail.addresser"),
	}

	resp.OkWithData(ctx, gin.H{"config": config})
}

// 配置邮箱
func SetEmailConfig(ctx *gin.Context) {
	var emailConfigReq dto.EmailConfigReq
	if err := ctx.Bind(&emailConfigReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	viper.Set("mail.user", emailConfigReq.User)
	viper.Set("mail.host", emailConfigReq.Host)
	viper.Set("mail.port", emailConfigReq.Port)
	viper.Set("mail.addresser", emailConfigReq.Addresser)

	if len(emailConfigReq.Pass) != 0 {
		viper.Set("mail.pass", emailConfigReq.Pass)
	}

	viper.WriteConfig()

	// 返回给前端
	resp.Ok(ctx)
}

// 获取存储配置信息
func GetStorageConfig(ctx *gin.Context) {
	config := vo.StorageConfigResp{
		MaxImgSize:   viper.GetInt("file.max_img_size"),
		MaxVideoSize: viper.GetInt("file.max_video_size"),

		Type:     viper.GetString("storage.oss_type"),
		KeyID:    viper.GetString("storage.key_id"),
		Bucket:   viper.GetString("storage.bucket"),
		Endpoint: viper.GetString("storage.endpoint"),
		AppID:    viper.GetString("storage.app_id"),
		Region:   viper.GetString("storage.region"),
		Domain:   viper.GetString("storage.domain"),
		Private:  viper.GetBool("storage.private"),

		UploadMp4File: viper.GetBool("storage.upload_mp4_file"),
	}

	resp.OkWithData(ctx, gin.H{"config": config})
}

// 修改存储配置
func SetStorageConfig(ctx *gin.Context) {
	var storageConfigReq dto.StorageConfigReq
	if err := ctx.Bind(&storageConfigReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
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
		viper.Set("storage.key_secret", storageConfigReq.KeySecret)
	}

	viper.WriteConfig()

	// 重新初始化OSS
	if storageConfigReq.Type != "local" {
		global.Storage = oss.InitStorage()
	}

	// 返回给前端
	resp.Ok(ctx)
}

// 获取其他配置信息
func GetOtherConfig(ctx *gin.Context) {
	config := vo.OtherConfigResp{
		AllowOrigin: viper.GetString("cors.allow_origin"),
		Prefix:      viper.GetString("user.prefix"),
	}

	resp.OkWithData(ctx, gin.H{"config": config})
}

// 修改其他配置信息
func SetOtherConfig(ctx *gin.Context) {
	var otherConfigReq dto.OtherConfigReq
	if err := ctx.Bind(&otherConfigReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	viper.Set("cors.allow_origin", otherConfigReq.AllowOrigin)
	viper.Set("user.prefix", otherConfigReq.Prefix)

	viper.WriteConfig()

	// 返回给前端
	resp.Ok(ctx)
}
