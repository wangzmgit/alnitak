package utils

import (
	"math/rand"
	"strconv"
	"time"
)

// 生成n位数字随机码
func GenerateNumberCode(length int) string {
	res := ""
	rand.Seed(time.Now().UnixNano())
	// 生成 4 个 [0, 9) 范围的真随机数。
	for i := 0; i < length; i++ {
		num := rand.Intn(10)
		res += strconv.Itoa(num)
	}
	return res
}
