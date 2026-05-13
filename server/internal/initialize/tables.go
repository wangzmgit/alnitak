package initialize

import (
	"fmt"

	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

func InitTables() {
	global.Mysql.AutoMigrate(&model.User{})           // 用户表
	global.Mysql.AutoMigrate(&model.UserBan{})        // 用户封禁表
	global.Mysql.AutoMigrate(&model.Role{})           // 角色表
	global.Mysql.AutoMigrate(&model.Menu{})           // 菜单表
	global.Mysql.AutoMigrate(&model.Api{})            // Api表
	global.Mysql.AutoMigrate(&model.CasbinRule{})     // casbin规则表
	global.Mysql.AutoMigrate(&model.Operate{})        // 操作日志表
	global.Mysql.AutoMigrate(&model.Partition{})      // 分区表
	global.Mysql.AutoMigrate(&model.Video{})          // 视频表
	global.Mysql.AutoMigrate(&model.VideoFile{})      // 视频文件表（全局去重）
	global.Mysql.AutoMigrate(&model.VideoFileRef{})   // 视频文件引用表
	global.Mysql.AutoMigrate(&model.Tag{})            // 标签
	global.Mysql.AutoMigrate(&model.VideoTag{})       // 视频-标签关联
	global.Mysql.AutoMigrate(&model.Resource{})       // 视频资源表
	global.Mysql.AutoMigrate(&model.VideoIndexFile{}) // 视频播放索引文件表
	global.Mysql.AutoMigrate(&model.Review{})         // 视频审核表
	global.Mysql.AutoMigrate(&model.Comment{})        // 评论回复表
	global.Mysql.AutoMigrate(&model.CommentLike{})   // 评论点赞表
	global.Mysql.AutoMigrate(&model.CommentDislike{}) // 评论点踩表
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
	global.Mysql.AutoMigrate(&model.Whisper{})           // 私信消息表
	global.Mysql.AutoMigrate(&model.MessageReadStatus{}) // 公告/点赞/回复/@ 已读进度表
	global.Mysql.AutoMigrate(&model.Carousel{})          // 轮播图表
	global.Mysql.AutoMigrate(&model.Article{})        // 文章表
	global.Mysql.AutoMigrate(&model.ImageFile{})      // 图片文件表
	global.Mysql.AutoMigrate(&model.Playlist{})       // 合集表
	global.Mysql.AutoMigrate(&model.PlaylistVideo{})  // 合集视频关联表
	global.Mysql.AutoMigrate(&model.PGCMedia{})       // PGC媒体表（media层）
	global.Mysql.AutoMigrate(&model.PGCContent{})     // PGC内容表
	global.Mysql.AutoMigrate(&model.PGCEpisode{})     // PGC剧集表
	global.Mysql.AutoMigrate(&model.SubtitleTrack{})  // 分P字幕轨
	// 认证相关表
	global.Mysql.AutoMigrate(&model.AuthType{})      // 认证类型配置表
	global.Mysql.AutoMigrate(&model.UserAuth{})       // 用户认证记录表
	// 初始化默认认证类型
	service.InitDefaultAuthTypes()

	// 历史数据补齐：为旧的 season 记录回填 media_id
	backfillPGCMedia()

	// 补填已有记录的空 ShortID
	backfillShortIDs()
}

// backfillPGCMedia 为历史 pgc_content（media_id=0）补齐 pgc_media 记录
func backfillPGCMedia() {
	var rows []model.PGCContent
	if err := global.Mysql.Where("media_id = 0").Find(&rows).Error; err != nil {
		utils.ErrorLog("查询待回填PGC media失败", "init", err.Error())
		return
	}
	if len(rows) == 0 {
		return
	}

	utils.InfoLog(fmt.Sprintf("【补齐PGC media】共%d条", len(rows)), "init")
	for _, r := range rows {
		mediaID := uint64(global.SnowflakeNode.Generate())
		media := model.PGCMedia{
			MediaID: mediaID,
			PGCType: r.PGCType,
			Title:   r.Title,
			Cover:   r.Cover,
			Desc:    r.Desc,
			Year:    r.Year,
			Area:    r.Area,
			Rating:  r.Rating,
			Status:  r.Status,
		}
		tx := global.Mysql.Begin()
		if err := tx.Create(&media).Error; err != nil {
			tx.Rollback()
			utils.ErrorLog("创建PGC media回填记录失败", "init", err.Error())
			continue
		}
		if err := tx.Model(&model.PGCContent{}).Where("id = ? AND media_id = 0", r.ID).Update("media_id", mediaID).Error; err != nil {
			tx.Rollback()
			utils.ErrorLog("更新PGC content media_id失败", "init", err.Error())
			continue
		}
		if err := tx.Commit().Error; err != nil {
			utils.ErrorLog("提交PGC media回填事务失败", "init", err.Error())
			continue
		}
	}
}

// backfillShortIDs 为已有的空 short_id 记录补填短ID
func backfillShortIDs() {
	backfillTable := func(tableName string) {
		var records []struct {
			ID uint
		}
		global.Mysql.Table(tableName).Where("short_id = '' OR short_id IS NULL").Select("id").Find(&records)
		if len(records) == 0 {
			return
		}
		utils.InfoLog(fmt.Sprintf("【补填ShortID】%s 共%d条空记录", tableName, len(records)), "init")
		for _, r := range records {
			var shortID string
			switch tableName {
			case "video":
				sid, err := service.AllocateUniqueVideoShortID()
				if err != nil {
					utils.ErrorLog("补填ShortID失败 video", "init", err.Error())
					continue
				}
				shortID = sid
			case "resource":
				sid, err := service.AllocateUniqueResourceShortID()
				if err != nil {
					utils.ErrorLog("补填ShortID失败 resource", "init", err.Error())
					continue
				}
				shortID = sid
			case "article":
				sid, err := service.AllocateUniqueArticleShortID()
				if err != nil {
					utils.ErrorLog("补填ShortID失败 article", "init", err.Error())
					continue
				}
				shortID = sid
			default:
				continue
			}
			global.Mysql.Table(tableName).Where("id = ?", r.ID).Update("short_id", shortID)
		}
		utils.InfoLog(fmt.Sprintf("【补填ShortID】%s 完成", tableName), "init")
	}

	backfillTable("video")
	backfillTable("resource")
	backfillTable("article")
}
