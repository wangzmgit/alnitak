package mysql

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"interastral-peace.com/alnitak/utils"
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
		utils.ErrorLog("mysql连接失败", "db", err.Error())
		panic(err)
	} else {
		zap.L().Info("mysql连接成功", zap.String("module", "db"))
		db = mysqlClient
		return mysqlClient
	}
}

func GetMysqlClient() *gorm.DB {
	return db
}
