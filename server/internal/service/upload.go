package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"gorm.io/gorm/clause"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	osslib "interastral-peace.com/alnitak/pkg/oss"
	"interastral-peace.com/alnitak/utils"
)

// truncateString 按字符数截断字符串（支持中文等多字节字符）
func truncateString(s string, maxLen int) string {
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxLen])
}

func UploadImg(ctx *gin.Context, file *multipart.FileHeader) (string, error) {
	suffix := path.Ext(file.Filename)
	userId := ctx.GetUint("userId")

	// 参数校验
	if !utils.IsImgType(suffix, global.Config.File.AllowedImgExts) { // 文件后缀
		return "", errors.New("文件类型错误")
	}

	// 魔数校验：读取文件头确认真实类型
	if isImg, err := utils.CheckImageMagicBytes(file); err != nil || !isImg {
		return "", errors.New("文件类型错误")
	}

	//文件大小限制
	if !utils.FileSize(file.Size, 1, global.Config.File.MaxImgSize) {
		return "", errors.New("文件大小超出限制")
	}

	// 计算文件hash
	fileHash, err := calculateFileHash(file)
	if err != nil {
		return "", errors.New("计算文件hash失败")
	}

	// 查询是否已存在相同hash的图片
	var existingFile model.ImageFile
	global.Mysql.Where("hash = ?", fileHash).First(&existingFile)
	if existingFile.ID != 0 {
		// 已存在相同图片，直接返回已有的URL
		url := generateFileUrl("image/" + existingFile.FileName)
		cache.SetUploadImage(url, userId)
		return url, nil
	}

	// 生成新文件名并保存
	fileName := generateImgFilename(suffix)
	objectKey := "image/" + fileName
	filePath := "./upload/image/" + fileName

	//保存文件
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		return "", errors.New("文件上传失败")
	}

	// 上传到 OSS（非 local）必须失败可感知，避免 DB 出现“无对象”的脏记录
	if global.Config.Storage.OssType != "local" {
		if err := global.Storage.PutObjectFromFile(objectKey, filePath); err != nil {
			// 本地临时文件清理，避免磁盘泄漏
			_ = os.Remove(filePath)
			utils.ErrorLog("图片上传到OSS失败", "upload", err.Error())
			return "", errors.New("文件上传失败")
		}
		// 上传到备用 OSS（带重试 + 失败持久化）
		go UploadToBackupWithRetry(objectKey, filePath, "image")
	}

	// 记录到数据库
	global.Mysql.Create(&model.ImageFile{
		Uid:      userId,
		FileName: fileName,
		Hash:     fileHash,
	})

	url := generateFileUrl(objectKey)

	// 缓存url
	cache.SetUploadImage(url, userId)

	return url, nil
}

// 计算上传文件的MD5 hash
func calculateFileHash(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, src); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// UploadVideoCreate 创建新视频稿件（全局去重版本）
func UploadVideoCreate(ctx *gin.Context, videoFileReq dto.VideoFileReq) (vo.ResourceResp, error) {
	userId := ctx.GetUint("userId")

	var fileInfo model.VideoFile

	// 支持 FileID(DirName) 和 Hash 两种方式查询
	if videoFileReq.FileID != "" {
		if err := global.Mysql.Unscoped().Where("dir_name = ?", videoFileReq.FileID).First(&fileInfo).Error; err != nil {
			utils.ErrorLog("视频文件信息不存在", "upload", fmt.Sprintf("fileID=%s", videoFileReq.FileID))
			return vo.ResourceResp{}, errors.New("视频文件不存在")
		}
	} else {
		// 【全局去重】按 hash + size 查询，包含软删除的记录
		if err := global.Mysql.Unscoped().Where("hash = ? AND size = ?", videoFileReq.Hash, videoFileReq.Size).First(&fileInfo).Error; err != nil {
			utils.ErrorLog("视频文件信息不存在", "upload", fmt.Sprintf("hash=%s, size=%d", videoFileReq.Hash, videoFileReq.Size))
			return vo.ResourceResp{}, errors.New("视频文件不存在")
		}
	}

	// 如果是软删除状态，恢复它
	if fileInfo.DeletedAt.Valid {
		if err := global.Mysql.Unscoped().Model(&fileInfo).Update("deleted_at", nil).Error; err != nil {
			utils.ErrorLog("恢复软删除记录失败", "upload", err.Error())
		}
	}

	// 先创建视频记录
	suffix := utils.GetFileSuffix(fileInfo.OriginalName)
	uploadVideoPath := "./upload/video/" + fileInfo.DirName + "/upload" + suffix
	vid, err := initVideo(userId, uploadVideoPath, fileInfo.OriginalName, videoFileReq.Cover)
	if err != nil || vid == 0 {
		utils.ErrorLog("创建视频失败", "upload", fmt.Sprintf("uid=%d, dirName=%s, originalName=%s, err=%v", userId, fileInfo.DirName, fileInfo.OriginalName, err))
		return vo.ResourceResp{}, errors.New("创建失败")
	}

	// 创建用户引用关系
	if err := createFileRef(userId, fileInfo.ID, 0); err != nil {
		return vo.ResourceResp{}, errors.New("创建文件引用失败")
	}

	// 构建前端元数据
	var meta *VideoMeta
	if videoFileReq.Duration > 0 {
		meta = &VideoMeta{
			Duration:  videoFileReq.Duration,
			Width:     videoFileReq.Width,
			Height:    videoFileReq.Height,
			CodecName: videoFileReq.CodecName,
		}
	}

	resource, err := CompleteUploadVideo(vid, userId, fileInfo.ID, fileInfo.DirName, fileInfo.OriginalName, fileInfo.Status == model.FileStatusReady, meta)
	if err != nil {
		// 补偿：创建资源失败时回滚临时引用关系，避免 resource_id=0 悬挂记录
		decreaseVideoFileRefCount(fileInfo.ID, userId, 0, fileInfo.DirName)
		return vo.ResourceResp{}, err
	}

	// 更新引用关系中的 ResourceID
	updateFileRefResourceID(userId, fileInfo.ID, resource.ID)

	return resource, nil
}

