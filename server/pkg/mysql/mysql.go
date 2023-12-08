package mysql

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

var db *gorm.DB

func Init() *gorm.DB {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.datasource"),
		viper.GetString("mysql.param"))

	logger := zapgorm2.New(zap.L())
	logger.SetAsDefault()
	if mysqlClient, err := gorm.Open(mysql.Open(dns), &gorm.Config{Logger: logger}); err != nil {
		zap.L().Error("mysql连接失败" + err.Error())
		panic(err)
	} else {
		zap.L().Info("mysql连接成功")
		db = mysqlClient
		return mysqlClient
	}
}

func GetMysqlClient() *gorm.DB {
	return db
}
