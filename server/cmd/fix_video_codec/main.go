package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/initialize"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/pkg/logger"
	"interastral-peace.com/alnitak/pkg/mysql"
	"interastral-peace.com/alnitak/pkg/oss"
)

// 历史 video_codec 对齐工具（不重编码）：
// - 仅处理 SegmentBase 行（有 video_file + video_init_range）
// - 对 *_video.m4s 等本地或 OSS 文件跑 ffprobe，生成 avc1 字符串
// - 仅当探测成功且 新值 != 当前库中值 时更新（或 dry-run 只打印）
//
// 请在 server 目录下执行（与 ./upload/video 相对路径一致），例如：
//
//	go run ./cmd/fix_video_codec -env=prod -dry-run
//	go run ./cmd/fix_video_codec -env=prod -yes
func main() {
	env := flag.String("env", "prod", "dev/prod")
	dryRun := flag.Bool("dry-run", false, "只打印将变更的行，不写库")
	yes := flag.Bool("yes", false, "跳过交互确认（非 dry-run 时）")
	uploadRoot := flag.String("upload", "./upload/video", "本地视频根目录")
	limit := flag.Int("limit", 0, "最多处理条数，0 表示不限制")
	useOss := flag.Bool("use-oss", false, "本地不存在时，才从 OSS 下载探测（默认仅用本地）")
	flag.Parse()

	fmt.Println("============================================")
	fmt.Println("  video_index_file.video_codec 对齐（ffprobe）")
	fmt.Println("============================================")
	fmt.Println("")

	initialize.InitConfig(*env)
	logger.InitLogger(global.Config)

	global.Mysql = mysql.Init(global.Config.Mysql)
	initialize.InitTables()

	if !*dryRun {
		fmt.Println("警告：将按 ffprobe 结果更新 video_codec（仅与当前值不一致时）")
		fmt.Println("建议先 -dry-run，并确保已备份数据库")
		fmt.Println("")
		if !*yes && !confirm() {
			fmt.Println("已取消")
			return
		}
	} else {
		fmt.Println("当前为 -dry-run：不会修改数据库")
		fmt.Println("")
	}

	if err := run(*uploadRoot, *dryRun, *limit, *useOss); err != nil {
		zap.L().Error("fix_video_codec 失败", zap.Error(err))
		os.Exit(1)
	}
}

func initOssSafely() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("oss 初始化 panic: %v", r)
		}
	}()
	global.Storage = oss.InitStorage(global.Config.Storage)
	return nil
}

func confirm() bool {
	fmt.Print("是否继续? (y/N): ")
	var input string
	_, _ = fmt.Scanln(&input)
	return strings.ToLower(strings.TrimSpace(input)) == "y"
}

func run(uploadRoot string, dryRun bool, limit int, useOss bool) error {
	var indexFiles []model.VideoIndexFile
	q := global.Mysql.Where("video_file != '' AND video_file IS NOT NULL AND video_init_range != '' AND video_init_range IS NOT NULL")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&indexFiles).Error; err != nil {
		return fmt.Errorf("查询 video_index_file: %w", err)
	}

	total := len(indexFiles)
	if total == 0 {
		fmt.Println("没有 SegmentBase 记录需要检查")
		return nil
	}

	fmt.Printf("待检查: %d 条\n\n", total)

	var updated, skippedSame, skippedNoFile, skippedProbe, failed int

	for i, row := range indexFiles {
		localPath := filepath.Join(uploadRoot, row.DirName, row.VideoFile)
		probePath := localPath
		tmpFile := ""
		if _, err := os.Stat(localPath); err != nil {
			if !useOss {
				fmt.Printf("跳过 id=%d resource=%d: 本地不存在且未启用 -use-oss: %s\n", row.ID, row.ResourceID, localPath)
				skippedNoFile++
				continue
			}

			if global.Storage == nil {
				if err := initOssSafely(); err != nil {
					fmt.Printf("跳过 id=%d resource=%d: OSS 初始化失败: %v\n", row.ID, row.ResourceID, err)
					skippedNoFile++
					continue
				}
			}

			f, err := os.CreateTemp("", "probe-video-*.m4s")
			if err != nil {
				fmt.Printf("失败 id=%d: 创建临时文件: %v\n", row.ID, err)
				failed++
				continue
			}
			tmpFile = f.Name()
			_ = f.Close()
			key := "video/" + row.DirName + "/" + row.VideoFile
			if err := global.Storage.GetObjectToFile(key, tmpFile); err != nil {
				_ = os.Remove(tmpFile)
				fmt.Printf("跳过 id=%d resource=%d: OSS 拉取失败 %s: %v\n", row.ID, row.ResourceID, key, err)
				skippedNoFile++
				continue
			}
			probePath = tmpFile
		}

		newCodec, err := service.ProbeH264Avc1CodecString(probePath)
		if tmpFile != "" {
			_ = os.Remove(tmpFile)
		}
		if err != nil {
			fmt.Printf("跳过 id=%d resource=%d quality=%s: ffprobe %v\n", row.ID, row.ResourceID, row.Quality, err)
			skippedProbe++
			continue
		}

		oldCodec := strings.TrimSpace(row.VideoCodec)
		if newCodec == oldCodec {
			skippedSame++
			if (i+1)%500 == 0 {
				fmt.Printf("进度 %d/%d (已更新 %d)\n", i+1, total, updated)
			}
			continue
		}

		fmt.Printf("变更 id=%d resource=%d quality=%s dir=%s | %q -> %q\n",
			row.ID, row.ResourceID, row.Quality, row.DirName, oldCodec, newCodec)

		if dryRun {
			updated++
			continue
		}

		if err := global.Mysql.Model(&model.VideoIndexFile{}).Where("id = ?", row.ID).
			Update("video_codec", newCodec).Error; err != nil {
			fmt.Printf("失败 id=%d 写库: %v\n", row.ID, err)
			failed++
			continue
		}
		updated++

		if (updated)%50 == 0 {
			fmt.Printf("进度: 已更新 %d 条（扫描 %d/%d）\n", updated, i+1, total)
		}
	}

	fmt.Println("")
	fmt.Println("============================================")
	fmt.Println("  完成")
	fmt.Println("============================================")
	fmt.Printf("  扫描: %d\n", total)
	if dryRun {
		fmt.Printf("  将更新（dry-run）: %d\n", updated)
	} else {
		fmt.Printf("  已更新: %d\n", updated)
	}
	fmt.Printf("  已一致跳过: %d\n", skippedSame)
	fmt.Printf("  无文件/OSS 失败跳过: %d\n", skippedNoFile)
	fmt.Printf("  ffprobe 失败跳过: %d\n", skippedProbe)
	fmt.Printf("  其它失败: %d\n", failed)
	fmt.Println("============================================")

	return nil
}
