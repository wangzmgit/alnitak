package initialize

import (
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
)

func InitTables() {
	global.Mysql.AutoMigrate(&model.User{})           // 用户表
	global.Mysql.AutoMigrate(&model.Role{})           // 角色表
	global.Mysql.AutoMigrate(&model.Menu{})           // 菜单表
	global.Mysql.AutoMigrate(&model.Api{})            // Api表
	global.Mysql.AutoMigrate(&model.CasbinRule{})     // casbin规则表
	global.Mysql.AutoMigrate(&model.Operate{})        // 操作日志表
	global.Mysql.AutoMigrate(&model.Partition{})      // 分区表
	global.Mysql.AutoMigrate(&model.Video{})          // 视频表
	global.Mysql.AutoMigrate(&model.Resource{})       // 视频资源表
	global.Mysql.AutoMigrate(&model.VideoIndexFile{}) // 视频播放索引文件表
	global.Mysql.AutoMigrate(&model.Comment{})        // 评论回复表
	global.Mysql.AutoMigrate(&model.Like{})           // 点赞表
	global.Mysql.AutoMigrate(&model.Collect{})        // 收藏表
	global.Mysql.AutoMigrate(&model.Collection{})     // 收藏夹表
	global.Mysql.AutoMigrate(&model.Relation{})       // 关系表
}
