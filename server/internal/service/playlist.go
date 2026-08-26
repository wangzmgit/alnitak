package service

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

const (
	maxPlaylistPerUser  = 50   // 每用户最多创建合集数
	maxVideoPerPlaylist = 200  // 每合集最多视频数
	sortStep            = 1000 // 排序步长
)

// AddPlaylist 创建合集
func AddPlaylist(ctx *gin.Context, req dto.AddPlaylistReq) (uint, error) {
	userId := ctx.GetUint("userId")

	if utils.VerifyStringLength(req.Title, ">", 100) {
		return 0, errors.New("合集标题不能超过100个字符")
	}

	// 检查用户合集数量限制
	var count int64
	global.Mysql.Model(&model.Playlist{}).Where("uid = ?", userId).Count(&count)
	if count >= maxPlaylistPerUser {
		return 0, fmt.Errorf("最多创建%d个合集", maxPlaylistPerUser)
	}

	playlist := model.Playlist{
		Uid:    userId,
		Title:  req.Title,
		Cover:  req.Cover,
		Desc:   req.Desc,
		Status: global.WAITING_REVIEW,
	}
	if err := global.Mysql.Create(&playlist).Error; err != nil {
		utils.ErrorLog("创建合集失败", "playlist", err.Error())
		return 0, errors.New("创建合集失败")
	}

	return playlist.ID, nil
}

// EditPlaylist 编辑合集
func EditPlaylist(ctx *gin.Context, req dto.EditPlaylistReq) error {
	userId := ctx.GetUint("userId")

	var playlist model.Playlist
	if err := global.Mysql.First(&playlist, req.ID).Error; err != nil {
		return errors.New("合集不存在")
	}

	if playlist.Uid != userId {
		return errors.New("无权操作")
	}

	if utils.VerifyStringLength(req.Title, ">", 100) {
		return errors.New("合集标题不能超过100个字符")
	}

	if err := global.Mysql.Model(&model.Playlist{}).Where("id = ?", req.ID).Updates(map[string]interface{}{
		"title": req.Title,
		"cover": req.Cover,
		"desc":  req.Desc,
	}).Error; err != nil {
		utils.ErrorLog("编辑合集失败", "playlist", err.Error())
		return errors.New("编辑失败")
	}

	return nil
}

// DeletePlaylist 删除合集
func DeletePlaylist(ctx *gin.Context, id uint) error {
	userId := ctx.GetUint("userId")

	var playlist model.Playlist
	if err := global.Mysql.First(&playlist, id).Error; err != nil {
		return errors.New("合集不存在")
	}

	if playlist.Uid != userId {
		return errors.New("无权操作")
	}

	// 硬删除关联视频（关联表无需保留历史）
	global.Mysql.Unscoped().Where("playlist_id = ?", id).Delete(&model.PlaylistVideo{})
	// 删除合集
	if err := global.Mysql.Where("id = ?", id).Delete(&model.Playlist{}).Error; err != nil {
		utils.ErrorLog("删除合集失败", "playlist", err.Error())
		return errors.New("删除失败")
	}

	return nil
}

// GetMyPlaylist 获取自己的合集列表
func GetMyPlaylist(ctx *gin.Context) (total int64, list []vo.PlaylistListResp) {
	userId := ctx.GetUint("userId")
	global.Mysql.Model(&model.Playlist{}).Where("uid = ?", userId).Count(&total)
	global.Mysql.Model(&model.Playlist{}).Select(vo.PLAYLIST_LIST_FIELD).
		Where("uid = ?", userId).Order("created_at desc").Scan(&list)
	return
}

// GetPlaylistInfo 获取合集详情
func GetPlaylistInfo(ctx *gin.Context, id uint) (result vo.PlaylistResp, err error) {
	if err = global.Mysql.Model(&model.Playlist{}).Select(vo.PLAYLIST_FIELD).
		Where("id = ?", id).Scan(&result).Error; err != nil {
		return result, errors.New("合集不存在")
	}

	if result.ID == 0 {
		return result, errors.New("合集不存在")
	}

	// 非公开或未审核通过的合集只有创建者可见
	userId := ctx.GetUint("userId")
	if result.Uid != userId && (!result.IsOpen || result.Status != global.AUDIT_APPROVED) {
		return result, errors.New("合集不存在")
	}

	result.Author = GetUserBaseInfo(result.Uid)

	// 增加浏览量（原子操作）
	global.Mysql.Model(&model.Playlist{}).Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + 1"))

	return
}

