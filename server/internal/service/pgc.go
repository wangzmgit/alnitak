package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func CreatePGCContent(req dto.CreatePGCReq) (uint64, error) {
	if req.Title == "" {
		return 0, errors.New("标题不能为空")
	}

	if len(req.Title) > 255 {
		return 0, errors.New("标题长度不能超过255字符")
	}

	if req.Cover == "" {
		return 0, errors.New("封面不能为空")
	}

	if req.PGCType <= 0 {
		return 0, errors.New("PGC类型不能为空")
	}

	validTypes := map[int]bool{
		global.PGCTypeCN:          true,
		global.PGCTypeJP:          true,
		global.PGCTypeDocumentary: true,
		global.PGCTypeMovie:       true,
		global.PGCTypeTVSeries:    true,
	}

	if !validTypes[req.PGCType] {
		return 0, errors.New("无效的PGC类型")
	}

	if req.Rating < 0 || req.Rating > 10 {
		return 0, errors.New("评分必须在0-10之间")
	}

	if req.PlannedTotalEpisodes != nil {
		if *req.PlannedTotalEpisodes < 0 {
			return 0, errors.New("计划总集数不能为负数")
		}
		if *req.PlannedTotalEpisodes > 999 {
			return 0, errors.New("计划总集数不能超过999")
		}
	}

	if len(req.Episodes) > 0 {
		positiveVids := uniquePositiveVidsFromEpisodeReq(req.Episodes)
		if len(positiveVids) > 0 {
			var count int64
			if err := global.Mysql.Model(&model.Video{}).
				Where("id IN ?", positiveVids).
				Count(&count).Error; err != nil {
				utils.ErrorLog("校验视频存在失败", "pgc", err.Error())
				return 0, errors.New("校验视频失败")
			}

			if count != int64(len(positiveVids)) {
				return 0, errors.New("部分视频不存在")
			}
		}

		episodeMap := make(map[int]bool)
		for _, episode := range req.Episodes {
			if episodeMap[episode.EpisodeNumber] {
				return 0, errors.New("剧集序号重复")
			}
			episodeMap[episode.EpisodeNumber] = true
		}
	}

	pgcID := uint64(global.SnowflakeNode.Generate())
	mediaID := uint64(global.SnowflakeNode.Generate())

	var totalEp, currEp int
	if len(req.Episodes) > 0 {
		totalEp = len(req.Episodes)
		currEp = countEpisodeReqWithBoundVideo(req.Episodes)
	} else {
		currEp = 0
		if req.PlannedTotalEpisodes != nil {
			totalEp = *req.PlannedTotalEpisodes
		}
	}

	tx := global.Mysql.Begin()
	committed := false
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			utils.ErrorLog("创建PGC内容panic", "pgc", "")
			return
		}
		if !committed {
			tx.Rollback()
		}
	}()

	pgcMedia := &model.PGCMedia{
		MediaID: mediaID,
		PGCType: req.PGCType,
		Title:   req.Title,
		Cover:   req.Cover,
		Desc:    req.Desc,
		Year:    req.Year,
		Area:    req.Area,
		Rating:  req.Rating,
		Status:  global.PGCAuditSubmitted,
	}
	if err := tx.Create(pgcMedia).Error; err != nil {
		utils.ErrorLog("创建PGC媒体失败", "pgc", err.Error())
		return 0, errors.New("创建PGC媒体失败")
	}

	pgcContent := &model.PGCContent{
		PGCID:           pgcID,
		MediaID:         mediaID,
		PGCType:         req.PGCType,
		Title:           req.Title,
		Cover:           req.Cover,
		Desc:            req.Desc,
		Year:            req.Year,
		Area:            req.Area,
		Rating:          req.Rating,
		IsOngoing:       req.IsOngoing,
		TotalEpisodes:   totalEp,
		CurrentEpisodes: currEp,
		Status:          global.PGCAuditSubmitted,
	}

	if err := tx.Create(pgcContent).Error; err != nil {
		utils.ErrorLog("创建PGC内容失败", "pgc", err.Error())
		return 0, errors.New("创建PGC内容失败")
	}

	for _, episode := range req.Episodes {
		publishTime := episode.PublishTime
		if episode.VID > 0 && publishTime == "" {
			publishTime = getVideoPublishTime(episode.VID)
		}
		pgcEpisode := &model.PGCEpisode{
			PGCID:         pgcID,
			EpisodeNumber: episode.EpisodeNumber,
			Title:         episode.Title,
			VID:           episode.VID,
			Duration:      episode.Duration,
			Status:        global.PGCEpisodeNormal,
			PublishTime:   publishTime,
		}

		if err := tx.Create(pgcEpisode).Error; err != nil {
			utils.ErrorLog("创建PGC剧集失败", "pgc", err.Error())
			return 0, errors.New("创建剧集失败")
		}
	}

	// 审核链路隔离：被绑定为 PGC 剧集的视频不再走 UGC 的「视频审核」队列
	if len(req.Episodes) > 0 {
		if err := markVideosAsPGCAttached(tx, uniquePositiveVidsFromEpisodeReq(req.Episodes)); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.ErrorLog("提交事务失败", "pgc", err.Error())
		return 0, errors.New("创建失败")
	}

	committed = true
	utils.InfoLog("创建PGC内容成功", "pgc")

	return pgcID, nil
}