// UploadVideoAdd 向已有视频添加分P（全局去重版本）
func UploadVideoAdd(ctx *gin.Context, vid uint, videoFileReq dto.VideoFileReq) (vo.ResourceResp, error) {
	userId := ctx.GetUint("userId")

	// 【全局去重】按 hash + size 查询，包含软删除的记录
	var fileInfo model.VideoFile
	if err := global.Mysql.Unscoped().Where("hash = ? AND size = ?", videoFileReq.Hash, videoFileReq.Size).First(&fileInfo).Error; err != nil {
		utils.ErrorLog("视频文件信息不存在", "upload", fmt.Sprintf("hash=%s, size=%d", videoFileReq.Hash, videoFileReq.Size))
		return vo.ResourceResp{}, errors.New("视频文件不存在")
	}

	// 如果是软删除状态，恢复它
	if fileInfo.DeletedAt.Valid {
		if err := global.Mysql.Unscoped().Model(&fileInfo).Update("deleted_at", nil).Error; err != nil {
			utils.ErrorLog("恢复软删除记录失败", "upload", err.Error())
		}
	}

	// 创建用户引用关系
	if err := createFileRef(userId, fileInfo.ID, 0); err != nil {
		return vo.ResourceResp{}, errors.New("创建文件引用失败")
	}

	var meta *VideoMeta
	if videoFileReq.Duration > 0 {
		meta = &VideoMeta{
			Duration:  videoFileReq.Duration,
			Width:     videoFileReq.Width,
			Height:    videoFileReq.Height,
			CodecName: videoFileReq.CodecName,
		}
	}

	resource, err := CompleteUploadVideo(vid, userId, fileInfo.ID, fileInfo.DirName, fileInfo.OriginalName, fileInfo.Status == model.FileStatusReady, meta, videoFileReq.ReplaceResourceID)
	if err != nil {
		// 补偿：创建资源失败时回滚临时引用关系，避免 resource_id=0 悬挂记录
		decreaseVideoFileRefCount(fileInfo.ID, userId, 0, fileInfo.DirName)
		return vo.ResourceResp{}, err
	}

	// 更新引用关系中的 ResourceID
	updateFileRefResourceID(userId, fileInfo.ID, resource.ID)

	return resource, nil
}

// UploadVideoCheck 检查视频上传进度（全局去重版本）
// 返回值: checks=已上传的分片索引, fileID=视频文件ID
// 返回 chunks=[-1] 表示秒传成功，客户端可以直接调用创建接口
func UploadVideoCheck(ctx *gin.Context, videoFileReq dto.VideoFileReq) (dto.VideoCheckResp, error) {
	resp := dto.VideoCheckResp{
		Chunks: []int{},
		FileID: "",
	}

	// 【全局去重】按 hash + size 查询，包含软删除的记录
	var fileInfo model.VideoFile
	result := global.Mysql.Unscoped().Where("hash = ? AND size = ?", videoFileReq.Hash, videoFileReq.Size).First(&fileInfo)

	// 文件不存在，返回空列表（需要从头上传）
	if result.Error != nil {
		utils.InfoLog(fmt.Sprintf("【秒传检测】hash=%s, size=%d, VideoFile不存在，需要上传", videoFileReq.Hash, videoFileReq.Size), "upload")
		return resp, nil
	}

	resp.FileID = fileInfo.DirName

	// 如果是软删除状态，恢复它
	if fileInfo.DeletedAt.Valid {
		if err := global.Mysql.Unscoped().Model(&fileInfo).Update("deleted_at", nil).Error; err != nil {
			utils.ErrorLog("恢复软删除记录失败", "upload", err.Error())
		} else {
			utils.InfoLog(fmt.Sprintf("【秒传检测】hash=%s, 恢复软删除记录 fileID=%s", videoFileReq.Hash, fileInfo.DirName), "upload")
		}
	}

	utils.InfoLog(fmt.Sprintf("【秒传检测】hash=%s, size=%d, fileID=%s, Status=%d, RefCount=%d, DirName=%s",
		videoFileReq.Hash, videoFileReq.Size, fileInfo.DirName, fileInfo.Status, fileInfo.RefCount, fileInfo.DirName), "upload")

	// 【秒传判断】根据文件状态决定处理方式
	switch fileInfo.Status {
	case model.FileStatusReady:
		// 检查本地文件是否存在
		suffix := utils.GetFileSuffix(fileInfo.OriginalName)
		uploadVideoPath := "./upload/video/" + fileInfo.DirName + "/upload" + suffix
		if !utils.IsFileExists(uploadVideoPath) {
			// 本地文件不存在，需要重新上传
			utils.InfoLog(fmt.Sprintf("【秒传检测】hash=%s, size=%d, 文件状态Ready但本地文件不存在，需要重新上传", videoFileReq.Hash, videoFileReq.Size), "upload")
			// 更新文件状态为上传中（文件已丢失，需要从头上传）
			global.Mysql.Model(&fileInfo).Update("status", model.FileStatusUploading)
		} else {
			// 检查是否存在关联的 Resource 记录（确保有转码信息可复用）
			// 优先查 file_id 关联，如果没有则通过 VideoIndexFile 的 DirName 关联查询
			var resourceCount int64
			global.Mysql.Model(&model.Resource{}).Where("file_id = ?", fileInfo.ID).Count(&resourceCount)
			utils.InfoLog(fmt.Sprintf("【秒传检测】通过FileID查询: fileID=%s 关联Resource数量=%d", fileInfo.DirName, resourceCount), "upload")

			if resourceCount > 0 {
				// 有已关联资源，可以秒传
				utils.InfoLog(fmt.Sprintf("【秒传成功】hash=%s, size=%d, fileID=%s, 返回[-1]", videoFileReq.Hash, videoFileReq.Size, fileInfo.DirName), "upload")
				resp.Chunks = []int{-1}
				return resp, nil
			}

			// 兼容旧数据：Resource.FileID 为 0，通过 VideoIndexFile.DirName 查找对应的资源
			var videoIndex model.VideoIndexFile
			if err := global.Mysql.Where("dir_name = ?", fileInfo.DirName).First(&videoIndex).Error; err == nil {
				utils.InfoLog(fmt.Sprintf("【秒传检测】通过DirName查询: DirName=%s, ResourceID=%d", fileInfo.DirName, videoIndex.ResourceID), "upload")
				if videoIndex.ResourceID > 0 {
					// 找到了通过 DirName 关联的资源，可以秒传
					// 同时修复 Resource.FileID 以便后续直接使用
					global.Mysql.Model(&model.Resource{}).Where("id = ? AND file_id = 0", videoIndex.ResourceID).Update("file_id", fileInfo.ID)
					utils.InfoLog(fmt.Sprintf("【秒传成功】hash=%s, size=%d, fileID=%s (通过DirName关联), 返回[-1]", videoFileReq.Hash, videoFileReq.Size, fileInfo.DirName), "upload")
					resp.Chunks = []int{-1}
					return resp, nil
				}
			}

			// 文件状态是Ready但没有Resource，本地文件存在，直接返回[-1]让前端创建
			// CompleteUploadVideo 会因查不到已有Resource而走正常转码流程
			utils.InfoLog(fmt.Sprintf("【秒传检测】hash=%s, size=%d, fileID=%s, 状态Ready但无可用Resource，返回[-1]走转码", videoFileReq.Hash, videoFileReq.Size, fileInfo.DirName), "upload")
			resp.Chunks = []int{-1}
			return resp, nil
		}

	case model.FileStatusMerged, model.FileStatusTranscoding:
		// 文件已合并完成（但可能提交失败没创建Resource，或转码中断）
		// 检查合并后的文件是否存在
		suffix := utils.GetFileSuffix(fileInfo.OriginalName)
		uploadVideoPath := "./upload/video/" + fileInfo.DirName + "/upload" + suffix
		if utils.IsFileExists(uploadVideoPath) {
			// 合并文件存在，前端可以直接调创建接口，不需要重新上传分片
			// 同时确保 OSS 上有源文件，避免远程 Worker 拉不到
			uploadMergedVideoToOSS(fileInfo.DirName, suffix, uploadVideoPath)
			utils.InfoLog(fmt.Sprintf("【秒传检测】hash=%s, Status=%d, 合并文件已存在，返回[-1]", videoFileReq.Hash, fileInfo.Status), "upload")
			resp.Chunks = []int{-1}
			return resp, nil
		}
		// 合并文件不存在，需要重新上传
		utils.InfoLog(fmt.Sprintf("【秒传检测】hash=%s, Status=%d, 合并文件不存在，重置为Uploading", videoFileReq.Hash, fileInfo.Status), "upload")
		global.Mysql.Model(&fileInfo).Update("status", model.FileStatusUploading)

	default:
		utils.InfoLog(fmt.Sprintf("【秒传检测】hash=%s, Status=%d，继续上传流程", videoFileReq.Hash, fileInfo.Status), "upload")
	}

	// 文件正在上传中，返回已上传的分片
	resp.Chunks = []int{}
	fileDir := "./upload/video/" + fileInfo.DirName
	for i := 0; i < fileInfo.ChunksCount; i++ {
		if utils.IsFileExists(fmt.Sprintf("%s/chunks/%d.part", fileDir, i)) {
			resp.Chunks = append(resp.Chunks, i)
		}
	}

	return resp, nil
}

