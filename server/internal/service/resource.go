package service

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func ModifyResourceTitle(ctx *gin.Context, modifyTitleReq dto.ModifyResourceTitleReq) error {
	var resource model.Resource
	userId := ctx.GetUint("userId")
	global.Mysql.Model(&model.Resource{}).Where("id = ? and uid = ?", modifyTitleReq.ID, userId).First(&resource)
	if resource.ID == 0 {
		return errors.New("资源不存在")
	}

	if err := global.Mysql.Model(&model.Resource{}).Where("id = ?", modifyTitleReq.ID).Updates(
		map[string]interface{}{
			"title": modifyTitleReq.Title,
		},
	).Error; err != nil {
		utils.ErrorLog("更新资源标题失败", "resource", err.Error())
		return errors.New("更新资源标题失败")
	}

	// 清除视频信息缓存（因为视频信息中包含资源列表）
	cache.DelVideoInfo(resource.Vid)

	return nil
}

// 删除资源
func DeleteResource(ctx *gin.Context, id uint, deleteDanmaku bool) error {
	var resource model.Resource
	userId := ctx.GetUint("userId")
	global.Mysql.Model(&model.Resource{}).Where("id = ? and uid = ?", id, userId).First(&resource)
	if resource.ID == 0 {
		return errors.New("资源不存在")
	}

	var count int64
	global.Mysql.Model(&model.Resource{}).Where("vid = ?", resource.Vid).Count(&count)
	if count <= 1 {
		return errors.New("至少需要一个视频")
	}

	// 先获取video_index_file记录以获取dir_name，用于删除video_file
	var indexFile model.VideoIndexFile
	global.Mysql.Where("resource_id = ?", id).First(&indexFile)

	if err := global.Mysql.Where("id = ?", id).Delete(&model.Resource{}).Error; err != nil {
		utils.ErrorLog("删除资源失败", "resource", err.Error())
		return errors.New("删除资源失败")
	}

	// 删除关联的m3u8索引文件记录
	if err := global.Mysql.Where("resource_id = ?", id).Delete(&model.VideoIndexFile{}).Error; err != nil {
		utils.ErrorLog("删除m3u8索引文件失败", "resource", err.Error())
	}

	// 处理 VideoFile 引用计数（全局去重）
	if resource.FileID != 0 {
		decreaseVideoFileRefCount(resource.FileID, userId, resource.ID, indexFile.DirName)
	} else if indexFile.DirName != "" {
		// 兼容旧数据：通过 DirName 查找 VideoFile
		if vf, err := findVideoFileByDirName(global.Mysql, indexFile.DirName); err == nil && vf != nil {
			decreaseVideoFileRefCount(vf.ID, userId, resource.ID, indexFile.DirName)
		}
	}

	// 删除视频信息缓存（删除后让下次查询时重新从数据库加载）
	cache.DelVideoInfo(resource.Vid)

	// 删除关联弹幕（如果用户选择删除）
	if deleteDanmaku && resource.ShortID != "" {
		if err := global.Mysql.Where("rid = ?", resource.ShortID).Delete(&model.Danmaku{}).Error; err != nil {
			utils.ErrorLog("删除弹幕失败", "resource", err.Error())
		}
	}

	return nil
}

// 获取视频资源
func GetVideoResourceByStatus(videoId uint, status int) (resources []vo.ResourceResp) {
	global.Mysql.Model(&model.Resource{}).
		Where("vid = ? and status = ?", videoId, status).
		Order("sort_order ASC, id ASC").
		Scan(&resources)

	return
}

// 获取对外可见的视频资源（仅返回 VisibleStatus=1 的分P，用于公开播放页）
func GetVisibleResources(videoId uint) (resources []vo.ResourceResp) {
	global.Mysql.Model(&model.Resource{}).
		Select("id, short_id, created_at, vid, title, duration, status, file_id, uid, sort_order").
		Where("vid = ? and visible_status = ?", videoId, global.VISIBLE_SHOWN).
		Order("sort_order ASC, id ASC").
		Scan(&resources)
	return
}