func UpdatePGCContent(req dto.UpdatePGCReq) error {
	if req.PGCID == 0 {
		return errors.New("PGC ID不能为空")
	}

	if req.Title != "" && len(req.Title) > 255 {
		return errors.New("标题长度不能超过255字符")
	}

	if req.Rating != nil && (*req.Rating < 0 || *req.Rating > 10) {
		return errors.New("评分必须在0-10之间")
	}
	if req.TotalEpisodes != nil && *req.TotalEpisodes <= 0 {
		return errors.New("总集数必须大于0")
	}

	var pgcContent model.PGCContent
	if err := global.Mysql.Where("pgc_id = ?", req.PGCID).First(&pgcContent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorLog("PGC内容不存在", "pgc", "")
			return errors.New("PGC内容不存在")
		}
		utils.ErrorLog("查询PGC内容失败", "pgc", err.Error())
		return errors.New("查询失败")
	}

	updates := make(map[string]interface{})

	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Cover != "" {
		updates["cover"] = req.Cover
	}
	if req.Desc != "" {
		updates["desc"] = req.Desc
	}
	if req.Year > 0 {
		updates["year"] = req.Year
	}
	if req.Area != "" {
		updates["area"] = req.Area
	}
	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}
	if req.IsOngoing != nil {
		updates["is_ongoing"] = *req.IsOngoing
	}
	if req.TotalEpisodes != nil {
		if *req.TotalEpisodes < pgcContent.CurrentEpisodes {
			return errors.New("总集数不能小于当前已更新集数")
		}
		updates["total_episodes"] = *req.TotalEpisodes
	}

	if len(updates) == 0 {
		return nil
	}

	if err := global.Mysql.Model(&model.PGCContent{}).
		Where("pgc_id = ?", req.PGCID).
		Updates(updates).Error; err != nil {
		utils.ErrorLog("更新PGC内容失败", "pgc", err.Error())
		return errors.New("更新失败")
	}

	// 同步更新 media 层（对齐 media/season 语义）
	var mediaID uint64
	_ = global.Mysql.Model(&model.PGCContent{}).
		Where("pgc_id = ?", req.PGCID).
		Select("media_id").
		Scan(&mediaID).Error
	if mediaID > 0 {
		mediaUpdates := make(map[string]interface{})
		if v, ok := updates["title"]; ok {
			mediaUpdates["title"] = v
		}
		if v, ok := updates["cover"]; ok {
			mediaUpdates["cover"] = v
		}
		if v, ok := updates["desc"]; ok {
			mediaUpdates["desc"] = v
		}
		if v, ok := updates["year"]; ok {
			mediaUpdates["year"] = v
		}
		if v, ok := updates["area"]; ok {
			mediaUpdates["area"] = v
		}
		if v, ok := updates["rating"]; ok {
			mediaUpdates["rating"] = v
		}
		if len(mediaUpdates) > 0 {
			if err := global.Mysql.Model(&model.PGCMedia{}).
				Where("media_id = ?", mediaID).
				Updates(mediaUpdates).Error; err != nil {
				utils.ErrorLog("同步更新PGC媒体失败", "pgc", err.Error())
			}
		}
	}

	utils.InfoLog("更新PGC内容成功", "pgc")

	return nil
}

func UpdatePGCStatus(pgcID uint64, status int) error {
	if pgcID == 0 {
		return errors.New("PGC ID不能为空")
	}
	if status != global.PGCAuditOffline && status != global.PGCAuditApproved {
		return errors.New("无效的状态值")
	}

	var pgcContent model.PGCContent
	if err := global.Mysql.Where("pgc_id = ?", pgcID).First(&pgcContent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("PGC内容不存在")
		}
		return errors.New("查询失败")
	}

	tx := global.Mysql.Begin()
	committed := false
	var touchedVids []uint
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&model.PGCContent{}).Where("pgc_id = ?", pgcID).Update("status", status).Error; err != nil {
		return errors.New("更新失败")
	}

	// 同步 media 状态
	if pgcContent.MediaID > 0 {
		if err := tx.Model(&model.PGCMedia{}).Where("media_id = ?", pgcContent.MediaID).Update("status", status).Error; err != nil {
			return errors.New("更新失败")
		}
	}

	// 置为通过时，联动通过剧集/视频/资源
	if status == global.PGCAuditApproved {
		var err error
		touchedVids, err = approvePGCLinkedData(tx, pgcID)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("更新失败")
	}
	committed = true

	// 清理视频缓存，避免状态延迟
	for _, vid := range touchedVids {
		cache.DelVideoInfo(vid)
	}
	return nil
}

