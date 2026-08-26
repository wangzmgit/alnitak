package service

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gorm.io/gorm"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

// 并发安全控制：防止 cron 与手动 API 同时执行清理
var (
	cleanupMu      sync.Mutex
	cleanupRunning int32
)

// CleanupItem 待清理项目
type CleanupItem struct {
	Type   string `json:"type"`   // video, image, subtitle
	Path   string `json:"path"`   // 本地路径或文件名
	Reason string `json:"reason"` // 清理原因
}

// CleanupResult 清理结果
type CleanupResult struct {
	CleanedVideoDirs  int           `json:"cleanedVideoDirs"`
	CleanedImages     int           `json:"cleanedImages"`
	CleanedSubtitles  int           `json:"cleanedSubtitles"`
	CleanedVideoFiles int           `json:"cleanedVideoFiles"`
	CleanedIndexFiles int           `json:"cleanedIndexFiles"`
	CleanedImageFiles int           `json:"cleanedImageFiles"`
	CleanedResources  int           `json:"cleanedResources"`
	Errors            []string      `json:"errors"`
	DryRun            bool          `json:"dryRun"`
	Items             []CleanupItem `json:"items"` // 待清理/已清理的文件列表
}

// CleanupOrphanedResources 清理孤立资源
// dryRun: 如果为true，只返回将要清理的内容，不实际执行删除
// 并发安全：dryRun 可并行，非 dryRun 互斥（cron + API 不会同时跑）
func CleanupOrphanedResources(dryRun bool) CleanupResult {
	if !dryRun {
		if !atomic.CompareAndSwapInt32(&cleanupRunning, 0, 1) {
			return CleanupResult{Errors: []string{"清理任务正在进行中"}}
		}
		defer atomic.StoreInt32(&cleanupRunning, 0)

		cleanupMu.Lock()
		defer cleanupMu.Unlock()
	}

	result := CleanupResult{
		Errors: make([]string, 0),
		Items:  make([]CleanupItem, 0),
		DryRun: dryRun,
	}

	// 1. 清理孤立的视频目录（基于完整关联校验）
	result.cleanOrphanedVideoDirs(dryRun)

	// 2. 清理孤立的图片文件
	result.cleanOrphanedImages(dryRun)

	// 3. 清理孤立的字幕文件
	result.cleanOrphanedSubtitles(dryRun)

	return result
}

// cleanOrphanedVideoDirs 清理孤立的视频目录
// 清理条件：本地视频目录对应的资源在数据库中不存在有效引用
// 校验链路：Video -> Resource -> VideoIndexFile(dirName) / VideoFile(dirName)
func (r *CleanupResult) cleanOrphanedVideoDirs(dryRun bool) {
	videoDir := "./upload/video"
	if !utils.IsFileExists(videoDir) {
		return
	}

	entries, err := os.ReadDir(videoDir)
	if err != nil {
		r.Errors = append(r.Errors, "读取视频目录失败: "+err.Error())
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirName := entry.Name()

		// 跳过 .gitkeep 等隐藏文件/目录
		if strings.HasPrefix(dirName, ".") {
			continue
		}

		// 检查这个目录是否有有效的关联
		reason := checkVideoDirValidity(dirName)
		if reason == "" {
			continue // 有效，保留
		}

		// 需要清理
		localDir := filepath.Join(videoDir, dirName)
		r.Items = append(r.Items, CleanupItem{
			Type:   "video",
			Path:   dirName,
			Reason: reason,
		})

		if !dryRun {
			// 删除OSS上的文件
			deleteVideoFromOSS(localDir, &r.Errors)
			// 删除本地目录
			if err := os.RemoveAll(localDir); err != nil {
				r.Errors = append(r.Errors, "删除视频目录失败: "+localDir+" - "+err.Error())
			} else {
				r.CleanedVideoDirs++
				utils.InfoLog("已删除视频目录: "+localDir+" 原因: "+reason, "cleanup")
			}

			// 清理相关的数据库记录
			cleanVideoDirDbRecords(dirName, r)
		} else {
			r.CleanedVideoDirs++
		}
	}
}

