package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/initialize"
	"interastral-peace.com/alnitak/pkg/logger"
	"interastral-peace.com/alnitak/pkg/mysql"
	"interastral-peace.com/alnitak/pkg/redis"
)

// getMP4InitRange 从 MP4 文件中提取初始化范围（ftyp + moov）
// 返回: initRange (0 到 moov 结束), indexRange (moov 开始到 moov 结束)
func getMP4InitRange(filePath string) (initRange, indexRange string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", "", err
	}
	fileSize := fileInfo.Size()

	// 遍历 MP4 box 结构，找到 ftyp 和 moov
	var ftypEnd int64 = 0
	var moovStart int64 = 0
	var moovEnd int64 = 0

	offset := int64(0)
	header := make([]byte, 8)

	for offset < fileSize {
		_, err := file.ReadAt(header, offset)
		if err != nil {
			break
		}

		// 解析 box size（4 字节大端）
		boxSize := int64(uint32(header[0])<<24 | uint32(header[1])<<16 | uint32(header[2])<<8 | uint32(header[3]))
		boxType := string(header[4:8])

		// 处理扩展 size（size == 1 表示使用 64 位 size）
		if boxSize == 1 {
			extHeader := make([]byte, 8)
			file.ReadAt(extHeader, offset+8)
			boxSize = int64(uint64(extHeader[0])<<56 | uint64(extHeader[1])<<48 | uint64(extHeader[2])<<40 | uint64(extHeader[3])<<32 |
				uint64(extHeader[4])<<24 | uint64(extHeader[5])<<16 | uint64(extHeader[6])<<8 | uint64(extHeader[7]))
		}

		// size == 0 表示 box 延伸到文件末尾
		if boxSize == 0 {
			boxSize = fileSize - offset
		}

		switch boxType {
		case "ftyp":
			ftypEnd = offset + boxSize
		case "moov":
			moovStart = offset
			moovEnd = offset + boxSize
		}

		// 找到 moov 就可以停止了
		if moovEnd > 0 {
			break
		}

		offset += boxSize
	}

	if moovEnd == 0 {
		return "", "", fmt.Errorf("moov box not found in %s", filePath)
	}

	// initRange: 从文件开头到 moov 结束（包含 ftyp + moov）
	// indexRange: moov box 的范围
	initRange = fmt.Sprintf("0-%d", moovEnd-1)
	indexRange = fmt.Sprintf("%d-%d", moovStart, moovEnd-1)

	// 如果 ftyp 在 moov 前面，确保 initRange 包含 ftyp
	if ftypEnd > 0 && ftypEnd < moovStart {
		initRange = fmt.Sprintf("0-%d", moovEnd-1)
	}

	return initRange, indexRange, nil
}

func main() {
	fmt.Println("============================================")
	fmt.Println("  视频索引范围修复工具 (SegmentBase)")
	fmt.Println("============================================")
	fmt.Println("")
	fmt.Println("此工具将重新计算 video/audio 的 init_range 和 index_range")
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

	fixRanges()
}

func confirm() bool {
	fmt.Println("警告：此操作会更新 video_index_file 表的 Range 字段")
	fmt.Println("建议在低峰期执行，并确保已备份数据库")
	fmt.Println("")
	fmt.Print("是否继续执行? (y/N): ")

	var input string
	fmt.Scanln(&input)
	return strings.ToLower(input) == "y"
}

func fixRanges() {
	// 查询所有 SegmentBase 模式的记录（有 video_file 字段）
	var indexFiles []model.VideoIndexFile
	global.Mysql.Where("video_file != '' AND video_file IS NOT NULL").Find(&indexFiles)

	total := len(indexFiles)
	if total == 0 {
		fmt.Println("没有需要修复的记录")
		return
	}

	fmt.Printf("开始修复，共 %d 条记录\n\n", total)

	basePath := "./upload/video"
	success := 0
	failed := 0
	skipped := 0

	for i, file := range indexFiles {
		videoPath := filepath.Join(basePath, file.DirName, file.VideoFile)
		audioPath := filepath.Join(basePath, file.DirName, file.AudioFile)

		// 检查视频文件
		if _, err := os.Stat(videoPath); os.IsNotExist(err) {
			fmt.Printf("跳过: id=%d 视频文件不存在: %s\n", file.ID, videoPath)
			skipped++
			continue
		}

		// 获取视频的 Range
		videoInitRange, videoIndexRange, err := getMP4InitRange(videoPath)
		if err != nil {
			fmt.Printf("失败: id=%d 无法解析视频: %s\n", file.ID, err.Error())
			failed++
			continue
		}

		// 获取音频的 Range（如果存在）
		audioInitRange := ""
		audioIndexRange := ""
		if file.AudioFile != "" {
			if _, err := os.Stat(audioPath); err == nil {
				audioInitRange, audioIndexRange, _ = getMP4InitRange(audioPath)
			}
		}

		// 更新数据库
		updates := map[string]interface{}{
			"video_init_range":  videoInitRange,
			"video_index_range": videoIndexRange,
		}
		if audioInitRange != "" {
			updates["audio_init_range"] = audioInitRange
			updates["audio_index_range"] = audioIndexRange
		}

		err = global.Mysql.Model(&model.VideoIndexFile{}).Where("id = ?", file.ID).Updates(updates).Error
		if err != nil {
			fmt.Printf("失败: id=%d 数据库更新错误: %s\n", file.ID, err.Error())
			failed++
			continue
		}

		fmt.Printf("成功: id=%d video=[%s, %s] audio=[%s, %s]\n",
			file.ID, videoInitRange, videoIndexRange, audioInitRange, audioIndexRange)
		success++

		if (i+1)%100 == 0 {
			fmt.Printf("进度: %d/%d (成功: %d, 失败: %d, 跳过: %d)\n", i+1, total, success, failed, skipped)
		}
	}

	fmt.Println("")
	fmt.Println("============================================")
	fmt.Println("  修复完成")
	fmt.Println("============================================")
	fmt.Printf("  总记录数: %d\n", total)
	fmt.Printf("  成功: %d\n", success)
	fmt.Printf("  失败: %d\n", failed)
	fmt.Printf("  跳过: %d\n", skipped)
	fmt.Println("============================================")
}