// UploadVideoChunk 上传视频分片
func UploadVideoChunk(ctx *gin.Context, file *multipart.FileHeader) error {
	userId := ctx.GetUint("userId")

	fileHash := strings.TrimSpace(ctx.PostForm("hash"))
	fileName := strings.TrimSpace(ctx.PostForm("name"))
	fileSize := strings.TrimSpace(ctx.PostForm("size"))

	chunkIndex, err := strconv.Atoi(ctx.PostForm("chunkIndex"))
	if err != nil {
		return errors.New("非法分片索引")
	}

	totalChunks, err := strconv.Atoi(ctx.PostForm("totalChunks"))
	if err != nil || totalChunks <= 0 {
		return errors.New("非法分片总数")
	}

	size, err := strconv.ParseInt(fileSize, 10, 64)
	if err != nil || size <= 0 {
		return errors.New("非法文件大小")
	}

	// 1️⃣ 校验文件类型
	suffix := path.Ext(fileName)
	if !utils.IsVideoType(suffix, global.Config.File.AllowedVideoExts) {
		return errors.New("不支持的视频格式")
	}

	// 首分片时校验魔数，确保文件头部包含真实视频格式标识
	if chunkIndex == 0 {
		if isVideo, err := utils.CheckVideoMagicBytes(file); err != nil || !isVideo {
			return errors.New("不支持的视频格式")
		}
	}

	// 2️⃣ 查询或创建文件记录（hash + size 全局唯一）
	var videoFile model.VideoFile
	err = global.Mysql.Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("hash = ? AND size = ?", fileHash, size).First(&videoFile)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			videoFile = model.VideoFile{
				Hash:         fileHash,
				Size:         size,
				DirName:      generateVideoFilename(),
				OriginalName: truncateString(fileName, 255),
				ChunksCount:  totalChunks,
				Status:       model.FileStatusUploading,
				UploaderUid:  userId,
			}
			if err := tx.Create(&videoFile).Error; err != nil {
				return err
			}
		} else if result.Error != nil {
			return result.Error
		}

		// 3️⃣ 校验 totalChunks 不允许被篡改
		if videoFile.ChunksCount != totalChunks {
			return errors.New("分片数量不匹配")
		}
		return nil
	})
	if err != nil {
		utils.ErrorLog("查询或创建视频记录失败", "upload", err.Error())
		return errors.New("文件上传失败")
	}

	// 4️⃣ 秒传判断
	if videoFile.Status == model.FileStatusReady {
		utils.InfoLog("秒传成功："+fileHash, "upload")
		chunkDir := filepath.Join("./upload/video", videoFile.DirName, "chunks")
		_ = os.RemoveAll(chunkDir)
		return nil
	}

	// 5️⃣ 校验 chunkIndex 合法性
	if chunkIndex < 0 || chunkIndex >= totalChunks {
		return errors.New("非法分片索引")
	}

	fileDir := filepath.Join("./upload/video", videoFile.DirName)
	chunkDir := filepath.Join(fileDir, "chunks")

	// 6️⃣ 创建目录
	if err := os.MkdirAll(chunkDir, 0755); err != nil {
		return errors.New("服务器存储异常")
	}

	chunkPath := filepath.Join(chunkDir, fmt.Sprintf("%d.part", chunkIndex))

	// 7️⃣ 幂等检查
	if utils.IsFileExists(chunkPath) {
		return nil
	}

	// 8️⃣ 保存分片
	if err := ctx.SaveUploadedFile(file, chunkPath); err != nil {
		return errors.New("分片保存失败")
	}

	return nil
}

