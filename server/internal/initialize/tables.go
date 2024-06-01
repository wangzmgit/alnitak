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
	global.Mysql.AutoMigrate(&model.Review{})         // 视频审核表
	global.Mysql.AutoMigrate(&model.Comment{})        // 评论回复表
	global.Mysql.AutoMigrate(&model.LikeVideo{})      // 视频点赞表
	global.Mysql.AutoMigrate(&model.LikeArticle{})    // 文章点赞表
	global.Mysql.AutoMigrate(&model.CollectVideo{})   // 视频收藏表
	global.Mysql.AutoMigrate(&model.CollectArticle{}) // 文章收藏表
	global.Mysql.AutoMigrate(&model.Collection{})     // 收藏夹表
	global.Mysql.AutoMigrate(&model.Relation{})       // 关系表
	global.Mysql.AutoMigrate(&model.Danmaku{})        // 弹幕表
	global.Mysql.AutoMigrate(&model.History{})        // 历史记录表
	global.Mysql.AutoMigrate(&model.Announce{})       // 公告表
	global.Mysql.AutoMigrate(&model.LikeMessage{})    // 点赞消息表
	global.Mysql.AutoMigrate(&model.AtMessage{})      // @消息表
	global.Mysql.AutoMigrate(&model.ReplyMessage{})   // 回复消息表
	global.Mysql.AutoMigrate(&model.Whisper{})        // 私信消息表
	global.Mysql.AutoMigrate(&model.Carousel{})       // 轮播图表
	global.Mysql.AutoMigrate(&model.Article{})        // 文章表
}