// GetUserPlaylist 获取用户的公开合集列表
func GetUserPlaylist(ctx *gin.Context, uid uint, page, pageSize int) (total int64, list []vo.PlaylistListResp) {
	global.Mysql.Model(&model.Playlist{}).Where("uid = ? AND is_open = ? AND status = ?", uid, true, global.AUDIT_APPROVED).Count(&total)
	global.Mysql.Model(&model.Playlist{}).Select(vo.PLAYLIST_LIST_FIELD).
		Where("uid = ? AND is_open = ? AND status = ?", uid, true, global.AUDIT_APPROVED).
		Order("created_at desc").
		Limit(pageSize).Offset((page - 1) * pageSize).Scan(&list)
	return
}

// AddPlaylistVideo 添加视频到合集
func AddPlaylistVideo(ctx *gin.Context, req dto.PlaylistVideoAddReq) error {
	userId := ctx.GetUint("userId")

	var playlist model.Playlist
	if err := global.Mysql.First(&playlist, req.PlaylistID).Error; err != nil {
		return errors.New("合集不存在")
	}

	if playlist.Uid != userId {
		return errors.New("无权操作")
	}

	// 检查视频数量限制
	if playlist.VideoCount+len(req.Vids) > maxVideoPerPlaylist {
		return fmt.Errorf("合集视频数量不能超过%d", maxVideoPerPlaylist)
	}

	// 获取当前最大排序值
	var maxSort int64
	global.Mysql.Model(&model.PlaylistVideo{}).Where("playlist_id = ?", req.PlaylistID).
		Select("COALESCE(MAX(sort), 0)").Scan(&maxSort)

	addCount := 0
	var skippedVids []uint
	for _, vid := range req.Vids {
		// 校验视频存在且审核通过
		var video model.Video
		if err := global.Mysql.Where("id = ? AND uid = ? AND status = ?", vid, userId, global.AUDIT_APPROVED).
			First(&video).Error; err != nil {
			continue // 跳过不符合条件的视频
		}

		// 检查视频是否已在用户的其他合集中
		var otherCount int64
		global.Mysql.Model(&model.PlaylistVideo{}).
			Joins("INNER JOIN playlist ON playlist.id = playlist_video.playlist_id AND playlist.deleted_at IS NULL").
			Where("playlist_video.vid = ? AND playlist.uid = ? AND playlist_video.playlist_id != ? AND playlist_video.deleted_at IS NULL",
				vid, userId, req.PlaylistID).
			Count(&otherCount)
		if otherCount > 0 {
			skippedVids = append(skippedVids, vid)
			continue // 已在其他合集中，跳过
		}

		// 检查是否已在合集中（包含软删除的记录）
		var existing model.PlaylistVideo
		err := global.Mysql.Unscoped().
			Where("playlist_id = ? AND vid = ?", req.PlaylistID, vid).First(&existing).Error
		if err == nil {
			// 记录存在
			if existing.DeletedAt.Valid {
				// 软删除的记录，恢复它
				maxSort += sortStep
				global.Mysql.Unscoped().Model(&existing).Updates(map[string]interface{}{
					"deleted_at": nil,
					"sort":       maxSort,
				})
				addCount++
			}
			// 未删除的记录，跳过
			continue
		}

		maxSort += sortStep
		global.Mysql.Create(&model.PlaylistVideo{
			PlaylistID: req.PlaylistID,
			Vid:        vid,
			Sort:       maxSort,
		})
		addCount++
	}

	if len(skippedVids) > 0 && addCount == 0 {
		return errors.New("所选视频已在其他合集中，无法重复添加")
	}

	if addCount > 0 {
		// 更新视频数量（原子操作）
		global.Mysql.Model(&model.Playlist{}).Where("id = ?", req.PlaylistID).
			UpdateColumn("video_count", gorm.Expr("video_count + ?", addCount))
	}

	return nil
}

// DelPlaylistVideo 从合集移除视频
func DelPlaylistVideo(ctx *gin.Context, req dto.PlaylistVideoDelReq) error {
	userId := ctx.GetUint("userId")

	var playlist model.Playlist
	if err := global.Mysql.First(&playlist, req.PlaylistID).Error; err != nil {
		return errors.New("合集不存在")
	}

	if playlist.Uid != userId {
		return errors.New("无权操作")
	}

	result := global.Mysql.Unscoped().Where("playlist_id = ? AND vid IN ?", req.PlaylistID, req.Vids).
		Delete(&model.PlaylistVideo{})

	if result.RowsAffected > 0 {
		// 更新视频数量
		var count int64
		global.Mysql.Model(&model.PlaylistVideo{}).Where("playlist_id = ?", req.PlaylistID).Count(&count)
		global.Mysql.Model(&model.Playlist{}).Where("id = ?", req.PlaylistID).
			Update("video_count", count)
	}

	return nil
}

