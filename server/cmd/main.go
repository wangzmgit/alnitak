package main

import (
	"flag"

	"github.com/spf13/viper"
	"interastral-peace.com/alnitak/internal/cron"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/initialize"
	"interastral-peace.com/alnitak/internal/routes"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/pkg/casbin"
	"interastral-peace.com/alnitak/pkg/config"
	"interastral-peace.com/alnitak/pkg/jigsaw"
	"interastral-peace.com/alnitak/pkg/logger"
	"interastral-peace.com/alnitak/pkg/mysql"
	"interastral-peace.com/alnitak/pkg/oss"
	"interastral-peace.com/alnitak/pkg/redis"
)

func main() {
	env := flag.String("env", "prod", "dev/prod")
	flag.Parse()

	// 初始化配置文件
	config.InitConfig(*env)
	// 初始化日志
	logger.InitLogger()
	// 初始化滑块验证码生成
	jigsaw.Jigsaw()
	// 初始化OSS
	if viper.GetString("storage.oss_type") != "local" {
		global.Storage = oss.InitStorage()
	}
	// 初始化雪花ID
	initialize.InitSnowflake()
	// 初始化mysql
	global.Mysql = mysql.Init()
	initialize.InitTables()
	initialize.InitDefaultData()
	// 初始化分区数据
	global.VideoPartitionMap = service.GetPartitionMap(global.CONTENT_TYPE_VIDEO)
	// 初始化缓存
	global.Redis = redis.Init()
	initialize.InitCacheData()
	// 初始化casbin
	global.Casbin = casbin.InitCasbin()

	// 手动执行一次刷新热点视频
	cron.RefreshPopular()
	// 启动定时任务
	go cron.StartCronTask()

	// 初始化路由
	routes.InitRouter()
}