// UploadVideoMerge 合并视频分片（全局去重版本）
func UploadVideoMerge(ctx *gin.Context, videoFileReq dto.VideoFileReq) error {
	var fileInfo model.VideoFile

	// 1️⃣ 查询文件：优先按 fileID(DirName)，fallback 到 hash+size
	if videoFileReq.FileID != "" {
		if err := global.Mysql.Where("dir_name = ?", videoFileReq.FileID).First(&fileInfo).Error; err != nil {
			return errors.New("视频文件不存在")
		}
	} else if videoFileReq.Hash != "" && videoFileReq.Size > 0 {
		if err := global.Mysql.Where("hash = ? AND size = ?", videoFileReq.Hash, videoFileReq.Size).First(&fileInfo).Error; err != nil {
			return errors.New("视频文件不存在")
		}
	} else {
		return errors.New("缺少文件标识参数")
	}

	fileDir := filepath.Join("./upload/video", fileInfo.DirName)
	chunkDir := filepath.Join(fileDir, "chunks")

	// 2️⃣ 原子切换状态：Uploading → Merging
	result := global.Mysql.Model(&fileInfo).
		Where("status = ?", model.FileStatusUploading).
		Update("status", model.FileStatusMerging)
	if result.RowsAffected == 0 {
		// 状态不是Uploading，可能已经合并过或正在合并
		// 只在已完成合并（Merged/Ready）时清理分片，避免干扰正在进行的合并
		if fileInfo.Status == model.FileStatusMerged || fileInfo.Status == model.FileStatusReady || fileInfo.Status == model.FileStatusTranscoding {
			_ = os.RemoveAll(chunkDir)
		}
		return nil
	}
	suffix := utils.GetFileSuffix(fileInfo.OriginalName)
	outputFile := filepath.Join(fileDir, "upload"+suffix)

	// 3️⃣ 校验分片数量
	files, err := os.ReadDir(chunkDir)
	if err != nil {
		global.Mysql.Model(&fileInfo).Update("status", model.FileStatusUploading)
		return errors.New("分片目录不存在")
	}
	if len(files) != fileInfo.ChunksCount {
		global.Mysql.Model(&fileInfo).
			Update("status", model.FileStatusUploading)
		return errors.New("分片数量不完整")
	}

	// 4️⃣ 执行流式合并
	if err := utils.MergeChunks(chunkDir, fileInfo.ChunksCount, outputFile); err != nil {
		global.Mysql.Model(&fileInfo).
			Update("status", model.FileStatusUploading)
		return errors.New("合并分片失败")
	}

	// 5️⃣ 二次校验 hash
	calculatedHash, err := utils.CalculateFileHash(outputFile)
	if err != nil || calculatedHash != fileInfo.Hash {
		os.Remove(outputFile)
		global.Mysql.Model(&fileInfo).
			Update("status", model.FileStatusUploading)
		return errors.New("文件校验失败")
	}

	// 6️⃣ 上传合并后的文件到 OSS（非 local 模式），供远程 Worker 拉取转码
	if global.Config.Storage.OssType != "local" {
		objectKey := fmt.Sprintf("video/%s/upload%s", fileInfo.DirName, suffix)
		if err := global.Storage.PutObjectFromFile(objectKey, outputFile); err != nil {
			utils.ErrorLog("上传视频文件到OSS失败", "upload", err.Error())
			// 不回滚合并状态：文件已在本地存在，可后续重试上传
			// 标记为 Merged 但无 OSS 对象，Worker 会报 key not found
			global.Mysql.Model(&fileInfo).
				Update("status", model.FileStatusMerged)
			return errors.New("上传视频文件到OSS失败")
		}
		// 上传到备用 OSS（带重试 + 失败持久化）
		go UploadToBackupWithRetry(objectKey, outputFile, "video")
	}

	// 7️⃣ 更新最终状态
	global.Mysql.Model(&fileInfo).
		Update("status", model.FileStatusMerged)

	// 8️⃣ 删除分片目录
	_ = os.RemoveAll(chunkDir)

	return nil
}

// createFileRef 创建用户-文件引用关系（事务性操作）
func createFileRef(uid, fileID, resourceID uint) error {
	if err := global.Mysql.Transaction(func(tx *gorm.DB) error {
		ref := model.VideoFileRef{
			Uid:        uid,
			FileID:     fileID,
			ResourceID: resourceID,
		}
		if err := tx.Create(&ref).Error; err != nil {
			return err
		}
		// 增加引用计数
		return tx.Model(&model.VideoFile{}).Where("id = ?", fileID).
			UpdateColumn("ref_count", gorm.Expr("ref_count + 1")).Error
	}); err != nil {
		utils.ErrorLog(fmt.Sprintf("创建文件引用失败, fileID=%d, uid=%d", fileID, uid), "upload", err.Error())
		return err
	}
	return nil
}

// updateFileRefResourceID 更新引用关系中的 ResourceID
func updateFileRefResourceID(uid, fileID, resourceID uint) {
	global.Mysql.Model(&model.VideoFileRef{}).
		Where("uid = ? AND file_id = ? AND resource_id = 0", uid, fileID).
		Update("resource_id", resourceID)
}

