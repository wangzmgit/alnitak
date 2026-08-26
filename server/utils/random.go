package utils

import (
	"crypto/rand"
	"math/big"
	"strconv"
)

// 生成n位数字随机码
func GenerateNumberCode(length int) string {
	res := ""
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			panic("crypto/rand failed: " + err.Error())
		}
		res += strconv.Itoa(int(num.Int64()))
	}
	return res
}