// checkVideoDirValidity 检查视频目录是否有效
// 返回空字符串表示有效，返回原因字符串表示需要清理
// 注意：同一个dirName可能有多条记录（不同画质），只要有一条是有效的，整个目录就应该保留
// 全局去重模式下，还需要检查 VideoFile 的引用计数和 VideoFileRef 表
func checkVideoDirValidity(dirName string) string {
	// 1. 检查 VideoIndexFile 表（获取所有记录）
	var indexFiles []model.VideoIndexFile
	global.Mysql.Unscoped().Where("dir_name = ?", dirName).Find(&indexFiles)

	if len(indexFiles) == 0 {
		// VideoIndexFile 不存在，检查 VideoFile
		var videoFile model.VideoFile
		global.Mysql.Unscoped().Where("dir_name = ?", dirName).First(&videoFile)

		if videoFile.ID == 0 {
			return "数据库无记录"
		}

		// VideoFile 存在，检查是否有有效引用（全局去重模式）
		// 检查 VideoFileRef 是否有未删除的引用
		var refCount int64
		global.Mysql.Model(&model.VideoFileRef{}).Where("file_id = ?", videoFile.ID).Count(&refCount)
		if refCount > 0 {
			// 有引用，保留
			return ""
		}

		// 检查 VideoFile 自身的引用计数
		if videoFile.RefCount > 0 {
			return ""
		}

		// 检查是否正在上传/转码中（状态 0=上传中, 1=已合并, 2=转码中）
		if videoFile.Status < model.FileStatusReady {
			// 还在处理中，保留
			return ""
		}

		// 检查是否有正在转码的资源引用了这个文件（重新转码场景）
		var processingCount int64
		global.Mysql.Model(&model.Resource{}).
			Where("status = ?", global.VIDEO_PROCESSING).
			Where("id IN (?)",
				global.Mysql.Model(&model.VideoFileRef{}).Select("resource_id").
					Where("file_id = ?", videoFile.ID),
			).Count(&processingCount)
		if processingCount > 0 {
			return ""
		}

		// VideoFile 存在但无引用且已完成，可以清理
		return "VideoFile无有效引用"
	}

	// 2. 遍历所有 VideoIndexFile 记录，只要有一条完整有效的链路，就保留目录
	var lastReason string
	for _, indexFile := range indexFiles {
		// 检查这条记录是否有效
		reason := checkSingleIndexFileValidity(indexFile)
		if reason == "" {
			// 找到一条有效记录，目录有效，保留
			return ""
		}
		lastReason = reason
	}

	// 3. 检查是否有正在转码的资源使用了该目录（防止重新转码时被误删）
	// 重新转码会先删旧 VideoIndexFile、软删旧 Resource，再创建新 Resource 开始转码
	// 此时新资源尚未生成 VideoIndexFile，但目录不能删
	var processingCount int64
	global.Mysql.Model(&model.Resource{}).
		Where("status = ?", global.VIDEO_PROCESSING).
		Where("id IN (?)",
			global.Mysql.Model(&model.VideoFileRef{}).Select("resource_id").
				Where("file_id IN (?)",
					global.Mysql.Model(&model.VideoFile{}).Select("id").Where("dir_name = ?", dirName),
				),
		).Count(&processingCount)
	if processingCount > 0 {
		return ""
	}

	// 所有记录都无效，返回最后一个无效原因
	return lastReason
}

// checkSingleIndexFileValidity 检查单条 VideoIndexFile 记录是否有效
func checkSingleIndexFileValidity(indexFile model.VideoIndexFile) string {
	// 检查 VideoIndexFile 是否被软删除
	if indexFile.DeletedAt.Valid {
		return "VideoIndexFile已删除"
	}

	// 检查对应的 Resource
	var resource model.Resource
	global.Mysql.Unscoped().Where("id = ?", indexFile.ResourceID).First(&resource)

	if resource.ID == 0 {
		return "Resource记录不存在"
	}

	if resource.DeletedAt.Valid {
		return "Resource已删除"
	}

	// 检查对应的 Video
	var video model.Video
	global.Mysql.Unscoped().Where("id = ?", resource.Vid).First(&video)

	if video.ID == 0 {
		return "Video记录不存在"
	}

	if video.DeletedAt.Valid {
		return "Video已删除"
	}

	// 这条记录完整有效
	return ""
}

