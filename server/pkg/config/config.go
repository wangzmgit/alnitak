package config

import (
	"github.com/spf13/viper"
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
	viper.WriteConfig()
}