// decreaseVideoFileRefCount 减少视频文件引用计数，如果计数为0则删除文件记录
// 用于删除稿件时正确处理全局去重的文件引用（事务性操作）
func decreaseVideoFileRefCount(fileID, uid, resourceID uint, dirName string) {
	if fileID == 0 {
		return
	}

	shouldDeleteDir := false
	if err := global.Mysql.Transaction(func(tx *gorm.DB) error {
		// 删除用户的引用记录
		if err := tx.Where("file_id = ? AND uid = ? AND resource_id = ?", fileID, uid, resourceID).
			Delete(&model.VideoFileRef{}).Error; err != nil {
			return fmt.Errorf("删除VideoFileRef失败: %w", err)
		}

		// 减少引用计数
		updateResult := tx.Model(&model.VideoFile{}).Where("id = ? AND ref_count > 0", fileID).
			UpdateColumn("ref_count", gorm.Expr("ref_count - 1"))
		if updateResult.Error != nil {
			return fmt.Errorf("减少引用计数失败: %w", updateResult.Error)
		}
		// 没有实际扣减时，直接退出，避免基于脏 ref_count 误删 VideoFile
		if updateResult.RowsAffected == 0 {
			utils.ErrorLog(fmt.Sprintf("【引用计数未扣减】fileID=%d, uid=%d, resourceID=%d", fileID, uid, resourceID), "upload", "")
			return nil
		}

		// 在事务内检查引用计数，使用行锁防止并发问题
		var vf model.VideoFile
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", fileID).First(&vf).Error; err != nil {
			return nil // 记录不存在，无需处理
		}

		if vf.RefCount <= 0 {
			if err := tx.Where("id = ?", fileID).Delete(&model.VideoFile{}).Error; err != nil {
				return fmt.Errorf("删除VideoFile失败: %w", err)
			}
			// 仅在“正常删除资源”路径（resourceID != 0）下同步清理目录，避免影响创建失败补偿路径。
			if resourceID != 0 && dirName != "" {
				shouldDeleteDir = true
			}
			utils.InfoLog(fmt.Sprintf("【删除VideoFile】fileID=%d, dirName=%s, 引用计数为0", fileID, dirName), "upload")
		} else {
			utils.InfoLog(fmt.Sprintf("【减少引用计数】fileID=%d, 剩余引用=%d", fileID, vf.RefCount), "upload")
		}
		return nil
	}); err != nil {
		utils.ErrorLog(fmt.Sprintf("减少文件引用计数事务失败, fileID=%d", fileID), "upload", err.Error())
		return
	}

	// resource 计数归零后，尝试加速清理本地目录（并在非 local 存储时同步删除 OSS 对应对象）
	if shouldDeleteDir {
		localDir := filepath.Join("./upload/video", dirName)
		var deleteErrors []string
		deleteVideoFromOSS(localDir, &deleteErrors)
		if err := os.RemoveAll(localDir); err != nil && !os.IsNotExist(err) {
			utils.ErrorLog("同步清理本地视频目录失败", "cleanup", err.Error())
		}
	}
}

// VideoMeta 前端传入的视频元数据（OSS 直传时本地无文件，前端截取）
type VideoMeta struct {
	Duration  float64 // 时长（秒）
	Width     int
	Height    int
	CodecName string
}

// CompleteUploadVideo 完成视频上传（支持全局去重）
// fileID: 关联的视频文件ID
// skipTranscode: 如果文件已转码完成，跳过转码直接使用
func CompleteUploadVideo(vid, userId, fileID uint, videoName, title string, skipTranscode bool, meta *VideoMeta, replaceResourceID ...uint) (vo.ResourceResp, error) {
	suffix := utils.GetFileSuffix(title)
	uploadVideoPath := "./upload/video/" + videoName + "/upload" + suffix

	// 去掉后缀名并截断过长标题
	titleWithoutExt := title[:len(title)-len(path.Ext(title))]
	titleWithoutExt = truncateString(titleWithoutExt, 255)

	// 处理替换场景
	replaceID := uint(0)
	sortOrder := -1 // 默认为append到末尾
	if len(replaceResourceID) > 0 && replaceResourceID[0] > 0 {
		replaceID = replaceResourceID[0]
		// 校验旧资源属于同一视频
		var oldResource model.Resource
		if err := global.Mysql.Where("id = ? AND vid = ?", replaceID, vid).First(&oldResource).Error; err != nil {
			return vo.ResourceResp{}, errors.New("被替换的资源不存在")
		}
		sortOrder = oldResource.SortOrder
		// ShortID 不提前转移：审核期间旧资源保留 ShortID，保证弹幕/字幕/历史正常
		// 审核通过时再由 ReviewVideoApproved 统一转移
	} else {
		// append模式：取当前最大排序序号+1
		global.Mysql.Model(&model.Resource{}).Where("vid = ?", vid).
			Select("COALESCE(MAX(sort_order), -1)").Scan(&sortOrder)
		sortOrder++ // sortOrder 现在是 maxOrder+1
	}

	// 分配新资源的 ShortID
	rSid, errSid := AllocateUniqueResourceShortID()
	if errSid != nil {
		return vo.ResourceResp{}, errSid
	}

	// 如果文件已就绪（秒传），直接复用已有转码结果
	if skipTranscode {
		// 查询已有的资源获取视频信息（duration, codecName等）
		var existingResource model.Resource
		if err := global.Mysql.Where("file_id = ?", fileID).First(&existingResource).Error; err == nil {
			// 【重要】秒传的资源也需要审核，状态设为 WAITING_REVIEW
			// 审核员决定是否通过，不能因为别人上传过就自动通过
			resource := model.Resource{
				Vid:            vid,
				Uid:            userId,
				Title:          titleWithoutExt,
				Status:         global.WAITING_REVIEW, // 等待审核，不直接通过
				VisibleStatus:  global.VISIBLE_HIDDEN, // 审核通过前对外隐藏
				Duration:       existingResource.Duration,
				FileID:         fileID,
				SortOrder:      sortOrder,
				ShortID:        rSid,
				CodecName:      existingResource.CodecName,
				ReplaceID:      replaceID,
			}
			if err := global.Mysql.Create(&resource).Error; err != nil {
				return vo.ResourceResp{}, errors.New("保存视频失败")
			}

			// 复制已有资源的 video_index_file 记录到新资源
		copyVideoIndexFiles(existingResource.ID, resource.ID)

		utils.InfoLog(fmt.Sprintf("【秒传成功】vid=%d, fileID=%d, uid=%d, 等待审核", vid, fileID, userId), "upload")
		return vo.ResourceToResourceResp(resource), nil
		}
		// 如果没找到已有资源，说明秒传判断有误，走正常转码流程
		utils.InfoLog(fmt.Sprintf("【秒传异常】fileID=%d 无已有资源记录，转为正常转码", fileID), "upload")
	}

	// 正常流程：读取视频信息并启动转码
	var transcodingInfo *dto.TranscodingInfo
	if meta != nil && meta.Duration > 0 {
		// 前端传入元数据（OSS 直传场景，本地无视频文件）
		transcodingInfo = &dto.TranscodingInfo{
			Width:        meta.Width,
			Height:       meta.Height,
			Duration:     meta.Duration,
			CodecName:    meta.CodecName,
			OriginalVideoStatus: -1,
		}
		utils.InfoLog(fmt.Sprintf("【前端元数据】fileID=%d, %dx%d, %.1fs, codec=%s", fileID, meta.Width, meta.Height, meta.Duration, meta.CodecName), "upload")
	} else {
		var err error
		transcodingInfo, err = ProcessVideoInfo(uploadVideoPath)
		if err != nil {
			return vo.ResourceResp{}, errors.New("读取视频信息失败")
		}
	}

	// 存入数据库
	// 如果视频已公开，新分P对外隐藏，等转码完成后再改为可见
	// 替换场景：新资源排在同一位置，复用旧 ShortID（弹幕/字幕/历史继承），旧资源等新资源成功后再隐藏
	resource := model.Resource{
		Vid:            vid,
		Uid:            userId,
		Title:          titleWithoutExt,
		CodecName:      transcodingInfo.CodecName,
		Status:         global.VIDEO_PROCESSING,
		VisibleStatus:  global.VISIBLE_HIDDEN,
		Duration:       utils.SecFromFloat(transcodingInfo.Duration),
		FileID:         fileID,
		SortOrder:      sortOrder,
		ShortID:        rSid,
		ReplaceID:      replaceID,
	}
	if err := global.Mysql.Create(&resource).Error; err != nil {
		return vo.ResourceResp{}, errors.New("保存视频失败")
	}

	// 启动转码服务
	transcodingInfo.VideoID = vid
	transcodingInfo.DirName = videoName
	transcodingInfo.ResourceID = resource.ID
	transcodingInfo.OutputDir = "./upload/video/" + videoName + "/"
	transcodingInfo.InputFile = transcodingInfo.OutputDir + "upload" + suffix
	transcodingInfo.Suffix = suffix
	if err := GetCurrentTranscoder().Enqueue(context.Background(), transcodingInfo); err != nil {
		utils.ErrorLog("转码入队失败", "upload",
			fmt.Sprintf("ResourceID=%d, err=%v", resource.ID, err))
	}

	return vo.ResourceToResourceResp(resource), nil
}

