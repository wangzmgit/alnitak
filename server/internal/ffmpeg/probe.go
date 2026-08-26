package ffmpeg

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunFFprobe 执行 ffprobe 并返回结构化结果。
func RunFFprobe(inputPath string) (*ProbeResult, error) {
	args := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		inputPath,
	}
	cmd := exec.Command("ffprobe", args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffprobe %s: %w\nstderr: %s", inputPath, err, stderr.String())
	}

	var result ProbeResult
	if err := json.Unmarshal([]byte(stdout.String()), &result); err != nil {
		return nil, fmt.Errorf("ffprobe json parse %s: %w", inputPath, err)
	}
	return &result, nil
}

// GetMP4InitRange 解析 fMP4 文件的 init range 和 index range。
func GetMP4InitRange(filePath string) (initRange, indexRange string, err error) {
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

	var moovEnd, sidxStart, sidxEnd int64
	offset := int64(0)
	for offset < fileSize {
		boxSize, boxType, err := readMP4Box(file, offset, fileSize)
		if err != nil || boxSize <= 0 {
			break
		}
		switch boxType {
		case "moov":
			moovEnd = offset + boxSize
		case "sidx":
			sidxStart = offset
			sidxEnd = offset + boxSize
		case "moof":
			goto ParseDone
		}
		offset += boxSize
		if sidxEnd > 0 {
			break
		}
	}
ParseDone:
	if moovEnd == 0 {
		return "", "", fmt.Errorf("moov box not found in %s", filePath)
	}
	initRange = fmt.Sprintf("0-%d", moovEnd-1)
	if sidxEnd > 0 {
		indexRange = fmt.Sprintf("%d-%d", sidxStart, sidxEnd-1)
	} else {
		indexRange = fmt.Sprintf("%d-%d", moovEnd, fileSize-1)
	}
	return initRange, indexRange, nil
}

// readMP4Box 读取 MP4 文件的一个 box header，返回 size 和 type。
func readMP4Box(file *os.File, offset, fileSize int64) (int64, string, error) {
	if offset+8 > fileSize {
		return 0, "", fmt.Errorf("offset out of range")
	}
	header := make([]byte, 8)
	if _, err := file.ReadAt(header, offset); err != nil {
		return 0, "", err
	}
	boxSize := int64(binary.BigEndian.Uint32(header[:4]))
	boxType := string(header[4:8])
	if boxSize == 1 {
		// 64-bit extended size
		if offset+16 > fileSize {
			return 0, "", fmt.Errorf("extended size offset out of range")
		}
		extHeader := make([]byte, 8)
		if _, err := file.ReadAt(extHeader, offset+8); err != nil {
			return 0, "", err
		}
		boxSize = int64(binary.BigEndian.Uint64(extHeader))
	}
	return boxSize, boxType, nil
}
