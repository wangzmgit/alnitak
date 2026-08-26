package main

import (
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/initialize"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/pkg/casbin"
	"interastral-peace.com/alnitak/pkg/jigsaw"
	"interastral-peace.com/alnitak/pkg/logger"
	"interastral-peace.com/alnitak/pkg/mysql"
	"interastral-peace.com/alnitak/pkg/oss"
	"interastral-peace.com/alnitak/pkg/redis"
)

// 一个简单的命令行工具，用于在不启动 HTTP 服务的情况下触发后端内置的重新转码逻辑。
// 使用方法示例（在项目根目录 e:\server\alnitak 下执行）:
//   go run ./server/cmd/retranscode -env=prod -vid=32
func main() {
	env := flag.String("env", "prod", "dev/prod")
	videoID := flag.Uint("vid", 0, "需要重新转码的视频ID")
	flag.Parse()

	if *videoID == 0 {
		fmt.Println("必须指定 -vid，例如: -vid=32")
		os.Exit(1)
	}

	// 初始化配置、日志、存储、数据库、缓存等，逻辑与主程序保持一致
	initialize.InitConfig(*env)
	logger.InitLogger(global.Config)
	jigsaw.Jigsaw()

	if global.Config.Storage.OssType != "local" {
		global.Storage = oss.InitStorage(global.Config.Storage)
	}
	// 初始化备用OSS（多源容灾）
	global.StorageBackup = oss.InitBackupStorage(global.Config.Storage)

	initialize.InitSnowflake()

	global.Mysql = mysql.Init(global.Config.Mysql)
	initialize.InitTables()
	initialize.InitDefaultData()

	global.VideoPartitionMap = service.GetPartitionMap(global.CONTENT_TYPE_VIDEO)
	global.Redis = redis.Init(global.Config.Redis)
	initialize.InitCacheData()
	global.Casbin = casbin.InitCasbin()

	zap.L().Info("开始执行重新转码", zap.Uint("video_id", *videoID))

	if err := service.ReTranscodeVideo(nil, uint(*videoID)); err != nil {
		zap.L().Error("重新转码触发失败", zap.Error(err))
		os.Exit(1)
	}

	zap.L().Info("重新转码任务已提交（实际转码在后台 goroutine 中执行）", zap.Uint("video_id", *videoID))
}