// GetPGCReviewList 后台待审列表：已提交、审核中
func GetPGCReviewList(page, pageSize int) (total int64, list []model.PGCContent, err error) {
	query := global.Mysql.Model(&model.PGCContent{}).
		Where("status IN ?", []int{global.PGCAuditSubmitted, global.PGCAuditProcessing})

	if err := query.Count(&total).Error; err != nil {
		utils.ErrorLog("查询PGC待审总数失败", "pgc", err.Error())
		return 0, nil, errors.New("查询失败")
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&list).Error; err != nil {
		utils.ErrorLog("查询PGC待审列表失败", "pgc", err.Error())
		return 0, nil, errors.New("查询失败")
	}
	return
}

// AdminReviewPGC 后台审核通过或驳回（仅待审状态可流转）
func AdminReviewPGC(pgcID uint64, targetStatus int) error {
	if pgcID == 0 {
		return errors.New("PGC ID不能为空")
	}
	if targetStatus != global.PGCAuditApproved && targetStatus != global.PGCAuditRejected {
		return errors.New("无效的状态值")
	}

	var pgcContent model.PGCContent
	if err := global.Mysql.Where("pgc_id = ?", pgcID).First(&pgcContent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("PGC内容不存在")
		}
		return errors.New("查询失败")
	}

	if pgcContent.Status != global.PGCAuditSubmitted && pgcContent.Status != global.PGCAuditProcessing {
		return errors.New("当前状态不可审核")
	}

	tx := global.Mysql.Begin()
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&model.PGCContent{}).Where("pgc_id = ?", pgcID).Update("status", targetStatus).Error; err != nil {
		return errors.New("更新失败")
	}

	if pgcContent.MediaID > 0 {
		_ = tx.Model(&model.PGCMedia{}).Where("media_id = ?", pgcContent.MediaID).Update("status", targetStatus).Error
	}

	// 审核通过时，联动通过剧集/视频/资源
	if targetStatus == global.PGCAuditApproved {
		vids, err := approvePGCLinkedData(tx, pgcID)
		if err != nil {
			return err
		}
		// 提交后再清缓存
		defer func(vs []uint) {
			if committed {
				for _, vid := range vs {
					cache.DelVideoInfo(vid)
				}
			}
		}(vids)
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("更新失败")
	}
	committed = true
	return nil
}

// approvePGCLinkedData 将 PGC 关联的剧集、视频、资源一并置为“通过可播”。
func approvePGCLinkedData(tx *gorm.DB, pgcID uint64) ([]uint, error) {
	// 剧集状态统一置为 normal（通过）
	if err := tx.Model(&model.PGCEpisode{}).
		Where("pgc_id = ?", pgcID).
		Update("status", global.PGCEpisodeNormal).Error; err != nil {
		utils.ErrorLog("更新PGC剧集状态失败", "pgc", err.Error())
		return nil, errors.New("更新失败")
	}

	var vids []uint
	if err := tx.Model(&model.PGCEpisode{}).Where("pgc_id = ?", pgcID).Pluck("vid", &vids).Error; err != nil {
		utils.ErrorLog("查询PGC关联视频失败", "pgc", err.Error())
		return nil, errors.New("更新失败")
	}
	positive := make([]uint, 0, len(vids))
	seen := make(map[uint]struct{}, len(vids))
	for _, id := range vids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		positive = append(positive, id)
	}
	if len(positive) == 0 {
		return []uint{}, nil
	}

	if err := tx.Model(&model.Video{}).
		Where("id IN ?", positive).
		Where("status NOT IN ?", []int{global.VIDEO_PROCESSING, global.PROCESSING_FAIL, global.CREATED_VIDEO}).
		Update("status", global.AUDIT_APPROVED).Error; err != nil {
		utils.ErrorLog("更新PGC关联视频状态失败", "pgc", err.Error())
		return nil, errors.New("更新失败")
	}

	if err := tx.Model(&model.Resource{}).
		Where("vid IN ?", positive).
		Where("status IN ?", []int{global.WAITING_REVIEW, global.SUBMIT_REVIEW}).
		Update("status", global.AUDIT_APPROVED).Error; err != nil {
		utils.ErrorLog("更新PGC关联资源状态失败", "pgc", err.Error())
		return nil, errors.New("更新失败")
	}
	return positive, nil
}

func UpdatePGCEpisodeStatus(pgcID uint64, episodeID uint64, status int) error {
	if pgcID == 0 || episodeID == 0 {
		return errors.New("参数不能为空")
	}
	if status != global.PGCEpisodeOffline && status != global.PGCEpisodeNormal {
		return errors.New("无效的状态值")
	}
	var ep model.PGCEpisode
	if err := global.Mysql.Where("id = ? AND pgc_id = ?", episodeID, pgcID).First(&ep).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("剧集不存在")
		}
		return errors.New("查询失败")
	}
	if err := global.Mysql.Model(&model.PGCEpisode{}).Where("id = ? AND pgc_id = ?", episodeID, pgcID).Update("status", status).Error; err != nil {
		return errors.New("更新失败")
	}
	var currEp int64
	if err := global.Mysql.Model(&model.PGCEpisode{}).
		Where("pgc_id = ? AND status = ? AND vid > ?", pgcID, global.PGCEpisodeNormal, 0).
		Count(&currEp).Error; err != nil {
		return errors.New("更新集数失败")
	}
	if err := global.Mysql.Model(&model.PGCContent{}).
		Where("pgc_id = ?", pgcID).
		UpdateColumn("current_episodes", int(currEp)).Error; err != nil {
		return errors.New("更新集数失败")
	}
	return nil
}