// SortPlaylistVideo 调整合集视频排序（前端传入完整的视频ID顺序数组）
func SortPlaylistVideo(ctx *gin.Context, req dto.PlaylistVideoSortReq) error {
	userId := ctx.GetUint("userId")

	var playlist model.Playlist
	if err := global.Mysql.First(&playlist, req.PlaylistID).Error; err != nil {
		return errors.New("合集不存在")
	}

	if playlist.Uid != userId {
		return errors.New("无权操作")
	}

	// 按传入的顺序重新分配排序值
	for i, vid := range req.Vids {
		newSort := int64((i + 1) * sortStep)
		global.Mysql.Model(&model.PlaylistVideo{}).
			Where("playlist_id = ? AND vid = ?", req.PlaylistID, vid).
			Update("sort", newSort)
	}

	return nil
}

// GetPlaylistListManage 获取全站合集列表（后台管理）
func GetPlaylistListManage(page, pageSize int) (total int64, list []vo.PlaylistResp) {
	global.Mysql.Model(&model.Playlist{}).Count(&total)
	global.Mysql.Model(&model.Playlist{}).Select(vo.PLAYLIST_FIELD).
		Order("created_at desc").
		Limit(pageSize).Offset((page - 1) * pageSize).Scan(&list)

	for i := 0; i < len(list); i++ {
		list[i].Author = GetUserBaseInfo(list[i].Uid)
	}

	return
}

// DeletePlaylistManage 删除合集（后台管理，无需校验所有者）
func DeletePlaylistManage(id uint) error {
	var playlist model.Playlist
	if err := global.Mysql.First(&playlist, id).Error; err != nil {
		return errors.New("合集不存在")
	}

	// 硬删除关联视频（关联表无需保留历史）
	global.Mysql.Unscoped().Where("playlist_id = ?", id).Delete(&model.PlaylistVideo{})
	// 删除合集
	if err := global.Mysql.Where("id = ?", id).Delete(&model.Playlist{}).Error; err != nil {
		utils.ErrorLog("删除合集失败", "playlist", err.Error())
		return errors.New("删除失败")
	}

	return nil
}

// GetPlaylistVideoList 获取合集视频列表
func GetPlaylistVideoList(ctx *gin.Context, playlistID uint, page, pageSize int) (total int64, list []vo.PlaylistVideoResp) {
	// 检查合集是否存在且可见
	var playlist model.Playlist
	if err := global.Mysql.First(&playlist, playlistID).Error; err != nil {
		return 0, nil
	}

	userId := ctx.GetUint("userId")
	if !playlist.IsOpen && playlist.Uid != userId {
		return 0, nil
	}

	global.Mysql.Model(&model.PlaylistVideo{}).Where("playlist_id = ?", playlistID).Count(&total)

	global.Mysql.Table("playlist_video").
		Select("video.id as vid, video.short_id, video.title, video.cover, video.duration, video.clicks, video.desc, video.created_at").
		Joins("LEFT JOIN video ON playlist_video.vid = video.id").
		Where("playlist_video.playlist_id = ? AND playlist_video.deleted_at IS NULL AND video.deleted_at IS NULL", playlistID).
		Order("playlist_video.sort ASC").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Scan(&list)

	for i := range list {
		if list[i].ShortID == "" {
			list[i].ShortID = utils.UintToString(list[i].Vid)
		}
	}

	return
}

// GetMyPlaylistVideoIds 获取当前用户所有合集中的视频ID（用于前端去重检查）
func GetMyPlaylistVideoIds(ctx *gin.Context) map[uint]uint {
	userId := ctx.GetUint("userId")
	result := make(map[uint]uint)

	var items []struct {
		Vid        uint
		PlaylistID uint
	}
	global.Mysql.Model(&model.PlaylistVideo{}).
		Select("playlist_video.vid, playlist_video.playlist_id").
		Joins("INNER JOIN playlist ON playlist.id = playlist_video.playlist_id AND playlist.deleted_at IS NULL").
		Where("playlist.uid = ? AND playlist_video.deleted_at IS NULL", userId).
		Scan(&items)

	for _, item := range items {
		result[item.Vid] = item.PlaylistID
	}

	return result
}

// GetVideoPlaylists 获取视频所属的公开合集列表
func GetVideoPlaylists(videoId uint) (total int64, list []vo.VideoPlaylistResp) {
	query := global.Mysql.Table("playlist p").
		Joins("INNER JOIN playlist_video pv ON p.id = pv.playlist_id AND pv.deleted_at IS NULL").
		Where("pv.vid = ? AND p.deleted_at IS NULL AND p.is_open = ? AND p.status = ?", videoId, true, global.AUDIT_APPROVED)

	query.Count(&total)
	query.Select("p.id, p.title, p.desc, p.cover, p.uid, p.created_at, p.is_open").
		Order("pv.sort ASC").
		Scan(&list)

	for i := 0; i < len(list); i++ {
		list[i].Author = GetUserBaseInfo(list[i].Uid)
	}

	return
}

