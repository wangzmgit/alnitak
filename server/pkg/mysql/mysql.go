package mysql

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"interastral-peace.com/alnitak/internal/config"
	"interastral-peace.com/alnitak/utils"
	"moul.io/zapgorm2"
)

var db *gorm.DB

func Init(c config.Mysql) *gorm.DB {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", c.Username, c.Password, c.Host, c.Port, c.Datasource, c.Param)

	// 配置zapgorm2，忽略ErrRecordNotFound错误
	zapLogger := zapgorm2.New(zap.L())
	zapLogger.SetAsDefault()
	zapLogger.IgnoreRecordNotFoundError = true // 关键配置：忽略ErrRecordNotFound错误

	if mysqlClient, err := gorm.Open(mysql.Open(dns), &gorm.Config{Logger: zapLogger}); err != nil {
		utils.ErrorLog("mysql连接失败", "db", err.Error())
		panic(err)
	} else {
		// 配置数据库连接池
		sqlDB, err := mysqlClient.DB()
		if err != nil {
			utils.ErrorLog("获取数据库连接失败", "db", err.Error())
			panic(err)
		}

		// 设置连接池参数
		sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
		sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
		sqlDB.SetConnMaxLifetime(time.Hour) // 设置了连接可复用的最大时间

		zap.L().Info("mysql连接成功", zap.String("module", "db"))
		db = mysqlClient
		return mysqlClient
	}
}

func GetMysqlClient() *gorm.DB {
	return db
}