func DeletePGCContent(pgcID uint64) error {
	if pgcID == 0 {
		return errors.New("PGC ID不能为空")
	}

	var pgcContent model.PGCContent
	if err := global.Mysql.Where("pgc_id = ?", pgcID).First(&pgcContent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("PGC内容不存在")
		}
		return errors.New("查询失败")
	}

	tx := global.Mysql.Begin()
	committed := false
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			utils.ErrorLog("删除PGC内容panic", "pgc", "")
			return
		}
		if !committed {
			tx.Rollback()
		}
	}()

	// 收集关联的视频ID和UID，用于删除后清理
	type vidInfo struct {
		VID uint
		UID uint
	}
	var episodeVids []uint
	tx.Model(&model.PGCEpisode{}).Where("pgc_id = ?", pgcID).Pluck("vid", &episodeVids)

	// 查询视频的 UID（deleteVideoAndRelatedData 需要）
	vidInfos := make([]vidInfo, 0, len(episodeVids))
	for _, vid := range episodeVids {
		var v model.Video
		if err := tx.Select("id, uid").Where("id = ?", vid).First(&v).Error; err == nil {
			vidInfos = append(vidInfos, vidInfo{VID: v.ID, UID: v.Uid})
		}
	}

	if err := tx.Unscoped().Delete(&model.PGCEpisode{}, "pgc_id = ?", pgcID).Error; err != nil {
		utils.ErrorLog("删除PGC剧集失败", "pgc", err.Error())
		return errors.New("删除剧集失败")
	}

	// 兼容 1:1 阶段：删除 season 时同步删除 media
	if pgcContent.MediaID > 0 {
		if err := tx.Unscoped().Delete(&model.PGCMedia{}, "media_id = ?", pgcContent.MediaID).Error; err != nil {
			utils.ErrorLog("删除PGC媒体失败", "pgc", err.Error())
			return errors.New("删除媒体失败")
		}
	}

	if err := tx.Unscoped().Delete(&model.PGCContent{}, "pgc_id = ?", pgcID).Error; err != nil {
		utils.ErrorLog("删除PGC内容失败", "pgc", err.Error())
		return errors.New("删除失败")
	}

	if err := tx.Commit().Error; err != nil {
		utils.ErrorLog("提交事务失败", "pgc", err.Error())
		return errors.New("删除失败")
	}

	committed = true

	// 事务提交后，删除关联的视频及其所有数据（文件、资源、弹幕等）
	for _, vi := range vidInfos {
		if err := deleteVideoAndRelatedData(vi.VID, vi.UID, nil); err != nil {
			utils.ErrorLog("删除PGC关联视频失败", "pgc", fmt.Sprintf("vid=%d err=%s", vi.VID, err.Error()))
			// 不中断，继续删除其他视频
		}
	}

	utils.InfoLog("删除PGC内容成功", "pgc")

	return nil
}

func GetPGCContentList(req dto.PGCListReq) (total int64, list []model.PGCContent, err error) {
	query := global.Mysql.Model(&model.PGCContent{})

	if req.PGCType > 0 {
		query = query.Where("pgc_type = ?", req.PGCType)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		// `desc` 是 SQL 关键字，查询时必须转义列名
		query = query.Where("title LIKE ? OR `desc` LIKE ?", keyword, keyword)
	}
	if req.Year > 0 {
		query = query.Where("year = ?", req.Year)
	}
	if req.Area != "" {
		query = query.Where("area = ?", req.Area)
	}
	if req.IsOngoing != nil {
		query = query.Where("is_ongoing = ?", *req.IsOngoing)
	}

	if err := query.Count(&total).Error; err != nil {
		utils.ErrorLog("查询PGC内容总数失败", "pgc", err.Error())
		return 0, nil, errors.New("查询失败")
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").
		Limit(req.PageSize).
		Offset(offset).
		Find(&list).Error; err != nil {
		utils.ErrorLog("查询PGC内容列表失败", "pgc", err.Error())
		return 0, nil, errors.New("查询失败")
	}

	return
}

func GetPGCContentDetail(pgcID uint64) (*model.PGCContent, error) {
	if pgcID == 0 {
		return nil, errors.New("PGC ID不能为空")
	}

	var pgcContent model.PGCContent
	if err := global.Mysql.Where("pgc_id = ?", pgcID).First(&pgcContent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorLog("PGC内容不存在", "pgc", "")
			return nil, errors.New("PGC内容不存在")
		}
		utils.ErrorLog("查询PGC内容详情失败", "pgc", err.Error())
		return nil, errors.New("查询失败")
	}

	return &pgcContent, nil
}

