package utils

import "os"

// 判断文件是否存在
func IsFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
