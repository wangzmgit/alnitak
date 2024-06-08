package config

import (
	"github.com/spf13/viper"
	"interastral-peace.com/alnitak/utils"
)

// 初始化配置文件
func InitConfig(env string) {
	viper.SetConfigName(getConfigFileName(env))
	viper.AddConfigPath("./conf")
	if err := viper.ReadInConfig(); err != nil {
		panic("配置文件读取失败" + err.Error())
	}

	initDefaultConfigItems()
}

// 获取配置文件名
func getConfigFileName(env string) string {
	switch env {
	case "dev":
		return "application.dev"
	case "prod":
		return "application.prod"
	default:
		return "application.dev"
	}
}

// 初始化默认配置
func initDefaultConfigItems() {
	if viper.GetString("security.access_jwt_secret") == "" {
		viper.Set("server.access_jwt_secret", utils.GenerateNumberCode(16))
	}
	if viper.GetString("security.refresh_jwt_secret") == "" {
		viper.Set("server.refresh_jwt_secret", utils.GenerateNumberCode(16))
	}
	if viper.GetString("security.user_id_salt") == "" {
		viper.Set("server.user_id_salt", utils.GenerateNumberCode(16))
	}

	viper.WriteConfig()
}