// cleanVideoDirDbRecords 清理视频目录相关的数据库记录
// 所有 DB 写入操作在一个事务内完成，避免崩溃导致数据不一致
func cleanVideoDirDbRecords(dirName string, r *CleanupResult) {
	// 先查询 VideoFile，获取其 ID 用于清理引用表（只读，在事务外）
	var videoFile model.VideoFile
	global.Mysql.Unscoped().Where("dir_name = ?", dirName).First(&videoFile)

	// 先收集 ResourceID（必须在删除 VideoIndexFile 之前，只读，在事务外）
	var indexFiles []model.VideoIndexFile
	global.Mysql.Unscoped().Where("dir_name = ?", dirName).Find(&indexFiles)

	// 在事务内执行所有 DELETE
	var (
		cleanedIndexFiles int64
		cleanedVideoFiles int64
		cleanedResources  int64
	)
	err := global.Mysql.Transaction(func(tx *gorm.DB) error {
		// 事务内重新校验：防止检查与删除之间有新转码完成创建了有效记录
		var freshVideoFile model.VideoFile
		tx.Unscoped().Where("dir_name = ?", dirName).First(&freshVideoFile)
		if freshVideoFile.ID != 0 {
			// 有 VideoFile → 检查是否有正在转码的资源引用
			var processingCount int64
			tx.Model(&model.Resource{}).
				Where("status = ?", global.VIDEO_PROCESSING).
				Where("file_id = ?", freshVideoFile.ID).
				Count(&processingCount)
			if processingCount > 0 {
				return nil // 有转码中的资源，跳过本次清理
			}
		}
		var freshIndexCount int64
		tx.Model(&model.VideoIndexFile{}).Where("dir_name = ?", dirName).Count(&freshIndexCount)
		if freshIndexCount > 0 {
			// 有新创建的 VideoIndexFile → 逐条检查是否有效
			var freshIndexFiles []model.VideoIndexFile
			tx.Unscoped().Where("dir_name = ?", dirName).Find(&freshIndexFiles)
			for _, fif := range freshIndexFiles {
				if !fif.DeletedAt.Valid {
					var res model.Resource
					if tx.Unscoped().Where("id = ?", fif.ResourceID).First(&res); res.ID != 0 && !res.DeletedAt.Valid {
						return nil // 有有效资源，跳过
					}
				}
			}
		}

		// 删除 VideoIndexFile 记录
		result := tx.Unscoped().Where("dir_name = ?", dirName).Delete(&model.VideoIndexFile{})
		cleanedIndexFiles = result.RowsAffected

		// 清理 VideoFileRef 引用记录（全局去重模式）
		if freshVideoFile.ID != 0 {
			tx.Unscoped().Where("file_id = ?", freshVideoFile.ID).Delete(&model.VideoFileRef{})
		}

		// 删除 VideoFile 记录
		result = tx.Unscoped().Where("dir_name = ?", dirName).Delete(&model.VideoFile{})
		cleanedVideoFiles = result.RowsAffected

		// 删除相关的 Resource 记录（包括通过 indexFiles 和通过 file_id 关联的）
		resourceIDs := make(map[uint]bool)
		for _, indexFile := range indexFiles {
			if indexFile.ResourceID > 0 {
				resourceIDs[indexFile.ResourceID] = true
			}
		}
		// 补充：通过 file_id 关联的 Resource（无 VideoIndexFile 的场景）
		if freshVideoFile.ID != 0 {
			var extraResources []model.Resource
			tx.Unscoped().Where("file_id = ?", freshVideoFile.ID).Find(&extraResources)
			for _, er := range extraResources {
				resourceIDs[er.ID] = true
			}
		}
		for rid := range resourceIDs {
			result = tx.Unscoped().Where("id = ?", rid).Delete(&model.Resource{})
			cleanedResources += result.RowsAffected
		}

		return nil
	})
	if err != nil {
		r.Errors = append(r.Errors, "事务清理数据库记录失败: "+dirName+" - "+err.Error())
		return
	}

	r.CleanedIndexFiles += int(cleanedIndexFiles)
	r.CleanedVideoFiles += int(cleanedVideoFiles)
	r.CleanedResources += int(cleanedResources)
}