func GetPGCEpisodeList(pgcID uint64, req dto.EpisodeListReq) ([]model.PGCEpisode, error) {
	if pgcID == 0 {
		utils.ErrorLog("获取PGC剧集列表参数错误", "pgc", "pgcID不能为空")
		return nil, errors.New("PGC ID不能为空")
	}

	var count int64
	if err := global.Mysql.Model(&model.PGCContent{}).
		Where("pgc_id = ?", pgcID).
		Count(&count).Error; err != nil {
		utils.ErrorLog("查询PGC内容失败", "pgc", err.Error())
		return nil, errors.New("查询失败")
	}

	if count == 0 {
		utils.ErrorLog("PGC内容不存在", "pgc", "")
		return nil, errors.New("PGC内容不存在")
	}

	var episodes []model.PGCEpisode
	query := global.Mysql.Model(&model.PGCEpisode{}).Where("pgc_id = ?", pgcID)

	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("episode_number ASC").
		Limit(req.PageSize).
		Offset(offset).
		Find(&episodes).Error; err != nil {
		utils.ErrorLog("查询PGC剧集列表失败", "pgc", err.Error())
		return nil, errors.New("查询失败")
	}

	return episodes, nil
}

func AddPGCEpisode(pgcID uint64, episode dto.EpisodeReq) error {
	if pgcID == 0 {
		utils.ErrorLog("添加PGC剧集参数错误", "pgc", "pgcID不能为空")
		return errors.New("PGC ID不能为空")
	}

	var pgcContent model.PGCContent
	if err := global.Mysql.Where("pgc_id = ?", pgcID).First(&pgcContent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorLog("PGC内容不存在", "pgc", "")
			return errors.New("PGC内容不存在")
		}
		utils.ErrorLog("查询PGC内容失败", "pgc", err.Error())
		return errors.New("查询失败")
	}

	if episode.VID > 0 {
		var count int64
		if err := global.Mysql.Model(&model.Video{}).
			Where("id = ?", episode.VID).
			Count(&count).Error; err != nil {
			utils.ErrorLog("查询视频失败", "pgc", err.Error())
			return errors.New("查询视频失败")
		}

		if count == 0 {
			return errors.New("视频不存在")
		}
	}

	publishTime := episode.PublishTime
	if episode.VID > 0 && publishTime == "" {
		publishTime = getVideoPublishTime(episode.VID)
	}

	pgcEpisode := &model.PGCEpisode{
		PGCID:         pgcID,
		EpisodeNumber: episode.EpisodeNumber,
		Title:         episode.Title,
		VID:           episode.VID,
		Duration:      episode.Duration,
		Status:        global.PGCEpisodeNormal,
		PublishTime:   publishTime,
	}

	tx := global.Mysql.Begin()
	committed := false
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			utils.ErrorLog("添加PGC剧集panic", "pgc", "")
			return
		}
		if !committed {
			tx.Rollback()
		}
	}()

	if err := tx.Create(pgcEpisode).Error; err != nil {
		// 幂等处理：并发/重试导致同一 (pgc_id, episode_number) 重复插入时，
		// 若已存在记录绑定的是同一 VID，则按成功处理，避免前端收到 500。
		if strings.Contains(err.Error(), "Duplicate entry") {
			var existed model.PGCEpisode
			e := tx.Where("pgc_id = ? AND episode_number = ?", pgcID, episode.EpisodeNumber).First(&existed).Error
			if e == nil && existed.VID == episode.VID {
				utils.InfoLog("PGC剧集重复提交，按幂等成功处理", "pgc")
			} else {
				utils.ErrorLog("创建PGC剧集失败", "pgc", err.Error())
				return errors.New("剧集序号已存在")
			}
		} else {
			utils.ErrorLog("创建PGC剧集失败", "pgc", err.Error())
			return errors.New("创建失败")
		}
	}

	var currEp int64
	if err := tx.Model(&model.PGCEpisode{}).
		Where("pgc_id = ? AND status = ? AND vid > ?", pgcID, global.PGCEpisodeNormal, 0).
		Count(&currEp).Error; err != nil {
		utils.ErrorLog("统计PGC当前集数失败", "pgc", err.Error())
		return errors.New("更新集数失败")
	}

	if err := tx.Model(&model.PGCContent{}).
		Where("pgc_id = ?", pgcID).
		UpdateColumn("current_episodes", int(currEp)).Error; err != nil {
		utils.ErrorLog("更新PGC内容集数失败", "pgc", err.Error())
		return errors.New("更新集数失败")
	}

	if episode.VID > 0 {
		if err := markVideosAsPGCAttached(tx, []uint{episode.VID}); err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.ErrorLog("提交事务失败", "pgc", err.Error())
		return errors.New("创建失败")
	}
	committed = true

	utils.InfoLog("添加PGC剧集成功", "pgc")

	return nil
}