// 获取视频资源（含fileId和uid用于全局去重）
func GetReviewResourceList(videoId uint) (resources []vo.ResourceResp) {
	global.Mysql.Model(&model.Resource{}).
		Select("id, short_id, created_at, vid, title, duration, status, file_id, uid, sort_order").
		Where("vid = ?", videoId).
		Order("sort_order ASC, id ASC").
		Scan(&resources)

	return
}

// GetRelatedResources 获取相同文件的关联稿件（用于审核时查看同文件其他稿件）
func GetRelatedResources(fileID uint, excludeVid uint) (resources []vo.RelatedResourceResp) {
	if fileID == 0 {
		return
	}

	// 查询相同 file_id 的其他资源（排除当前视频）
	global.Mysql.Table("resource").
		Select("resource.id as resource_id, resource.vid, resource.uid, resource.title, resource.status, resource.created_at, user.username as author_name").
		Joins("LEFT JOIN user ON resource.uid = user.id").
		Where("resource.file_id = ? AND resource.vid != ? AND resource.deleted_at IS NULL", fileID, excludeVid).
		Order("resource.created_at ASC").
		Scan(&resources)

	return
}

// ResourceQualityInfo 资源清晰度信息（包含类型）
type ResourceQualityInfo struct {
	Quality      []string `json:"quality"`
	SupportsDash bool     `json:"supportsDash"`           // 是否支持 DASH（SegmentBase 模式）
	QualityOrder []string `json:"qualityOrder,omitempty"` // 按码率升序排列（对应 dash.js Representation 索引）
}

// 获取视频资源支持的分辨率信息
func GetResourceQuality(ctx *gin.Context, id uint) (*ResourceQualityInfo, error) {
	if !IsResourceExist(id) {
		return nil, errors.New("资源不存在")
	}

	return buildResourceQualityInfo(id)
}

// 获取视频资源支持的分辨率信息(后台管理，不校验审核状态)
func GetResourceQualityManage(ctx *gin.Context, id uint) (*ResourceQualityInfo, error) {
	return buildResourceQualityInfo(id)
}

// buildResourceQualityInfo 构建资源清晰度信息（公共逻辑）
func buildResourceQualityInfo(id uint) (*ResourceQualityInfo, error) {
	var files []model.VideoIndexFile
	if err := global.Mysql.Where("resource_id = ?", id).Find(&files).Error; err != nil {
		utils.ErrorLog("分辨率信息获取失败", "resource", err.Error())
		return nil, errors.New("获取失败")
	}

	result := &ResourceQualityInfo{
		Quality:      make([]string, 0, len(files)),
		SupportsDash: false,
	}

	for _, f := range files {
		result.Quality = append(result.Quality, f.Quality)
		if f.IsSegmentBase() {
			result.SupportsDash = true
		}
	}

	// 生成按码率升序排列的清晰度列表（与统一 MPD 中 Representation 顺序一致）
	if result.SupportsDash {
		sbFiles := make([]model.VideoIndexFile, 0, len(files))
		for _, f := range files {
			if f.IsSegmentBase() {
				sbFiles = append(sbFiles, f)
			}
		}
		sort.Slice(sbFiles, func(i, j int) bool {
			return dashRepresentationLess(sbFiles[i], sbFiles[j])
		})
		result.QualityOrder = make([]string, 0, len(sbFiles))
		for _, f := range sbFiles {
			result.QualityOrder = append(result.QualityOrder, f.Quality)
		}
	}

	return result, nil
}