// cleanOrphanedImages 清理孤立的图片文件
// 清理条件（保守策略）：
// 1. 本地文件存在但数据库ImageFile表中无记录（真正的孤立文件）
// 2. 且该图片没有被任何内容引用
// 注意：如果ImageFile表中有记录，说明是正常上传的图片，即使暂时没被引用也保留
func (r *CleanupResult) cleanOrphanedImages(dryRun bool) {
	imageDir := "./upload/image"
	if !utils.IsFileExists(imageDir) {
		return
	}

	entries, err := os.ReadDir(imageDir)
	if err != nil {
		r.Errors = append(r.Errors, "读取图片目录失败: "+err.Error())
		return
	}

	// 支持的图片扩展名
	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".webp": true, ".bmp": true, ".ico": true, ".svg": true,
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()

		// 只处理图片文件，跳过 .gitkeep 等非图片文件
		ext := strings.ToLower(filepath.Ext(fileName))
		if !imageExts[ext] {
			continue
		}

		// 检查数据库ImageFile表中是否有记录
		var imageFileCount int64
		global.Mysql.Model(&model.ImageFile{}).Where("file_name = ?", fileName).Count(&imageFileCount)

		// 如果数据库有记录，说明是正常上传的图片，保留不删除
		if imageFileCount > 0 {
			continue
		}

		// 数据库无记录，检查是否被任何内容引用
		imageUrl := "/api/image/" + fileName
		if isImageReferenced(imageUrl) {
			continue
		}

		// 数据库无记录且无引用，可以安全删除
		localPath := filepath.Join(imageDir, fileName)
		r.Items = append(r.Items, CleanupItem{
			Type:   "image",
			Path:   fileName,
			Reason: "无数据库记录且无引用",
		})

		if !dryRun {
			// 删除OSS上的图片
			if global.Config.Storage.OssType != "local" {
				objectKey := "image/" + fileName
				utils.InfoLog("删除OSS图片: "+objectKey, "cleanup")
				if err := global.Storage.DeleteObject(objectKey); err != nil {
					r.Errors = append(r.Errors, "删除OSS图片失败: "+objectKey+" - "+err.Error())
				}
				// 同时删除备用OSS
				if global.StorageBackup != nil {
					if err := global.StorageBackup.DeleteObject(objectKey); err != nil {
						r.Errors = append(r.Errors, "删除备用OSS图片失败: "+objectKey+" - "+err.Error())
					}
				}
			}
			// 删除本地文件
			if err := os.Remove(localPath); err != nil {
				r.Errors = append(r.Errors, "删除孤立图片失败: "+localPath+" - "+err.Error())
			} else {
				r.CleanedImages++
				utils.InfoLog("已删除孤立图片(无数据库记录且无引用): "+localPath, "cleanup")
			}
		} else {
			r.CleanedImages++
		}
	}
}

