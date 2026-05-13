package service

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func GetDanmaku(ctx *gin.Context, videoId uint, part uint) (danmakus []vo.DanmakuResp, err error) {
	// 通过 video_short_id 查询
	var video model.Video
	if err := global.Mysql.Where("id = ?", videoId).First(&video).Error; err != nil || video.ShortID == "" {
		return danmakus, errors.New("视频不存在")
	}

	if err := global.Mysql.Model(&model.Danmaku{}).Select(vo.DANMAKU_FIELD).
		Where("video_short_id = ? and part = ?", video.ShortID, part).Scan(&danmakus).Error; err != nil {
		utils.ErrorLog("获取弹幕失败", "danmaku", err.Error())
		return danmakus, errors.New("获取失败")
	}

	return danmakus, nil
}

// GetDanmakuByRid 通过 video_short_id + rid 获取弹幕（精准匹配）
// 如果没有匹配，回退到按 part 查询（兼容旧数据）
func GetDanmakuByRid(ctx *gin.Context, videoId uint, rid string, part uint) (danmakus []vo.DanmakuResp, err error) {
	// 先获取 video_short_id
	var video model.Video
	if err := global.Mysql.Where("id = ?", videoId).First(&video).Error; err != nil || video.ShortID == "" {
		utils.ErrorLog("视频不存在或无short_id", "danmaku", fmt.Sprintf("videoId=%d, err=%v", videoId, err))
		return danmakus, errors.New("视频不存在")
	}

	utils.InfoLog(fmt.Sprintf("【弹幕查询】videoId=%d, shortID=%s, rid=%s, part=%d", videoId, video.ShortID, rid, part), "danmaku")

	// 优先通过 video_short_id + rid 精准匹配
	if err := global.Mysql.Model(&model.Danmaku{}).Select(vo.DANMAKU_FIELD).
		Where("video_short_id = ? AND resource_short_id = ?", video.ShortID, rid).Scan(&danmakus).Error; err != nil {
		utils.ErrorLog("通过rid获取弹幕失败", "danmaku", err.Error())
	}
	utils.InfoLog(fmt.Sprintf("【弹幕查询】通过rid查询结果: %d条", len(danmakus)), "danmaku")

	// 如果没有匹配（rid 为空或不存在），回退到按 part 查询
	if len(danmakus) == 0 && part > 0 {
		if err := global.Mysql.Model(&model.Danmaku{}).Select(vo.DANMAKU_FIELD).
			Where("video_short_id = ? AND part = ?", video.ShortID, part).Scan(&danmakus).Error; err != nil {
			utils.ErrorLog("回退到part获取弹幕失败", "danmaku", err.Error())
		}
		utils.InfoLog(fmt.Sprintf("【弹幕查询】回退到part查询结果: %d条", len(danmakus)), "danmaku")
	}
	return danmakus, nil
}

// 获取视频弹幕数量
func GetDanmakuCount(videoId uint) int64 {
	var count int64
	// 通过 video_short_id 查询
	var video model.Video
	if err := global.Mysql.Where("id = ?", videoId).First(&video).Error; err != nil || video.ShortID == "" {
		return count
	}
	global.Mysql.Model(&model.Danmaku{}).Where("video_short_id = ?", video.ShortID).Count(&count)
	return count
}

func SendDanmaku(ctx *gin.Context, danmakuReq dto.DanmakuReq) error {
	videoId, err := ParseVideoID(danmakuReq.Vid)
	if err != nil {
		return errors.New("视频不存在")
	}
	video := GetVideoInfo(videoId)
	if video.ID == 0 {
		return errors.New("视频不存在")
	}

	userId := ctx.GetUint("userId")
	if danmakuReq.Part == 0 {
		danmakuReq.Part = 1
	}

	// 获取 rid 用于精准绑定
	rid := ""
	if danmakuReq.Rid != "" {
		rid = danmakuReq.Rid
	} else {
		rid, _ = GetResourceShortIDByPart(videoId, danmakuReq.Part)
	}

	// 保存到数据库（使用 video_short_id）
	danmaku := model.Danmaku{
		VideoShortID:     video.ShortID,
		Uid:              userId,
		Part:             danmakuReq.Part,
		ResourceShortID:  rid,
		Time:             danmakuReq.Time,
		Type:             danmakuReq.Type,
		Color:            danmakuReq.Color,
		Text:             danmakuReq.Text,
	}
	if err := global.Mysql.Create(&danmaku).Error; err != nil {
		utils.ErrorLog("保存弹幕失败", "danmaku", err.Error())
		return errors.New("发送失败")
	}

	// 广播弹幕消息给同视频的所有用户
	BroadcastDanmaku(videoId, danmaku)

	return nil
}