// copyVideoIndexFiles 将 srcResourceID 的 video_index_file 记录复制给 dstResourceID（秒传场景）
func copyVideoIndexFiles(srcResourceID, dstResourceID uint) {
	var origFiles []model.VideoIndexFile
	if err := global.Mysql.Where("resource_id = ?", srcResourceID).Find(&origFiles).Error; err != nil {
		utils.ErrorLog("复制video_index_file失败", "upload", fmt.Sprintf("src=%d dst=%d err=%s", srcResourceID, dstResourceID, err.Error()))
		return
	}
	for _, f := range origFiles {
		newFile := f
		newFile.Model = gorm.Model{} // 清除 ID/时间，让 GORM 插入新行
		newFile.ResourceID = dstResourceID
		if err := global.Mysql.Create(&newFile).Error; err != nil {
			utils.ErrorLog("复制video_index_file记录失败", "upload", fmt.Sprintf("src=%d dst=%d quality=%s err=%s", srcResourceID, dstResourceID, f.Quality, err.Error()))
		}
	}
}

// 生成文件url（存库用，必须短且不过期）
// 公开OSS: 直接拼接公开URL
// 私有OSS/local: 存 /api/ 路径，访问时后端302到签名URL
func generateFileUrl(objectKey string) string {
	if global.Config.Storage.OssType != "local" && !global.Config.Storage.Private {
		return global.GetOssUrl(objectKey)
	}

	return "/api/" + objectKey
}

// 初始化视频
func initVideo(userId uint, videoPath, title string, coverObjectKey string) (uint, error) {
	var coverUrl string

	if coverObjectKey != "" {
		coverUrl = generateFileUrl(coverObjectKey)
		utils.InfoLog(fmt.Sprintf("initVideo: 使用前端封面 coverObjectKey=%s", coverObjectKey), "upload")
	} else {
		utils.InfoLog("initVideo: 无封面，等待用户后续上传", "upload")
	}

	// 去掉后缀名并截断过长标题
	titleWithoutExt := title[:len(title)-len(path.Ext(title))]
	titleWithoutExt = truncateString(titleWithoutExt, 255)

	utils.InfoLog(fmt.Sprintf("initVideo: uid=%d, title=%s, cover=%s", userId, titleWithoutExt, coverUrl), "upload")

	videoId, err := CreateVideo(&model.Video{
		Uid:       userId,
		Cover:     coverUrl,
		Title:     titleWithoutExt,
		Copyright: global.CopyrightReprint,
		Status:    global.CREATED_VIDEO,
	})
	if err != nil {
		utils.ErrorLog("CreateVideo失败", "upload", fmt.Sprintf("uid=%d, title=%s, err=%v", userId, titleWithoutExt, err))
		return 0, err
	}

	return videoId, nil
}

// 随机生成图片文件名
func generateImgFilename(suffix string) string {
	id := global.SnowflakeNode.Generate()
	return id.String() + suffix
}

// 随机视频文件名
func generateVideoFilename() string {
	id := global.SnowflakeNode.Generate()
	return id.String()
}

// uploadMergedVideoToOSS 确保已合并的本地视频文件已上传到 OSS，
// 供远程 Worker 拉取转码。上传失败仅记日志，不阻塞业务流程。
func uploadMergedVideoToOSS(dirName, suffix, localPath string) {
	if global.Config.Storage.OssType == "local" {
		return
	}
	objectKey := fmt.Sprintf("video/%s/upload%s", dirName, suffix)
	if err := global.Storage.PutObjectFromFile(objectKey, localPath); err != nil {
		utils.ErrorLog("上传视频到OSS失败(秒传场景)", "upload",
			fmt.Sprintf("key=%s, err=%v", objectKey, err))
		return
	}
	utils.InfoLog(fmt.Sprintf("上传视频到OSS成功(秒传场景): key=%s", objectKey), "upload")
	go UploadToBackupWithRetry(objectKey, localPath, "video")
}

// ═══════════════════════════════════════════════════════════════════
// 直传 OSS 服务函数
// ═══════════════════════════════════════════════════════════════════

// presignURLTTL 预签名 URL 有效期。
// 图片直传：15分钟足够；视频分片直传：大文件可能超过15分钟，给24小时兜底
const (
	presignURLTTL        = 15 * time.Minute // 图片直传
	presignVideoChunkTTL = 24 * time.Hour   // 视频分片直传
	// presignBatchSize 每次分批签名的分片数量（防预签名URL被滥用）
	presignBatchSize = 20
)