// BindPGCEpisodeVideo 将占位剧集绑定到已存在的视频（vid>0）。
func BindPGCEpisodeVideo(pgcID uint64, episodeID uint64, req dto.BindPGCEpisodeVideoReq) error {
	if pgcID == 0 || episodeID == 0 {
		return errors.New("参数不能为空")
	}
	if req.VID == 0 {
		return errors.New("视频ID不能为空")
	}
	var ep model.PGCEpisode
	if err := global.Mysql.Where("id = ? AND pgc_id = ?", episodeID, pgcID).First(&ep).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("剧集不存在")
		}
		return errors.New("查询失败")
	}
	if ep.VID == req.VID {
		return nil
	}
	if ep.VID > 0 {
		return errors.New("该剧集已绑定视频，请先解绑或不使用该接口")
	}

	var count int64
	if err := global.Mysql.Model(&model.Video{}).Where("id = ?", req.VID).Count(&count).Error; err != nil {
		return errors.New("查询视频失败")
	}
	if count == 0 {
		return errors.New("视频不存在")
	}

	publishTime := req.PublishTime
	if publishTime == "" {
		publishTime = getVideoPublishTime(req.VID)
	}
	duration := ep.Duration
	if req.Duration != nil {
		duration = *req.Duration
	}

	tx := global.Mysql.Begin()
	committed := false
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			utils.ErrorLog("绑定PGC剧集视频panic", "pgc", "")
			return
		}
		if !committed {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&model.PGCEpisode{}).
		Where("id = ? AND pgc_id = ?", episodeID, pgcID).
		Updates(map[string]interface{}{
			"vid":          req.VID,
			"duration":     duration,
			"publish_time": publishTime,
			"status":       global.PGCEpisodeNormal,
		}).Error; err != nil {
		utils.ErrorLog("绑定PGC剧集视频失败", "pgc", err.Error())
		return errors.New("绑定失败")
	}

	var currEp int64
	if err := tx.Model(&model.PGCEpisode{}).
		Where("pgc_id = ? AND status = ? AND vid > ?", pgcID, global.PGCEpisodeNormal, 0).
		Count(&currEp).Error; err != nil {
		utils.ErrorLog("统计PGC当前集数失败", "pgc", err.Error())
		return errors.New("更新集数失败")
	}
	if err := tx.Model(&model.PGCContent{}).
		Where("pgc_id = ?", pgcID).
		UpdateColumn("current_episodes", int(currEp)).Error; err != nil {
		utils.ErrorLog("更新PGC内容集数失败", "pgc", err.Error())
		return errors.New("更新集数失败")
	}

	if err := markVideosAsPGCAttached(tx, []uint{req.VID}); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		utils.ErrorLog("提交事务失败", "pgc", err.Error())
		return errors.New("绑定失败")
	}
	committed = true
	return nil
}

func UpdatePGCEpisode(pgcID uint64, episodeID uint64, req dto.UpdatePGCEpisodeReq) error {
	if pgcID == 0 || episodeID == 0 {
		return errors.New("参数不能为空")
	}
	var ep model.PGCEpisode
	if err := global.Mysql.Where("id = ? AND pgc_id = ?", episodeID, pgcID).First(&ep).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("剧集不存在")
		}
		return errors.New("查询失败")
	}
	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if len(updates) == 0 {
		return nil
	}
	if err := global.Mysql.Model(&model.PGCEpisode{}).
		Where("id = ? AND pgc_id = ?", episodeID, pgcID).
		Updates(updates).Error; err != nil {
		return errors.New("更新失败")
	}
	return nil
}

func getVideoPublishTime(vid uint) string {
	var v model.Video
	if err := global.Mysql.Select("id,created_at").Where("id = ?", vid).First(&v).Error; err != nil {
		return ""
	}
	if v.CreatedAt.IsZero() {
		return ""
	}
	return v.CreatedAt.In(time.Local).Format("2006-01-02 15:04:05")
}

// markVideosAsPGCAttached 标记视频已绑定 PGC，并将 UGC 待审状态迁出。
// 仅迁移 WAITING_REVIEW 到 SUBMIT_REVIEW，避免误标记为已通过（AUDIT_APPROVED）。
func markVideosAsPGCAttached(db *gorm.DB, vids []uint) error {
	if len(vids) == 0 {
		return nil
	}
	set := make(map[uint]struct{}, len(vids))
	uniq := make([]uint, 0, len(vids))
	for _, vid := range vids {
		if vid == 0 {
			continue
		}
		if _, ok := set[vid]; ok {
			continue
		}
		set[vid] = struct{}{}
		uniq = append(uniq, vid)
	}
	if len(uniq) == 0 {
		return nil
	}

	if err := db.Model(&model.Video{}).
		Where("id IN ?", uniq).
		Update("pgc_attached", true).Error; err != nil {
		utils.ErrorLog("标记PGC绑定视频失败", "pgc", err.Error())
		return errors.New("更新关联视频状态失败")
	}

	if err := db.Model(&model.Video{}).
		Where("id IN ?", uniq).
		Where("status = ?", global.WAITING_REVIEW).
		Update("status", global.SUBMIT_REVIEW).Error; err != nil {
		utils.ErrorLog("迁移PGC绑定视频审核状态失败", "pgc", err.Error())
		return errors.New("更新关联视频状态失败")
	}

	// 立刻从 UGC 分区/热门缓存移除，避免首页等列表继续出现已绑定 PGC 的视频
	type videoCacheMeta struct {
		ID          uint `gorm:"column:id"`
		PartitionID uint `gorm:"column:partition_id"`
	}
	var metas []videoCacheMeta
	if err := db.Model(&model.Video{}).Select("id, partition_id").Where("id IN ?", uniq).Scan(&metas).Error; err == nil {
		for _, m := range metas {
			cache.DelVideoInfo(m.ID)
			cache.DelVideoId(global.VideoPartitionMap[m.PartitionID], m.ID)
			cache.DelSingleHotVideoId(m.ID)
		}
	}
	return nil
}

