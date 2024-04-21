package utils

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
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

// 随机生成图片文件名
func GenerateImgFilename(suffix string) string {
	node, err := snowflake.NewNode(1)
	if err != nil {
		ErrorLog("图片文件名生成失败", "random", err.Error())
		return ""
	}

	id := node.Generate()
	return id.String() + suffix
}

// 随机视频文件名
func GenerateVideoFilename() string {
	node, err := snowflake.NewNode(1)
	if err != nil {
		ErrorLog("视频文件名生成失败", "random", err.Error())
		return ""
	}

	id := node.Generate()
	return id.String()
}