// isImageReferenced 检查图片是否被引用
func isImageReferenced(imageUrl string) bool {
	// 30天前的时间（软删除30天内的数据也要保留图片）
	expireTime := time.Now().AddDate(0, 0, -30)

	// 检查Video.cover（包括30天内软删除的）
	var videoCount int64
	global.Mysql.Unscoped().Model(&model.Video{}).
		Where("cover = ? AND (deleted_at IS NULL OR deleted_at > ?)", imageUrl, expireTime).
		Count(&videoCount)
	if videoCount > 0 {
		return true
	}

	// 检查Article.cover
	var articleCount int64
	global.Mysql.Unscoped().Model(&model.Article{}).
		Where("cover = ? AND (deleted_at IS NULL OR deleted_at > ?)", imageUrl, expireTime).
		Count(&articleCount)
	if articleCount > 0 {
		return true
	}

	// 检查User.avatar和space_cover（用户不会软删除）
	var userCount int64
	global.Mysql.Model(&model.User{}).
		Where("avatar = ? OR space_cover = ?", imageUrl, imageUrl).
		Count(&userCount)
	if userCount > 0 {
		return true
	}

	// 检查Playlist.cover（包括30天内软删除的）
	var playlistCount int64
	global.Mysql.Unscoped().Model(&model.Playlist{}).
		Where("cover = ? AND (deleted_at IS NULL OR deleted_at > ?)", imageUrl, expireTime).
		Count(&playlistCount)
	if playlistCount > 0 {
		return true
	}

	// 检查Carousel.img（轮播图不会软删除）
	var carouselCount int64
	global.Mysql.Model(&model.Carousel{}).
		Where("img = ?", imageUrl).
		Count(&carouselCount)
	if carouselCount > 0 {
		return true
	}

	// 检查PGCMedia.cover（包括30天内软删除的）
	var pgcMediaCount int64
	global.Mysql.Unscoped().Model(&model.PGCMedia{}).
		Where("cover = ? AND (deleted_at IS NULL OR deleted_at > ?)", imageUrl, expireTime).
		Count(&pgcMediaCount)
	if pgcMediaCount > 0 {
		return true
	}

	// 检查PGCContent.cover（包括30天内软删除的）
	var pgcContentCount int64
	global.Mysql.Unscoped().Model(&model.PGCContent{}).
		Where("cover = ? AND (deleted_at IS NULL OR deleted_at > ?)", imageUrl, expireTime).
		Count(&pgcContentCount)
	if pgcContentCount > 0 {
		return true
	}

	// 检查Collection.cover（收藏夹会软删除，给30天宽限期）
	var collectionCount int64
	global.Mysql.Unscoped().Model(&model.Collection{}).
		Where("cover = ? AND (deleted_at IS NULL OR deleted_at > ?)", imageUrl, expireTime).
		Count(&collectionCount)
	if collectionCount > 0 {
		return true
	}

	return false
}

// deleteVideoFromOSS 删除OSS上的视频文件
func deleteVideoFromOSS(localDir string, errors *[]string) {
	if global.Config.Storage.OssType == "local" {
		return
	}

	// 从localDir中提取目录名 (例如: "./upload/video/abc123" -> "abc123")
	dirName := filepath.Base(localDir)
	if dirName == "" || dirName == "." {
		return
	}

	// 遍历本地目录，删除OSS上对应的文件
	filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// 直接使用 video/dirName/fileName 格式构建objectKey
		fileName := filepath.Base(path)
		objectKey := "video/" + dirName + "/" + fileName

		utils.InfoLog("删除OSS文件: "+objectKey, "cleanup")
		if err := global.Storage.DeleteObject(objectKey); err != nil {
			errMsg := "删除OSS文件失败: " + objectKey + " - " + err.Error()
			utils.ErrorLog(errMsg, "cleanup", err.Error())
			if errors != nil {
				*errors = append(*errors, errMsg)
			}
		}
		// 同时删除备用OSS
		if global.StorageBackup != nil {
			if err := global.StorageBackup.DeleteObject(objectKey); err != nil {
				errMsg := "删除备用OSS文件失败: " + objectKey + " - " + err.Error()
				utils.ErrorLog(errMsg, "cleanup", err.Error())
				if errors != nil {
					*errors = append(*errors, errMsg)
				}
			}
		}
		return nil
	})
}