func DeletePGCEpisode(pgcID uint64, episodeID uint64) error {
	if pgcID == 0 || episodeID == 0 {
		return errors.New("参数不能为空")
	}

	var episode model.PGCEpisode
	if err := global.Mysql.Where("id = ? AND pgc_id = ?", episodeID, pgcID).First(&episode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("剧集不存在")
		}
		return errors.New("查询失败")
	}

	// 查询关联视频的 UID，用于后续清理
	var video model.Video
	var videoUID uint
	if err := global.Mysql.Select("id, uid").Where("id = ?", episode.VID).First(&video).Error; err == nil {
		videoUID = video.Uid
	}

	tx := global.Mysql.Begin()
	committed := false
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			utils.ErrorLog("删除PGC剧集panic", "pgc", "")
			return
		}
		if !committed {
			tx.Rollback()
		}
	}()

	if err := tx.Unscoped().Delete(&model.PGCEpisode{}, "id = ? AND pgc_id = ?", episodeID, pgcID).Error; err != nil {
		utils.ErrorLog("删除PGC剧集失败", "pgc", err.Error())
		return errors.New("删除失败")
	}

	var currEp int64
	if err := tx.Model(&model.PGCEpisode{}).
		Where("pgc_id = ? AND status = ? AND vid > ?", pgcID, global.PGCEpisodeNormal, 0).
		Count(&currEp).Error; err != nil {
		utils.ErrorLog("统计PGC当前集数失败", "pgc", err.Error())
		return errors.New("更新集数失败")
	}

	if err := tx.Model(&model.PGCContent{}).
		Where("pgc_id = ?", pgcID).
		UpdateColumn("current_episodes", int(currEp)).Error; err != nil {
		utils.ErrorLog("更新PGC内容集数失败", "pgc", err.Error())
		return errors.New("更新集数失败")
	}

	if err := tx.Commit().Error; err != nil {
		utils.ErrorLog("提交事务失败", "pgc", err.Error())
		return errors.New("删除失败")
	}
	committed = true

	// 事务提交后，删除关联的视频及其所有数据
	if episode.VID != 0 {
		if err := deleteVideoAndRelatedData(episode.VID, videoUID, nil); err != nil {
			utils.ErrorLog("删除PGC关联视频失败", "pgc", fmt.Sprintf("vid=%d err=%s", episode.VID, err.Error()))
		}
	}

	utils.InfoLog("删除PGC剧集成功", "pgc")

	return nil
}

func SearchPGC(keyword string, pgcType int, page, pageSize int) (total int64, list []model.PGCContent, err error) {
	query := global.Mysql.Model(&model.PGCContent{})

	if keyword != "" {
		likeKeyword := "%" + keyword + "%"
		// `desc` 是 SQL 关键字，查询时必须转义列名
		query = query.Where("title LIKE ? OR `desc` LIKE ?", likeKeyword, likeKeyword)
	}

	if pgcType > 0 {
		query = query.Where("pgc_type = ?", pgcType)
	}

	query = query.Where("status = ?", global.PGCAuditApproved)

	if err := query.Count(&total).Error; err != nil {
		utils.ErrorLog("搜索PGC内容失败", "pgc", err.Error())
		return 0, nil, errors.New("搜索失败")
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&list).Error; err != nil {
		utils.ErrorLog("搜索PGC内容失败", "pgc", err.Error())
		return 0, nil, errors.New("搜索失败")
	}

	return
}

func GetPGCByType(pgcType int, page, pageSize int) (total int64, list []model.PGCContent, err error) {
	if pgcType <= 0 {
		utils.ErrorLog("获取PGC类型列表参数错误", "pgc", "pgcType不能为空")
		return 0, nil, errors.New("PGC类型不能为空")
	}

	query := global.Mysql.Model(&model.PGCContent{}).
		Where("pgc_type = ? AND status = ?", pgcType, global.PGCAuditApproved)

	if err := query.Count(&total).Error; err != nil {
		utils.ErrorLog("获取PGC类型总数失败", "pgc", err.Error())
		return 0, nil, errors.New("查询失败")
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&list).Error; err != nil {
		utils.ErrorLog("获取PGC类型列表失败", "pgc", err.Error())
		return 0, nil, errors.New("查询失败")
	}

	return
}

func GetOngoingPGC(page, pageSize int) (total int64, list []model.PGCContent, err error) {
	query := global.Mysql.Model(&model.PGCContent{}).
		Where("is_ongoing = ? AND status = ?", true, global.PGCAuditApproved)

	if err := query.Count(&total).Error; err != nil {
		utils.ErrorLog("获取连载PGC失败", "pgc", err.Error())
		return 0, nil, errors.New("查询失败")
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&list).Error; err != nil {
		utils.ErrorLog("获取连载PGC失败", "pgc", err.Error())
		return 0, nil, errors.New("查询失败")
	}

	return
}

