package main

import (
	"flag"

	"interastral-peace.com/alnitak/internal/cron"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/initialize"
	"interastral-peace.com/alnitak/internal/routes"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/pkg/casbin"
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
	initialize.InitConfig(*env)
	// 初始化日志
	logger.InitLogger()
	// 初始化滑块验证码生成
	jigsaw.Jigsaw()
	// 初始化OSS
	if global.Config.Storage.OssType != "local" {
		global.Storage = oss.InitStorage(global.Config.Storage)
	}
	// 初始化雪花ID
	initialize.InitSnowflake()
	// 初始化mysql
	global.Mysql = mysql.Init(global.Config.Mysql)
	initialize.InitTables()
	initialize.InitDefaultData()
	// 初始化分区数据
	global.VideoPartitionMap = service.GetPartitionMap(global.CONTENT_TYPE_VIDEO)
	// 初始化缓存
	global.Redis = redis.Init(global.Config.Redis)
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
