package utils

import (
	"regexp"
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
func IsImgType(suffix string) bool {
	pattern := `.(png|jpeg|jpg)`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(suffix)
}

// 验证是否为视频
func IsVideoType(suffix string) bool {
	pattern := `.(mp4)`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(suffix)
}

func FileSize(fileSize int64, count int64, targetSize int64) bool {
	if (fileSize * count) > (targetSize * MB) {
		return false
	}
	return true
}