// PresignImageUpload 为图片直传生成预签名 URL。
// local 模式返回错误（请使用原始上传接口）。
func PresignImageUpload(ctx *gin.Context, req dto.PresignImageReq) (dto.PresignImageResp, error) {
	if global.Config.Storage.OssType == "local" || global.Storage == nil {
		return dto.PresignImageResp{}, errors.New("当前存储模式不支持直传，请使用原始上传接口")
	}

	suffix := path.Ext(req.FileName)
	if !utils.IsImgType(suffix, global.Config.File.AllowedImgExts) {
		return dto.PresignImageResp{}, errors.New("不支持的图片格式")
	}

	fileName := generateImgFilename(suffix)
	objectKey := "image/" + fileName

	presignURL, err := global.Storage.PresignPutObject(objectKey, presignURLTTL)
	if err != nil {
		utils.ErrorLog("图片预签名失败", "upload", err.Error())
		return dto.PresignImageResp{}, errors.New("生成上传地址失败")
	}

	utils.InfoLog(fmt.Sprintf("【直传图片】objectKey=%s, uid=%d", objectKey, ctx.GetUint("userId")), "upload")
	return dto.PresignImageResp{
		PresignURL: presignURL,
		ObjectKey:  objectKey,
	}, nil
}

// ConfirmImageUpload 确认图片已直传到 OSS 并写入数据库。
func ConfirmImageUpload(ctx *gin.Context, req dto.ConfirmImageReq) (string, error) {
	if global.Config.Storage.OssType == "local" || global.Storage == nil {
		return "", errors.New("当前存储模式不支持直传")
	}

	userId := ctx.GetUint("userId")
	fileName := path.Base(req.ObjectKey)

	// 检查是否已存在相同 hash 的图片
	var existingFile model.ImageFile
	global.Mysql.Where("hash = ?", req.Hash).First(&existingFile)
	if existingFile.ID != 0 {
		url := generateFileUrl("image/" + existingFile.FileName)
		cache.SetUploadImage(url, userId)
		return url, nil
	}

	// 写入数据库
	global.Mysql.Create(&model.ImageFile{
		Uid:      userId,
		FileName: fileName,
		Hash:     req.Hash,
	})

	url := generateFileUrl(req.ObjectKey)
	cache.SetUploadImage(url, userId)
	utils.InfoLog(fmt.Sprintf("【直传图片确认】objectKey=%s, uid=%d", req.ObjectKey, userId), "upload")
	return url, nil
}

// InitVideoUpload 发起视频分片直传 OSS（仅 OSS 非 local 模式）。
// 创建 VideoFile 记录 + 发起 OSS 分片上传 + 返回每个分片的预签名 URL。
func InitVideoUpload(ctx *gin.Context, req dto.InitVideoUploadReq) (dto.InitVideoUploadResp, error) {
	if global.Config.Storage.OssType == "local" || global.Storage == nil {
		return dto.InitVideoUploadResp{}, errors.New("当前存储模式不支持直传，请使用原始上传接口")
	}

	userId := ctx.GetUint("userId")

	suffix := path.Ext(req.FileName)
	if !utils.IsVideoType(suffix, global.Config.File.AllowedVideoExts) {
		return dto.InitVideoUploadResp{}, errors.New("不支持的视频格式")
	}

	// 查询或创建文件记录（hash + size 全局唯一）
	var videoFile model.VideoFile
	err := global.Mysql.Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("hash = ? AND size = ?", req.Hash, req.Size).First(&videoFile)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			videoFile = model.VideoFile{
				Hash:         req.Hash,
				Size:         req.Size,
				DirName:      generateVideoFilename(),
				OriginalName: truncateString(req.FileName, 255),
				ChunksCount:  req.TotalChunks,
				Status:       model.FileStatusUploading,
				UploaderUid:  userId,
			}
			if err := tx.Create(&videoFile).Error; err != nil {
				return err
			}
		} else if result.Error != nil {
			return result.Error
		}

		// 校验分片数量一致
		if videoFile.ChunksCount != req.TotalChunks {
			return errors.New("分片数量不匹配")
		}
		return nil
	})
	if err != nil {
		utils.ErrorLog("创建视频文件记录失败", "upload", err.Error())
		return dto.InitVideoUploadResp{}, errors.New("初始化上传失败")
	}

	// 如果已上传完成（秒传），直接返回
	if videoFile.Status == model.FileStatusReady {
		utils.InfoLog(fmt.Sprintf("【直传视频秒传】fileID=%s, hash=%s", videoFile.DirName, req.Hash), "upload")
		return dto.InitVideoUploadResp{
			FileID:      videoFile.DirName,
			TotalChunks: 0, // 0 表示秒传
		}, nil
	}

	objectKey := fmt.Sprintf("video/%s/upload%s", videoFile.DirName, suffix)

	// 发起 OSS 分片上传
	uploadID, err := global.Storage.InitiateMultipartUpload(objectKey)
	if err != nil {
		utils.ErrorLog("发起OSS分片上传失败", "upload", err.Error())
		return dto.InitVideoUploadResp{}, errors.New("初始化上传失败")
	}

	// 缓存 uploadID 到 Redis（续签接口需要），TTL 与分片签名有效期一致
	uploadIDKey := fmt.Sprintf("video:upload_id:%s", videoFile.DirName)
	global.Redis.RawClient().Set(ctx, uploadIDKey, uploadID, presignVideoChunkTTL)

	// 分批签名：只签第一批 presignBatchSize 个分片，前端上传完后续签
	firstBatchEnd := req.TotalChunks
	if firstBatchEnd > presignBatchSize {
		firstBatchEnd = presignBatchSize
	}

	chunks := make([]dto.PresignChunkResp, firstBatchEnd)
	for i := 0; i < firstBatchEnd; i++ {
		partNumber := i + 1
		presignURL, err := global.Storage.PresignUploadPart(uploadID, objectKey, partNumber, presignVideoChunkTTL)
		if err != nil {
			utils.ErrorLog("分片预签名失败", "upload", fmt.Sprintf("part=%d, err=%v", partNumber, err))
			return dto.InitVideoUploadResp{}, errors.New("生成分片上传地址失败")
		}
		chunks[i] = dto.PresignChunkResp{
			Index:      i,
			PartNumber: partNumber,
			PresignURL: presignURL,
		}
	}

	// 计算下一批起始 index
	nextBatchStart := -1 // -1 表示全部签完
	if firstBatchEnd < req.TotalChunks {
		nextBatchStart = firstBatchEnd
	}

	utils.InfoLog(fmt.Sprintf("【直传视频初始化】fileID=%s, uploadID=%s, objectKey=%s, chunks=%d/%d, uid=%d",
		videoFile.DirName, uploadID, objectKey, firstBatchEnd, req.TotalChunks, userId), "upload")

	return dto.InitVideoUploadResp{
		FileID:         videoFile.DirName,
		UploadID:       uploadID,
		ObjectKey:      objectKey,
		TotalChunks:    req.TotalChunks,
		Chunks:         chunks,
		NextBatchStart: nextBatchStart,
	}, nil
}