func GetRecommendedPGC(limit int) ([]model.PGCContent, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	var list []model.PGCContent
	if err := global.Mysql.Model(&model.PGCContent{}).
		Where("status = ?", global.PGCAuditApproved).
		Order("created_at DESC").
		Limit(limit).
		Find(&list).Error; err != nil {
		utils.ErrorLog("获取推荐PGC失败", "pgc", err.Error())
		return nil, errors.New("查询失败")
	}

	return list, nil
}

func GetPGCWithEpisodes(pgcID uint64) (*model.PGCContent, []model.PGCEpisode, error) {
	if pgcID == 0 {
		return nil, nil, errors.New("PGC ID不能为空")
	}

	var pgcContent model.PGCContent
	// 管理端/编辑页需要查看任意状态的 PGC，不做状态过滤
	if err := global.Mysql.Where("pgc_id = ?", pgcID).First(&pgcContent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorLog("PGC内容不存在", "pgc", "")
			return nil, nil, errors.New("PGC内容不存在")
		}
		utils.ErrorLog("查询PGC内容失败", "pgc", err.Error())
		return nil, nil, errors.New("查询失败")
	}

	var episodes []model.PGCEpisode
	// 同上：不做剧集状态过滤，直接返回该 PGC 下的所有剧集
	if err := global.Mysql.Model(&model.PGCEpisode{}).
		Where("pgc_id = ?", pgcID).
		Order("episode_number ASC").
		Find(&episodes).Error; err != nil {
		utils.ErrorLog("查询PGC剧集失败", "pgc", err.Error())
		return &pgcContent, nil, errors.New("查询失败")
	}

	return &pgcContent, episodes, nil
}

func GetPGCEpisodeDetail(epID uint64) (*model.PGCEpisode, error) {
	if epID == 0 {
		return nil, errors.New("剧集ID不能为空")
	}
	var ep model.PGCEpisode
	if err := global.Mysql.Where("id = ?", epID).First(&ep).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("剧集不存在")
		}
		return nil, errors.New("查询失败")
	}
	return &ep, nil
}

// GetPGCPlayPanelByVideo 按当前播放 vid 返回 PGC 播放面板数据：
// - current: 当前所在 season
// - seasons: 同 media 下所有可用 season（按年份/创建时间排序）
// - episodes: 选中 season 的剧集（按集号升序）
func GetPGCPlayPanelByVideo(vid uint, seasonID uint64) (*model.PGCContent, []model.PGCContent, []model.PGCEpisode, uint64, error) {
	if vid == 0 {
		return nil, nil, nil, 0, errors.New("视频ID不能为空")
	}

	var bindEp model.PGCEpisode
	// 这里按 vid 识别绑定关系，不限制剧集状态，避免因状态流转导致“已绑定却识别失败”。
	if err := global.Mysql.Where("vid = ?", vid).
		Order("id DESC").First(&bindEp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, 0, errors.New("当前视频未绑定PGC")
		}
		utils.ErrorLog("PGC面板查询绑定剧集失败", "pgc", err.Error())
		return nil, nil, nil, 0, errors.New("查询失败")
	}

	var current model.PGCContent
	if err := global.Mysql.Where("pgc_id = ?", bindEp.PGCID).First(&current).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, 0, errors.New("PGC内容不存在")
		}
		utils.ErrorLog("PGC面板查询内容失败", "pgc", err.Error())
		return nil, nil, nil, 0, errors.New("查询失败")
	}

	activeSeasonID := current.PGCID
	if seasonID > 0 {
		activeSeasonID = seasonID
	}

	var seasons []model.PGCContent
	if current.MediaID > 0 {
		if err := global.Mysql.
			Where("media_id = ? AND status = ?", current.MediaID, global.PGCAuditApproved).
			Order("year ASC, created_at ASC").
			Find(&seasons).Error; err != nil {
			utils.ErrorLog("PGC面板查询季列表失败", "pgc", err.Error())
			return nil, nil, nil, 0, errors.New("查询失败")
		}
	} else {
		seasons = []model.PGCContent{current}
	}
	if len(seasons) == 0 {
		seasons = []model.PGCContent{current}
	}

	var exists bool
	for _, s := range seasons {
		if s.PGCID == activeSeasonID {
			exists = true
			break
		}
	}
	if !exists {
		activeSeasonID = current.PGCID
	}

	var episodes []model.PGCEpisode
	// 播放面板仅展示已绑定视频且可参与播放的剧集（vid=0 的占位集不出现在用户播放列表）
	if err := global.Mysql.
		Where("pgc_id = ? AND vid > ? AND status = ?", activeSeasonID, 0, global.PGCEpisodeNormal).
		Order("episode_number ASC, id ASC").
		Find(&episodes).Error; err != nil {
		utils.ErrorLog("PGC面板查询剧集列表失败", "pgc", err.Error())
		return nil, nil, nil, 0, errors.New("查询失败")
	}
	return &current, seasons, episodes, activeSeasonID, nil
}