// 资源排序
func ReorderResources(ctx *gin.Context, req dto.ReorderResourceReq) error {
	userId := ctx.GetUint("userId")

	// 验证视频属于当前用户
	var video model.Video
	global.Mysql.Model(&model.Video{}).Where("id = ? and uid = ?", req.Vid, userId).First(&video)
	if video.ID == 0 {
		return errors.New("视频不存在")
	}

	tx := global.Mysql.Begin()
	if tx.Error != nil {
		utils.ErrorLog("开启排序事务失败", "resource", tx.Error.Error())
		return errors.New("更新排序失败")
	}

	// 校验：请求中的资源ID必须与该视频资源集合完全一致（防止部分更新造成顺序错乱）
	var dbResourceIds []uint
	if err := tx.Model(&model.Resource{}).
		Where("vid = ? AND uid = ?", req.Vid, userId).
		Order("id ASC").
		Pluck("id", &dbResourceIds).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("查询资源列表失败", "resource", err.Error())
		return errors.New("更新排序失败")
	}
	if len(dbResourceIds) == 0 {
		tx.Rollback()
		return errors.New("资源不存在")
	}
	if len(dbResourceIds) != len(req.ResourceIds) {
		tx.Rollback()
		return errors.New("资源数量不匹配")
	}

	dbIdSet := make(map[uint]struct{}, len(dbResourceIds))
	for _, id := range dbResourceIds {
		dbIdSet[id] = struct{}{}
	}
	reqIdSet := make(map[uint]struct{}, len(req.ResourceIds))
	for _, id := range req.ResourceIds {
		if _, exists := reqIdSet[id]; exists {
			tx.Rollback()
			return errors.New("资源ID重复")
		}
		reqIdSet[id] = struct{}{}
		if _, exists := dbIdSet[id]; !exists {
			tx.Rollback()
			return errors.New("存在无效资源ID")
		}
	}

	// 批量更新 sort_order
	for i, resourceId := range req.ResourceIds {
		result := tx.Model(&model.Resource{}).
			Where("id = ? AND vid = ? AND uid = ?", resourceId, req.Vid, userId).
			Update("sort_order", i)
		if result.Error != nil {
			tx.Rollback()
			utils.ErrorLog("更新资源排序失败", "resource", result.Error.Error())
			return errors.New("更新排序失败")
		}
		if result.RowsAffected != 1 {
			tx.Rollback()
			utils.ErrorLog("更新资源排序影响行数异常", "resource", fmt.Sprintf("resourceId=%d rows=%d", resourceId, result.RowsAffected))
			return errors.New("更新排序失败")
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.ErrorLog("提交排序事务失败", "resource", err.Error())
		return errors.New("更新排序失败")
	}

	// 清除视频信息缓存
	cache.DelVideoInfo(req.Vid)

	return nil
}

func IsResourceExist(resourceId uint) bool {
	var resource model.Resource
	global.Mysql.Model(&model.Resource{}).Select("id").
		Where("id = ? and `status` = ?", resourceId, global.AUDIT_APPROVED).First(&resource)

	return resource.ID != 0
}

// 检查资源是否需要替换（比较hash）
func CheckReplaceResource(ctx *gin.Context, replaceReq dto.ReplaceResourceReq) error {
	userId := ctx.GetUint("userId")

	// 验证资源是否属于当前用户
	var oldResource model.Resource
	global.Mysql.Model(&model.Resource{}).Where("id = ? and uid = ?", replaceReq.ResourceID, userId).First(&oldResource)
	if oldResource.ID == 0 {
		return errors.New("资源不存在")
	}

	// 通过 Resource.FileID 或 VideoIndexFile 查找旧视频文件
	var oldFileInfo model.VideoFile
	if oldResource.FileID != 0 {
		global.Mysql.Where("id = ?", oldResource.FileID).First(&oldFileInfo)
	}
	if oldFileInfo.ID == 0 {
		// 兼容旧数据：通过 VideoIndexFile 查找
		var oldIndexFile model.VideoIndexFile
		global.Mysql.Where("resource_id = ?", replaceReq.ResourceID).First(&oldIndexFile)
		if oldIndexFile.DirName != "" {
			global.Mysql.Where("dir_name = ?", oldIndexFile.DirName).First(&oldFileInfo)
		}
	}

	// 检查哈希值是否相同
	if oldFileInfo.ID != 0 && oldFileInfo.Hash == replaceReq.Hash {
		// 哈希值相同，检查源文件和转码产物是否都存在
		suffix := utils.GetFileSuffix(oldFileInfo.OriginalName)
		localPath := "./upload/video/" + oldFileInfo.DirName + "/upload" + suffix
		if !utils.IsFileExists(localPath) {
			// 源文件已丢失，允许重新替换上传
			return nil
		}
		// 源文件存在，还需检查转码产物是否完整
		var indexCount int64
		global.Mysql.Model(&model.VideoIndexFile{}).
			Where("resource_id = ?", replaceReq.ResourceID).Count(&indexCount)
		if indexCount == 0 {
			// 转码记录不存在（转码产物丢失），允许重新替换上传
			return nil
		}
		return errors.New("视频文件相同，无需替换")
	}

	// 可以替换
	return nil
}

