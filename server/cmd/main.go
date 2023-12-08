package main

import (
	"flag"

	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/initialize"
	"interastral-peace.com/alnitak/internal/routes"
	"interastral-peace.com/alnitak/pkg/casbin"
	"interastral-peace.com/alnitak/pkg/config"
	"interastral-peace.com/alnitak/pkg/jigsaw"
	"interastral-peace.com/alnitak/pkg/logger"
	"interastral-peace.com/alnitak/pkg/mysql"
	"interastral-peace.com/alnitak/pkg/redis"
)

func main() {
	env := flag.String("env", "dev", "dev/prod")
	flag.Parse()

	// 初始化配置文件
	config.InitConfig(*env)
	// 初始化日志
	logger.InitLogger()
	// 初始化滑块验证码生成
	jigsaw.Jigsaw()
	// 初始化OSS
	// initialize.Oss()
	// 初始化mysql
	global.Mysql = mysql.Init()
	initialize.InitTables()
	// 初始化缓存
	global.Redis = redis.Init()
	// 初始化casbin
	global.Casbin = casbin.InitCasbin()

	// 初始化路由
	routes.InitRouter()
}
