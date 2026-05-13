package main

import (
	"flag"

	"go.uber.org/zap"
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
	clearCache := flag.Bool("clear-cache", false, "是否清空所有Redis缓存")
	migrateDanmaku := flag.Bool("migrate-danmaku", false, "迁移弹幕和历史记录的resource_short_id")
	// 保留 -api 仅兼容旧命令行；API 与 Casbin 规则已在 InitDefaultData→SyncApiData 中每次启动自动同步
	_ = flag.Bool("api", false, "兼容旧参数（无额外效果）")
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
	// 初始化GPU检测

	// 初始化mysql
	global.Mysql = mysql.Init(global.Config.Mysql)
	initialize.InitTables()
	initialize.InitDefaultData()
	// 初始化分区数据
	global.VideoPartitionMap = service.GetPartitionMap(global.CONTENT_TYPE_VIDEO)
	// 初始化缓存
	global.Redis = redis.Init(global.Config.Redis)

	// 如果指定了clear-cache参数，清空所有Redis缓存
	if *clearCache {
		global.Redis.FlushDB()
		zap.L().Info("已清空所有Redis缓存", zap.String("module", "cache"))
	}

	// 如果指定了migrate-danmaku参数，执行弹幕和历史记录的resource_short_id迁移
	if *migrateDanmaku {
		zap.L().Info("开始迁移弹幕和历史记录的resource_short_id...", zap.String("module", "migrate"))
		service.MigrateAllDanmakuAndHistoryResourceShortID()
		zap.L().Info("迁移完成", zap.String("module", "migrate"))
	}

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