// 替换资源
func ReplaceResource(ctx *gin.Context, replaceReq dto.ReplaceResourceReq) (vo.ResourceResp, error) {
	userId := ctx.GetUint("userId")

	// 验证资源是否属于当前用户
	var oldResource model.Resource
	global.Mysql.Model(&model.Resource{}).Where("id = ? and uid = ?", replaceReq.ResourceID, userId).First(&oldResource)
	if oldResource.ID == 0 {
		return vo.ResourceResp{}, errors.New("资源不存在")
	}

	// 获取旧资源的索引文件信息
	var oldIndexFile model.VideoIndexFile
	global.Mysql.Where("resource_id = ?", replaceReq.ResourceID).First(&oldIndexFile)

	// 通过 Resource.FileID 或 VideoIndexFile 查找旧视频文件
	var oldFileInfo model.VideoFile
	if oldResource.FileID != 0 {
		global.Mysql.Where("id = ?", oldResource.FileID).First(&oldFileInfo)
	}
	if oldFileInfo.ID == 0 && oldIndexFile.DirName != "" {
		global.Mysql.Where("dir_name = ?", oldIndexFile.DirName).First(&oldFileInfo)
	}
	// 获取旧文件的 DirName（用于后续清理）
	oldDirName := oldIndexFile.DirName
	if oldDirName == "" && oldFileInfo.ID != 0 {
		oldDirName = oldFileInfo.DirName
	}

	// 获取新视频文件信息
	var newFileInfo model.VideoFile
	// 全局去重：替换资源时按文件哈希查找，不限制上传者
	if err := global.Mysql.Where("hash = ?", replaceReq.Hash).Order("id DESC").First(&newFileInfo).Error; err != nil || newFileInfo.ID == 0 {
		utils.ErrorLog("新视频文件信息不存在", "resource", replaceReq.Hash)
		return vo.ResourceResp{}, errors.New("视频文件不存在")
	}

	// 检查哈希值是否相同且旧文件和转码产物都完整
	if oldFileInfo.ID != 0 && oldFileInfo.Hash == replaceReq.Hash {
		suffix := utils.GetFileSuffix(oldFileInfo.OriginalName)
		localPath := "./upload/video/" + oldFileInfo.DirName + "/upload" + suffix
		var indexCount int64
		global.Mysql.Model(&model.VideoIndexFile{}).
			Where("resource_id = ?", replaceReq.ResourceID).Count(&indexCount)
		if utils.IsFileExists(localPath) && indexCount > 0 {
			// 源文件和转码产物都存在，无需转码，直接更新为待审核状态
			if err := global.Mysql.Model(&model.Resource{}).Where("id = ?", replaceReq.ResourceID).Update("status", global.WAITING_REVIEW).Error; err != nil {
				utils.ErrorLog("更新资源状态失败", "resource", err.Error())
				return vo.ResourceResp{}, errors.New("更新资源状态失败")
			}
			var newResource model.Resource
			global.Mysql.Where("id = ?", replaceReq.ResourceID).First(&newResource)
			utils.InfoLog(fmt.Sprintf("【替换资源】hash相同且文件完整，跳过转码，设为待审核 resourceID=%d", replaceReq.ResourceID), "resource")
			return vo.ResourceToResourceResp(newResource), nil
		}
		// 源文件或转码产物丢失，需要重新转码
		utils.InfoLog(fmt.Sprintf("【替换资源】hash相同但文件不完整(源文件=%v,转码记录=%d)，重新转码 resourceID=%d",
			utils.IsFileExists(localPath), indexCount, replaceReq.ResourceID), "resource")
	}

	// 处理旧文件的引用计数（使用统一的引用计数机制）
	if oldResource.FileID != 0 {
		decreaseVideoFileRefCount(oldResource.FileID, userId, replaceReq.ResourceID, oldDirName)
	} else if oldDirName != "" {
		// 兼容旧数据：通过 DirName 查找 VideoFile
		var oldVf model.VideoFile
		if global.Mysql.Where("dir_name = ?", oldDirName).First(&oldVf).Error == nil {
			decreaseVideoFileRefCount(oldVf.ID, userId, replaceReq.ResourceID, oldDirName)
		}
	}

	// 删除旧的索引文件记录
	if err := global.Mysql.Where("resource_id = ?", replaceReq.ResourceID).Delete(&model.VideoIndexFile{}).Error; err != nil {
		utils.ErrorLog("删除旧索引文件失败", "resource", err.Error())
	}

	// 创建新文件的引用
	if err := createFileRef(userId, newFileInfo.ID, replaceReq.ResourceID); err != nil {
		return vo.ResourceResp{}, errors.New("创建文件引用失败")
	}

	// 读取新视频信息
	suffix := utils.GetFileSuffix(newFileInfo.OriginalName)
	uploadVideoPath := "./upload/video/" + newFileInfo.DirName + "/upload" + suffix
	transcodingInfo, err := ProcessVideoInfo(uploadVideoPath)
	if err != nil {
		return vo.ResourceResp{}, errors.New("读取视频信息失败")
	}

	// 更新资源记录（包含 file_id）
	if err := global.Mysql.Model(&model.Resource{}).Where("id = ?", replaceReq.ResourceID).Updates(
		map[string]interface{}{
			"codec_name": transcodingInfo.CodecName,
			"status":     global.VIDEO_PROCESSING,
			"duration":   utils.SecFromFloat(transcodingInfo.Duration),
			"file_id":    newFileInfo.ID,
		},
	).Error; err != nil {
		utils.ErrorLog("更新资源记录失败", "resource", err.Error())
		return vo.ResourceResp{}, errors.New("更新资源失败")
	}

	// 启动转码服务
	transcodingInfo.VideoID = oldResource.Vid
	transcodingInfo.DirName = newFileInfo.DirName
	transcodingInfo.ResourceID = replaceReq.ResourceID
	transcodingInfo.OutputDir = "./upload/video/" + newFileInfo.DirName + "/"
	transcodingInfo.InputFile = transcodingInfo.OutputDir + "upload" + suffix
	transcodingInfo.Suffix = suffix
	if err := GetCurrentTranscoder().Enqueue(context.Background(), transcodingInfo); err != nil {
		utils.ErrorLog("资源替换转码入队失败", "resource",
			fmt.Sprintf("ResourceID=%d, err=%v", replaceReq.ResourceID, err))
	}

	// 清除视频信息缓存
	cache.DelVideoInfo(oldResource.Vid)

	// 返回更新后的资源信息
	var updatedResource model.Resource
	global.Mysql.First(&updatedResource, replaceReq.ResourceID)
	return vo.ResourceToResourceResp(updatedResource), nil
}

// GetResourceShortIDByPart 根据分P序号获取资源的ShortID
// 用于历史记录绑定到具体资源，不受排序影响
func GetResourceShortIDByPart(videoId uint, part uint) (string, error) {
	if part == 0 {
		part = 1
	}
	// 按 sort_order 和 id 排序，获取第 part-1 个资源
	var resource model.Resource
	err := global.Mysql.Model(&model.Resource{}).
		Where("vid = ?", videoId).
		Order("sort_order ASC, id ASC").
		Offset(int(part - 1)).
		First(&resource).Error
	if err != nil {
		return "", err
	}
	return resource.ShortID, nil
}
