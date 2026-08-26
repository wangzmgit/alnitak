package utils

import (
	"bytes"
	"io"
	"mime/multipart"
	"regexp"
	"strings"
)

const MB = 1024 * 1024

// 是否不为空
func VerifyNotEmpty(val interface{}) bool {
	switch v := val.(type) {
	case string:
		return len(v) > 0
	case int:
		return v != 0
	case uint:
		return v != 0
	default:
		return false
	}
}

// 验证字符串长度
func VerifyStringLength(val, op string, length int) bool {
	switch op {
	case "<":
		return len(val) < length
	case ">":
		return len(val) > length
	case "=":
		return len(val) == length
	default:
		return false
	}
}

// 验证邮箱
func VerifyEmail(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// 验证是否为图片
func IsImgType(suffix string, allowedExts []string) bool {
	// 如果配置为空则使用默认值
	if len(allowedExts) == 0 {
		allowedExts = []string{"png", "jpeg", "jpg"}
	}

	// 标准化后缀格式（去除开头的点，转小写）
	suffix = strings.ToLower(strings.TrimPrefix(suffix, "."))

	// 检查后缀是否在允许列表中
	for _, ext := range allowedExts {
		if suffix == strings.ToLower(ext) {
			return true
		}
	}
	return false
}

// 验证是否为视频
func IsVideoType(suffix string, allowedExts []string) bool {
	// 如果配置为空则使用默认值
	if len(allowedExts) == 0 {
		allowedExts = []string{"mp4", "avi", "mkv", "mov", "flv", "wmv", "webm", "m4v", "mpeg", "mpg", "3gp"}
	}

	// 标准化后缀格式（去除开头的点，转小写）
	suffix = strings.ToLower(strings.TrimPrefix(suffix, "."))

	// 检查后缀是否在允许列表中
	for _, ext := range allowedExts {
		if suffix == strings.ToLower(ext) {
			return true
		}
	}
	return false
}

// IsAllASCIIDigits 判断字符串是否非空且每一位均为 ASCII 数字（用于区分数字主键与 shortId）。
func IsAllASCIIDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func FileSize(fileSize int64, count int64, targetSize int64) bool {
	if (fileSize * count) > (targetSize * MB) {
		return false
	}
	return true
}

// ---- 文件魔数校验 ----

// 图片魔数列表
var imageMagicPrefixes = [][]byte{
	{0xFF, 0xD8, 0xFF},                       // JPEG
	{0x89, 0x50, 0x4E, 0x47},                 // PNG
	{0x47, 0x49, 0x46, 0x38},                 // GIF87a / GIF89a
	{0x52, 0x49, 0x46, 0x46},                 // WebP (RIFF)
	{0x42, 0x4D},                              // BMP
}

// 视频魔数列表
var videoMagicPrefixes = [][]byte{
	{0x00, 0x00, 0x00, 0x18, 0x66, 0x74, 0x79, 0x70}, // MP4 (mp42 / isom)
	{0x00, 0x00, 0x00, 0x1C, 0x66, 0x74, 0x79, 0x70}, // MP4 variant
	{0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70}, // MP4 (iso5 / dash)
	{0x1A, 0x45, 0xDF, 0xA3},                            // MKV / WebM
	{0x52, 0x49, 0x46, 0x46},                            // AVI (RIFF)
}

// CheckFileMagicBytes 读取文件头字节并与魔数列表比对
func CheckFileMagicBytes(file *multipart.FileHeader, magicList [][]byte) (bool, error) {
	src, err := file.Open()
	if err != nil {
		return false, err
	}
	defer src.Close()

	// 取最长的魔数长度作为读取字节数
	maxLen := 0
	for _, m := range magicList {
		if len(m) > maxLen {
			maxLen = len(m)
		}
	}
	header := make([]byte, maxLen)
	if _, err := io.ReadFull(src, header); err != nil {
		return false, err
	}

	for _, magic := range magicList {
		if bytes.HasPrefix(header, magic) {
			return true, nil
		}
	}
	return false, nil
}

// CheckImageMagicBytes 校验文件头是否为已知图片格式
func CheckImageMagicBytes(file *multipart.FileHeader) (bool, error) {
	return CheckFileMagicBytes(file, imageMagicPrefixes)
}

// CheckVideoMagicBytes 校验文件头是否为已知视频格式
func CheckVideoMagicBytes(file *multipart.FileHeader) (bool, error) {
	return CheckFileMagicBytes(file, videoMagicPrefixes)
}
