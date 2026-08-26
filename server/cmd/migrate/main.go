package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/initialize"
	"interastral-peace.com/alnitak/pkg/logger"
	"interastral-peace.com/alnitak/pkg/mysql"
	"interastral-peace.com/alnitak/pkg/redis"
)

func main() {
	fmt.Println("============================================")
	fmt.Println("  VideoIndexFile 元数据迁移工具")
	fmt.Println("============================================")
	fmt.Println("")
	fmt.Println("此工具将旧的 m3u8 Content 字段解析为元数据字段")
	fmt.Println("")

	env := flag.String("env", "prod", "dev/prod")
	flag.Parse()

	initialize.InitConfig(*env)
	logger.InitLogger(global.Config)
	global.Mysql = mysql.Init(global.Config.Mysql)
	global.Redis = redis.Init(global.Config.Redis)
	initialize.InitTables()

	if !confirm() {
		fmt.Println("已取消")
		return
	}

	migrate()
}

func confirm() bool {
	fmt.Println("警告：此操作会更新 video_index_file 表的元数据字段")
	fmt.Println("建议在低峰期执行，并确保已备份数据库")
	fmt.Println("")
	fmt.Print("是否继续执行? (y/N): ")

	var input string
	fmt.Scanln(&input)
	return strings.ToLower(input) == "y"
}

func migrate() {
	// 只查询有 Content 但没有 SegmentCount 的记录（旧数据）
	var indexFiles []model.VideoIndexFile
	global.Mysql.Where("content != '' AND content IS NOT NULL AND segment_count = 0").Find(&indexFiles)

	total := len(indexFiles)
	if total == 0 {
		fmt.Println("无需迁移，所有记录已经包含元数据")
		return
	}

	fmt.Printf("开始迁移，共 %d 条旧格式记录\n\n", total)

	success := 0
	failed := 0

	for i, file := range indexFiles {
		meta := parseM3U8Meta(file.Content)
		width, height, bandwidth, frameRate := parseQualityInfo(file.Quality)

		// 更新元数据字段
		err := global.Mysql.Model(&model.VideoIndexFile{}).Where("id = ?", file.ID).Updates(map[string]interface{}{
			"segment_duration":      meta.SegmentDuration,
			"segment_count":         meta.SegmentCount,
			"total_duration":        meta.TotalDuration,
			"last_segment_duration": meta.LastSegmentDuration,
			"init_file":             meta.InitFile,
			"bandwidth":             bandwidth,
			"width":                 width,
			"height":                height,
			"frame_rate":            frameRate,
		}).Error

		if err != nil {
			fmt.Printf("失败: id=%d resource_id=%d quality=%s error=%s\n",
				file.ID, file.ResourceID, file.Quality, err.Error())
			failed++
		} else {
			success++
		}

		if (i+1)%100 == 0 {
			fmt.Printf("进度: %d/%d (成功: %d, 失败: %d)\n", i+1, total, success, failed)
		}
	}

	fmt.Println("")
	fmt.Println("============================================")
	fmt.Println("  迁移完成")
	fmt.Println("============================================")
	fmt.Printf("  总记录数: %d\n", total)
	fmt.Printf("  成功: %d\n", success)
	fmt.Printf("  失败: %d\n", failed)
	fmt.Println("============================================")
}

// M3U8Meta m3u8 元数据
type M3U8Meta struct {
	SegmentDuration     float64
	SegmentCount        int
	TotalDuration       float64
	LastSegmentDuration float64
	InitFile            string
}

// parseM3U8Meta 解析 m3u8 文件提取元数据
func parseM3U8Meta(content string) M3U8Meta {
	meta := M3U8Meta{}
	lines := strings.Split(content, "\n")

	var durations []float64
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 解析初始化文件 (fMP4)
		if strings.HasPrefix(line, "#EXT-X-MAP:") {
			if idx := strings.Index(line, "URI=\""); idx != -1 {
				start := idx + 5
				end := strings.Index(line[start:], "\"")
				if end != -1 {
					meta.InitFile = line[start : start+end]
				}
			}
		}

		// 解析分片时长
		if strings.HasPrefix(line, "#EXTINF:") {
			durationStr := strings.TrimPrefix(line, "#EXTINF:")
			durationStr = strings.TrimSuffix(durationStr, ",")
			if d, err := strconv.ParseFloat(durationStr, 64); err == nil {
				durations = append(durations, d)
			}
		}
	}

	meta.SegmentCount = len(durations)
	if meta.SegmentCount > 0 {
		for _, d := range durations {
			meta.TotalDuration += d
		}
		meta.SegmentDuration = durations[0]
		meta.LastSegmentDuration = durations[meta.SegmentCount-1]
	}

	return meta
}

// parseQualityInfo 从 quality 字符串解析分辨率、码率、帧率
func parseQualityInfo(quality string) (width, height, bandwidth int, frameRate float64) {
	width, height = 1920, 1080
	bandwidth = 3000000
	frameRate = 30

	parts := strings.Split(quality, "_")
	if len(parts) >= 1 {
		resParts := strings.Split(parts[0], "x")
		if len(resParts) == 2 {
			if w, err := strconv.Atoi(resParts[0]); err == nil {
				width = w
			}
			if h, err := strconv.Atoi(resParts[1]); err == nil {
				height = h
			}
		}
	}

	if len(parts) >= 2 {
		rateStr := strings.TrimSuffix(parts[1], "k")
		if rate, err := strconv.Atoi(rateStr); err == nil {
			bandwidth = rate * 1000
		}
	}

	if len(parts) >= 3 {
		if fps, err := strconv.ParseFloat(parts[2], 64); err == nil {
			frameRate = fps
		}
	}

	return
}