// GetVideoPrimaryPlaylist 获取视频的主要合集
func GetVideoPrimaryPlaylist(videoId uint, userId uint) (*vo.VideoPlaylistResp, error) {
	// 获取所有合集
	_, list := GetVideoPlaylists(videoId)

	if len(list) == 0 {
		return nil, errors.New("视频未加入任何合集")
	}

	// 智能选择主合集
	primary := selectPrimaryPlaylist(list, userId)

	return &primary, nil
}

// selectPrimaryPlaylist 智能选择主合集
func selectPrimaryPlaylist(playlists []vo.VideoPlaylistResp, userId uint) vo.VideoPlaylistResp {
	if len(playlists) == 0 {
		return vo.VideoPlaylistResp{}
	}

	// 1. 优先选择用户自己的合集
	for _, playlist := range playlists {
		if playlist.Uid == userId {
			return playlist
		}
	}

	// 2. 选择第一个合集
	return playlists[0]
}

// GetPlaylistVideoListWithParts 获取合集视频列表（包含多分P展开和当前视频分P）
func GetPlaylistVideoListWithParts(ctx *gin.Context, videoId uint) gin.H {
	// 1. 获取视频所属的合集
	_, playlists := GetVideoPlaylists(videoId)

	if len(playlists) == 0 {
		// 没有合集，返回当前视频的分P列表
		resources := GetVideoResources(videoId)
		currentParts := make([]gin.H, 0)
		for i, r := range resources {
			currentParts = append(currentParts, gin.H{
				"p":        i + 1,
				"title":    r.Title,
				"duration": r.Duration,
				"shortId":  r.ShortID,
			})
		}
		return gin.H{
			"playlist":          nil,
			"videos":            []interface{}{},
			"currentVideoParts": currentParts,
			"hasCollection":     false,
		}
	}

	// 2. 获取合集信息
	first := playlists[0]
	playlistInfo := gin.H{
		"id":    first.ID,
		"title": first.Title,
	}

	// 3. 获取合集视频列表
	_, rawVideos := GetPlaylistVideoList(ctx, first.ID, 1, 200)

// 4. 展开多分P视频
	videos := make([]gin.H, 0)
	for _, v := range rawVideos {
		// 获取该视频的资源列表
		resources := GetVideoResources(v.Vid)

		if len(resources) > 1 {
			// 多分P，展开为多项
			for i, r := range resources {
				videos = append(videos, gin.H{
					"vid":         v.Vid,
					"shortId":     v.ShortID,
					"resourceRid": r.ShortID, // 资源的 shortId，用于精确匹配分P
					"title":       v.Title,
					"cover":       v.Cover,
					"duration":    r.Duration,
					"clicks":      v.Clicks,
					"desc":        v.Desc,
					"p":           i + 1,
					"partTitle":   r.Title,
				})
			}
		} else if len(resources) == 1 {
			// 单分P
			videos = append(videos, gin.H{
				"vid":         v.Vid,
				"shortId":     v.ShortID,
				"resourceRid": resources[0].ShortID, // 资源的 shortId
				"title":       v.Title,
				"cover":       v.Cover,
				"duration":    resources[0].Duration,
				"clicks":      v.Clicks,
				"desc":        v.Desc,
				"p":           1,
				"partTitle":   nil,
			})
		} else {
			// 无资源
			videos = append(videos, gin.H{
				"vid":         v.Vid,
				"shortId":     v.ShortID,
				"resourceRid": nil,
				"title":       v.Title,
				"cover":       v.Cover,
				"duration":    0,
				"clicks":      v.Clicks,
				"desc":        v.Desc,
				"p":           nil,
				"partTitle":   nil,
			})
		}
	}

	// 5. 获取当前视频的分P列表
	currentParts := make([]gin.H, 0)
	currentResources := GetVideoResources(videoId)
	for i, r := range currentResources {
		currentParts = append(currentParts, gin.H{
			"p":        i + 1,
			"title":    r.Title,
			"duration": r.Duration,
			"shortId":  r.ShortID,
		})
	}

	return gin.H{
		"playlist":          playlistInfo,
		"videos":            videos,
		"currentVideoParts": currentParts,
		"hasCollection":     true,
	}
}

// GetVideoResources 获取视频的资源列表（分P）
func GetVideoResources(videoId uint) []model.Resource {
	var resources []model.Resource
	global.Mysql.Where("vid = ? AND status = ?", videoId, global.AUDIT_APPROVED).
		Order("sort_order ASC, id ASC").
		Find(&resources)
	return resources
}
