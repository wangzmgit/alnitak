package initialize

import (
	"github.com/bwmarrin/snowflake"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func InitSnowflake() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		utils.ErrorLog("雪花ID初始化失败", "snowflake", err.Error())
	}

	global.SnowflakeNode = node
}
