// db_inspect 读取 conf 配置连接数据库，输出排障所需的关键查询结果。
//
// 用法（在 server 目录下运行，保证 ./conf 存在）：
//   go run ./cmd/db_inspect -env=dev -vid=73
//   go run ./cmd/db_inspect -env=dev -pgc_id=2038788819825725440
//   go run ./cmd/db_inspect -env=dev -vid=73 -pgc_id=2038788819825725440
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/initialize"
	"interastral-peace.com/alnitak/pkg/logger"
	"interastral-peace.com/alnitak/pkg/mysql"
)

type Row map[string]any

func main() {
	env := flag.String("env", "dev", "配置环境：dev / prod（对应 conf/application.{env}.yaml）")
	vid := flag.Uint("vid", 0, "视频ID（video.id / resource.vid）")
	pgcID := flag.String("pgc_id", "", "PGC ID（snowflake，按字符串传入）")
	// 兼容旧参数（无额外效果）
	_ = flag.Bool("api", false, "兼容旧参数（无额外效果）")
	flag.Parse()

	initialize.InitConfig(*env)
	logger.InitLogger(global.Config)
	global.Mysql = mysql.Init(global.Config.Mysql)

	fmt.Println("OK: 已连接数据库。")

	if *vid == 0 && *pgcID == "" {
		fmt.Fprintln(os.Stderr, "ERR: 请至少传入 -vid 或 -pgc_id")
		os.Exit(2)
	}

	if *vid != 0 {
		runVideoInspect(uint(*vid))
	}
	if *pgcID != "" {
		runPGCInspect(*pgcID)
	}
}

func runVideoInspect(vid uint) {
	fmt.Println("\n== VIDEO ==")
	fmt.Printf("vid=%d\n", vid)

	printQuery("video", `
SELECT id, status, deleted_at, created_at, updated_at
FROM video
WHERE id = ?
LIMIT 1;
`, vid)

	printQuery("resource (by vid)", `
SELECT id, vid, file_id, status, deleted_at, created_at
FROM resource
WHERE vid = ?
ORDER BY id DESC
LIMIT 50;
`, vid)

	printQuery("pgc_episode references (by vid)", `
SELECT id, pgc_id, vid, status, deleted_at, created_at
FROM pgc_episode
WHERE vid = ?
ORDER BY id DESC
LIMIT 50;
`, vid)
}

func runPGCInspect(pgcID string) {
	fmt.Println("\n== PGC ==")
	fmt.Printf("pgc_id=%s\n", pgcID)

	printQuery("pgc_content", `
SELECT id, pgc_id, media_id, status, operator_id, deleted_at, created_at, updated_at
FROM pgc_content
WHERE pgc_id = ?
LIMIT 1;
`, pgcID)

	// media_id 由上面的查询决定，这里用子查询避免额外解析
	printQuery("pgc_media (via pgc_content.media_id)", `
SELECT id, media_id, pgc_type, title, status, deleted_at, created_at, updated_at
FROM pgc_media
WHERE media_id = (SELECT media_id FROM pgc_content WHERE pgc_id = ? LIMIT 1)
LIMIT 1;
`, pgcID)

	printQuery("pgc_episode (by pgc_id)", `
SELECT id, pgc_id, episode_number, vid, status, deleted_at, created_at, updated_at
FROM pgc_episode
WHERE pgc_id = ?
ORDER BY episode_number ASC
LIMIT 200;
`, pgcID)

	printQuery("pgc_review_queue (status IN 100,200)", `
SELECT pgc_id, status, deleted_at, created_at
FROM pgc_content
WHERE status IN (100,200)
ORDER BY id DESC
LIMIT 20;
`)
}

func printQuery(title string, sql string, args ...any) {
	fmt.Printf("\n-- %s --\n", title)
	var rows []Row
	if err := global.Mysql.Raw(sql, args...).Scan(&rows).Error; err != nil {
		fmt.Println("ERR:", err.Error())
		return
	}
	b, _ := json.MarshalIndent(rows, "", "  ")
	fmt.Println(string(b))
}