// cleanOrphanedSubtitles 清理孤立的字幕文件
// 清理条件：
// 1. 本地字幕文件对应的 SubtitleTrack 记录不存在或被软删除
// 2. SubtitleTrack 指向的 Resource 或 Video 已被删除（物理或软删除）
func (r *CleanupResult) cleanOrphanedSubtitles(dryRun bool) {
	subtitleDir := "./upload/subtitle"
	if !utils.IsFileExists(subtitleDir) {
		return
	}

	entries, err := os.ReadDir(subtitleDir)
	if err != nil {
		r.Errors = append(r.Errors, "读取字幕目录失败: "+err.Error())
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if !strings.HasSuffix(strings.ToLower(fileName), ".vtt") {
			continue
		}

		objectKey := "subtitle/" + fileName
		localPath := filepath.Join(subtitleDir, fileName)

		// 检查 SubtitleTrack 表是否有对应记录
		var track model.SubtitleTrack
		if err := global.Mysql.Unscoped().Where("object_key = ?", objectKey).First(&track).Error; err != nil || track.ID == 0 {
			r.markSubtitleOrphaned(objectKey, localPath, "SubtitleTrack记录不存在", 0, dryRun)
			continue
		}

		// 记录存在但已被软删除
		if track.DeletedAt.Valid {
			r.markSubtitleOrphaned(objectKey, localPath, "SubtitleTrack已软删除", track.ID, dryRun)
			continue
		}

		// 检查关联的 Resource 是否存在
		var resource model.Resource
		if err := global.Mysql.Unscoped().Where("short_id = ?", track.ResourceShortID).First(&resource).Error; err != nil || resource.ID == 0 {
			r.markSubtitleOrphaned(objectKey, localPath, "关联Resource不存在", track.ID, dryRun)
			continue
		}

		// Resource 已被软删除
		if resource.DeletedAt.Valid {
			r.markSubtitleOrphaned(objectKey, localPath, "关联Resource已软删除", track.ID, dryRun)
			continue
		}

		// 检查关联的 Video 是否存在
		var video model.Video
		if err := global.Mysql.Unscoped().Where("id = ?", resource.Vid).First(&video).Error; err != nil || video.ID == 0 {
			r.markSubtitleOrphaned(objectKey, localPath, "关联Video不存在", track.ID, dryRun)
			continue
		}

		// Video 已被软删除
		if video.DeletedAt.Valid {
			r.markSubtitleOrphaned(objectKey, localPath, "关联Video已软删除", track.ID, dryRun)
			continue
		}
	}
}

// markSubtitleOrphaned 记录一条孤立字幕并执行清理
// trackID 为 0 表示 DB 无记录（不执行 DB 删除），>0 则同时清理 SubtitleTrack
func (r *CleanupResult) markSubtitleOrphaned(objectKey, localPath, reason string, trackID uint, dryRun bool) {
	r.Items = append(r.Items, CleanupItem{
		Type:   "subtitle",
		Path:   objectKey,
		Reason: reason,
	})
	if !dryRun {
		if trackID > 0 {
			global.Mysql.Unscoped().Delete(&model.SubtitleTrack{}, trackID)
		}
		deleteOrphanedSubtitleFile(objectKey, localPath, r)
	}
	r.CleanedSubtitles++
}

// deleteOrphanedSubtitleFile 删除孤立的字幕文件（本地+OSS）
func deleteOrphanedSubtitleFile(objectKey, localPath string, r *CleanupResult) {
	// 删除本地文件
	if err := os.Remove(localPath); err != nil && !os.IsNotExist(err) {
		r.Errors = append(r.Errors, "删除字幕本地文件失败: "+localPath+" - "+err.Error())
	}

	// 删除OSS上的文件
	if global.Config.Storage.OssType != "local" {
		if err := global.Storage.DeleteObject(objectKey); err != nil {
			r.Errors = append(r.Errors, "删除OSS字幕文件失败: "+objectKey+" - "+err.Error())
		}
		// 同时删除备用OSS
		if global.StorageBackup != nil {
			if err := global.StorageBackup.DeleteObject(objectKey); err != nil {
				r.Errors = append(r.Errors, "删除备用OSS字幕文件失败: "+objectKey+" - "+err.Error())
			}
		}
	}
}

// GetCleanupPreview 获取清理预览（不执行实际删除）
func GetCleanupPreview() CleanupResult {
	return CleanupOrphanedResources(true)
}

// ExecuteCleanup 执行清理
func ExecuteCleanup() CleanupResult {
	return CleanupOrphanedResources(false)
}