// BroadcastDanmaku 广播弹幕消息
func BroadcastDanmaku(videoId uint, danmaku model.Danmaku) {
	// 构造弹幕响应
	danmakuResp := vo.DanmakuResp{
		ID:               danmaku.ID,
		Time:             danmaku.Time,
		Type:             danmaku.Type,
		Color:            danmaku.Color,
		Text:             danmaku.Text,
		Part:             danmaku.Part,
		ResourceShortID:  danmaku.ResourceShortID,
		CreatedAt:        danmaku.CreatedAt,
	}

	// 广播给该视频所有分组的观众（包括无rid和有rid的）
	// 无rid分组（只传vid没传rid的用户）
	setMessageAllClient(groupVideoRid(videoId, ""), gin.H{
		"type":    "danmaku",
		"danmaku": danmakuResp,
	})
	// 有rid分组（传了rid的用户）
	if danmaku.ResourceShortID != "" {
		setMessageAllClient(groupVideoRid(videoId, danmaku.ResourceShortID), gin.H{
			"type":    "danmaku",
			"danmaku": danmakuResp,
		})
	}
}

// MigrateAllDanmakuAndHistoryResourceShortID 迁移所有视频的弹幕和历史记录的 resource_short_id
// 启动参数 -migrate-danmaku 时调用
func MigrateAllDanmakuAndHistoryResourceShortID() error {
	// 获取所有视频
	var videos []model.Video
	if err := global.Mysql.Unscoped().Where("deleted_at IS NULL").Find(&videos).Error; err != nil {
		utils.ErrorLog("获取视频列表失败", "migrate", err.Error())
		return err
	}

	updatedDanmaku := int64(0)
	updatedHistory := int64(0)

	for _, video := range videos {
		// 跳过没有 short_id 的视频
		if video.ShortID == "" {
			utils.InfoLog(fmt.Sprintf("【跳过】vid=%d 无short_id", video.ID), "migrate")
			continue
		}

		// 通过 vid 字段匹配，正确设置 video_short_id
		dmResult := global.Mysql.Model(&model.Danmaku{}).
			Where("vid = ? AND (video_short_id IS NULL OR video_short_id = '')", video.ID).
			Updates(map[string]interface{}{
				"video_short_id": video.ShortID,
			})
		if dmResult.Error != nil {
			utils.ErrorLog(fmt.Sprintf("迁移弹幕 video_short_id 失败 vid=%d", video.ID), "migrate", dmResult.Error.Error())
		} else if dmResult.RowsAffected > 0 {
			updatedDanmaku += dmResult.RowsAffected
			utils.InfoLog(fmt.Sprintf("【弹幕迁移】vid=%d -> video_short_id=%s, 更新了%d条", video.ID, video.ShortID, dmResult.RowsAffected), "migrate")
		}

		// 通过 vid 字段匹配，正确设置 video_short_id
		hisResult := global.Mysql.Model(&model.History{}).
			Where("vid = ? AND (video_short_id IS NULL OR video_short_id = '')", video.ID).
			Updates(map[string]interface{}{
				"video_short_id": video.ShortID,
			})
		if hisResult.Error != nil {
			utils.ErrorLog(fmt.Sprintf("迁移历史记录 video_short_id 失败 vid=%d", video.ID), "migrate", hisResult.Error.Error())
		} else if hisResult.RowsAffected > 0 {
			updatedHistory += hisResult.RowsAffected
			utils.InfoLog(fmt.Sprintf("【历史记录迁移】vid=%d -> video_short_id=%s, 更新了%d条", video.ID, video.ShortID, hisResult.RowsAffected), "migrate")
		}
	}

	utils.InfoLog(fmt.Sprintf("【迁移完成】共更新弹幕%d条, 历史记录%d条", updatedDanmaku, updatedHistory), "migrate")
	return nil
}