// PresignUploadChunks 前端上传完一批分片后续签下一批预签名 URL。
// 通过 fileID 查找 uploadID 和 objectKey，为指定范围的分片生成预签名 URL。
func PresignUploadChunks(ctx *gin.Context, req dto.PresignUploadChunksReq) (dto.PresignUploadChunksResp, error) {
	if global.Config.Storage.OssType == "local" || global.Storage == nil {
		return dto.PresignUploadChunksResp{}, errors.New("当前存储模式不支持直传")
	}

	// 限制每次最大请求数量
	if req.Count <= 0 || req.Count > presignBatchSize {
		req.Count = presignBatchSize
	}
	if req.Start < 0 {
		return dto.PresignUploadChunksResp{}, errors.New("无效的起始位置")
	}

	// 查找文件记录
	var videoFile model.VideoFile
	if err := global.Mysql.Where("dir_name = ?", req.FileID).First(&videoFile).Error; err != nil {
		return dto.PresignUploadChunksResp{}, errors.New("文件记录不存在")
	}

	suffix := path.Ext(videoFile.OriginalName)
	objectKey := fmt.Sprintf("video/%s/upload%s", videoFile.DirName, suffix)

	// 需要 uploadID 来签名，从 Redis 获取（InitVideoUpload 时应缓存）
	// 或者重新查找：如果文件已在上传中，uploadID 存在 Redis
	uploadIDKey := fmt.Sprintf("video:upload_id:%s", videoFile.DirName)
	uploadID, err := global.Redis.RawClient().Get(ctx, uploadIDKey).Result()
	if err != nil {
		return dto.PresignUploadChunksResp{}, errors.New("上传会话已过期，请重新发起上传")
	}

	// 计算本批签名范围
	end := req.Start + req.Count
	if end > videoFile.ChunksCount {
		end = videoFile.ChunksCount
	}

	chunks := make([]dto.PresignChunkResp, 0, end-req.Start)
	for i := req.Start; i < end; i++ {
		partNumber := i + 1
		presignURL, err := global.Storage.PresignUploadPart(uploadID, objectKey, partNumber, presignVideoChunkTTL)
		if err != nil {
			utils.ErrorLog("续签分片失败", "upload", fmt.Sprintf("fileID=%s, part=%d, err=%v", req.FileID, partNumber, err))
			return dto.PresignUploadChunksResp{}, errors.New("续签分片地址失败")
		}
		chunks = append(chunks, dto.PresignChunkResp{
			Index:      i,
			PartNumber: partNumber,
			PresignURL: presignURL,
		})
	}

	// 计算下一批起始 index
	nextBatchStart := -1
	if end < videoFile.ChunksCount {
		nextBatchStart = end
	}

	utils.InfoLog(fmt.Sprintf("【直传视频续签】fileID=%s, start=%d, count=%d, next=%d",
		req.FileID, req.Start, len(chunks), nextBatchStart), "upload")

	return dto.PresignUploadChunksResp{
		Chunks:         chunks,
		NextBatchStart: nextBatchStart,
	}, nil
}

// CompleteVideoUpload 完成视频分片直传 OSS。
// 验证分片完整性、完成 OSS 合并、更新文件状态为 Merged。
func CompleteVideoUpload(ctx *gin.Context, req dto.CompleteVideoUploadReq) error {
	if global.Config.Storage.OssType == "local" || global.Storage == nil {
		return errors.New("当前存储模式不支持直传")
	}

	userId := ctx.GetUint("userId")

	var videoFile model.VideoFile
	if err := global.Mysql.Where("dir_name = ?", req.FileID).First(&videoFile).Error; err != nil {
		return errors.New("视频文件不存在")
	}

	// 权限校验：只有上传者可以完成
	if videoFile.UploaderUid != userId {
		return errors.New("无权操作")
	}

	// 校验分片数量
	if len(req.Parts) != videoFile.ChunksCount {
		return fmt.Errorf("分片数量不匹配：期望 %d，实际 %d", videoFile.ChunksCount, len(req.Parts))
	}

	suffix := utils.GetFileSuffix(videoFile.OriginalName)
	objectKey := fmt.Sprintf("video/%s/upload%s", videoFile.DirName, suffix)

	// 转换为 oss.CompletePart
	ossParts := make([]osslib.CompletePart, 0, len(req.Parts))
	for _, p := range req.Parts {
		ossParts = append(ossParts, osslib.CompletePart{
			PartNumber: p.PartNumber,
			ETag:       p.ETag,
		})
	}

	// 完成 OSS 分片上传
	if err := global.Storage.CompleteMultipartUpload(req.UploadID, objectKey, ossParts); err != nil {
		utils.ErrorLog("完成OSS分片上传失败", "upload",
			fmt.Sprintf("fileID=%s, err=%v", req.FileID, err))
		return errors.New("合并分片失败")
	}

	// 更新文件状态为 Merged
	if err := global.Mysql.Model(&videoFile).Update("status", model.FileStatusMerged).Error; err != nil {
		utils.ErrorLog("更新视频文件状态失败", "upload",
			fmt.Sprintf("fileID=%s, err=%v", req.FileID, err))
	}

	utils.InfoLog(fmt.Sprintf("【直传视频完成】fileID=%s, objectKey=%s, uid=%d", req.FileID, objectKey, userId), "upload")
	return nil
}
