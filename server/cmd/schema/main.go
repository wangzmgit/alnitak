// schema 仅连接数据库并执行 initialize.InitTables()（GORM AutoMigrate + shortId 补填等），不启动 HTTP。
// 用法（在 server 目录下，保证 ./conf 存在）：
//
//	go run ./cmd/schema -env=dev
//	go run ./cmd/schema -env=prod
package main

import (
	"flag"
	"fmt"
	"os"

	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/initialize"
	"interastral-peace.com/alnitak/pkg/logger"
	"interastral-peace.com/alnitak/pkg/mysql"
)

func main() {
	env := flag.String("env", "dev", "配置环境：dev / prod（对应 conf/application.{env}.yaml）")
	flag.Parse()

	initialize.InitConfig(*env)
	logger.InitLogger(global.Config)
	global.Mysql = mysql.Init(global.Config.Mysql)
	initialize.InitSnowflake()
	initialize.InitTables()

	fmt.Println("OK: 数据库结构已同步（AutoMigrate），并完成 InitTables 中的补填逻辑。")
	os.Exit(0)
}
