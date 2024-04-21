package service

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func UploadVideoInfo(ctx *gin.Context, uploadVideoReq dto.UploadVideoReq) error {
	userId := ctx.GetUint("userId")
	if cache.GetUploadImage(uploadVideoReq.Cover) != userId {
		// 查询是否与旧封面图一致
		if v, _ := FindVideoById(uploadVideoReq.Vid); v.Cover != uploadVideoReq.Cover {
			return errors.New("文件链接无效")
		}
	}

	if !IsSubpartition(uploadVideoReq.PartitionId) {
		return errors.New("分区不存在")
	}

	if err := global.Mysql.Model(&model.Video{}).Where("id = ?", uploadVideoReq.Vid).Updates(
		map[string]interface{}{
			"title":        uploadVideoReq.Title,
			"cover":        uploadVideoReq.Cover,
			"desc":         uploadVideoReq.Desc,
			"tags":         uploadVideoReq.Tags,
			"copyright":    uploadVideoReq.Copyright,
			"partition_id": uploadVideoReq.PartitionId,
			"status":       getVideoStatus(uploadVideoReq.Vid),
		},
	).Error; err != nil {
		utils.ErrorLog("修改视频失败", "video", err.Error())
		return errors.New("修改失败")
	}

	return nil
}

func GetVideoStatus(ctx *gin.Context, vid uint) (video vo.VideoStatusResp, err error) {
	userId := ctx.GetUint("userId")
	global.Mysql.Model(&model.Video{}).Select(vo.VIDEO_STATUS_FIELD).Where("id = ? and uid = ?", vid, userId).Scan(&video)
	if video.ID == 0 {
		return video, errors.New("视频不存在")
	}

	//查询分区下的视频资源
	video.Resources = GetReviewResourceList(vid)

	return video, nil
}

// 获取视频文件
func GetVideoFile(ctx *gin.Context, resourceId uint, quality string) (string, error) {
	var file model.VideoIndexFile
	global.Mysql.Where("resource_id = ? and quality = ?", resourceId, quality).First(&file)

	res := ""
	key := uuid.New().String()
	cache.SetVideoSlice(key, file.DirName)
	for _, line := range strings.Split(file.Content, "\n") {
		//根据关键词覆盖当前行
		if strings.Contains(line, ".ts") {
			res += "/api/v1/video/slice/" + line + "?key=" + key + "\n"
		} else {
			res += line + "\n"
		}
	}

	return res, nil
}

// 获取视频切所在文件目录
func GetVideoSliceDir(key string) string {
	return cache.GetVideoSlice(key)
}

// 获取自己上传的视频
func GetUploadVideoList(ctx *gin.Context, page, pageSize int) (total int64, videos []vo.UploadVideoResp) {
	userId := ctx.GetUint("userId")

	global.Mysql.Model(&model.Video{}).Where("uid = ?", userId).Count(&total)
	global.Mysql.Model(&model.Video{}).Select(vo.UPLOAD_VIDEO_FIELD).
		Where("uid = ?", userId).Limit(pageSize).Offset((page - 1) * pageSize).Scan(&videos)

	// 更新播放量数据
	for i := 0; i < len(videos); i++ {
		videos[i].Clicks = GetVideoClicks(videos[i].ID)
	}

	return
}

func EditVideoInfo(ctx *gin.Context, editVideoReq dto.EditVideoReq) error {
	userId := ctx.GetUint("userId")
	if cache.GetUploadImage(editVideoReq.Cover) != userId {
		// 查询是否与旧封面图一致
		if v, _ := FindVideoById(editVideoReq.Vid); v.Cover != editVideoReq.Cover {
			return errors.New("文件链接无效")
		}
	}

	if err := global.Mysql.Model(&model.Video{}).Where("id = ?", editVideoReq.Vid).Updates(
		map[string]interface{}{
			"title":  editVideoReq.Title,
			"cover":  editVideoReq.Cover,
			"desc":   editVideoReq.Desc,
			"tags":   editVideoReq.Tags,
			"status": getVideoStatus(editVideoReq.Vid),
		},
	).Error; err != nil {
		utils.ErrorLog("修改视频失败", "video", err.Error())
		return errors.New("修改失败")
	}

	// 删除视频信息缓存
	cache.DelVideoInfo(editVideoReq.Vid)

	return nil
}

func DeleteVideo(ctx *gin.Context, id uint) error {
	var video model.Video
	userId := ctx.GetUint("userId")
	global.Mysql.Model(&model.Video{}).Where("id = ? and uid = ?", id, userId).First(&video)
	if video.ID == 0 {
		return errors.New("视频不存在")
	}

	if err := global.Mysql.Where("id = ?", id).Delete(&model.Video{}).Error; err != nil {
		utils.ErrorLog("删除视频失败", "video", err.Error())
		return errors.New("删除视频失败")
	}

	// 删除视频信息缓存
	cache.DelVideoInfo(id)

	return nil
}

// 获取视频信息
func GetVideoById(ctx *gin.Context, videoId uint) (vo.VideoResp, error) {
	video := GetVideoInfo(videoId)
	if video.ID == 0 {
		return video, errors.New("视频信息不存在")
	}

	// 增加播放量(一个ip在同一个视频下，每30分钟可重新增加1播放量)
	AddVideoClicks(videoId, ctx.ClientIP())
	video.Clicks = GetVideoClicks(video.ID)

	return video, nil
}

// 获取所有的视频列表
func GetAllVideoList(ctx *gin.Context) (videos []vo.AllVideoResp) {
	userId := ctx.GetUint("userId")
	global.Mysql.Model(&model.Video{}).Select("`id`,`title`").Where("uid = ?", userId).Scan(&videos)

	return
}

// 获取用户视频
func GetVideoByUser(ctx *gin.Context, userId uint, page, pageSize int) (total int64, videos []vo.UploadVideoResp) {
	global.Mysql.Model(&model.Video{}).
		Where("uid = ? and status = ?", userId, global.AUDIT_APPROVED).Count(&total)
	global.Mysql.Model(&model.Video{}).Select(vo.UPLOAD_VIDEO_FIELD).
		Where("uid = ? and status = ?", userId, global.AUDIT_APPROVED).
		Limit(pageSize).Offset((page - 1) * pageSize).Scan(&videos)

	// 更新播放量数据
	for i := 0; i < len(videos); i++ {
		videos[i].Clicks = GetVideoClicks(videos[i].ID)
	}

	return
}

// 获取视频列表(后台管理)
func GetVideoListManage(videoListReq dto.VideoListReq) (total int64, videos []vo.VideoInfoManageResp) {
	global.Mysql.Model(&model.Video{}).Where("status = ?", global.AUDIT_APPROVED).Count(&total)
	global.Mysql.Model(&model.Video{}).Where("status = ?", global.AUDIT_APPROVED).
		Limit(videoListReq.PageSize).Offset((videoListReq.Page - 1) * videoListReq.PageSize).Scan(&videos)

	// 更新播放量和作者数据
	for i := 0; i < len(videos); i++ {
		videos[i].Clicks = GetVideoClicks(videos[i].ID)
		videos[i].Author = GetUserBaseInfo(videos[i].Uid)
	}

	return
}

// 删除视频(后台管理)
func DeleteVideoManage(ctx *gin.Context, id uint) error {
	if err := global.Mysql.Where("id = ?", id).Delete(&model.Video{}).Error; err != nil {
		utils.ErrorLog("删除视频失败", "video", err.Error())
		return errors.New("删除视频失败")
	}

	// 删除视频信息缓存
	cache.DelVideoInfo(id)

	return nil
}

// 获取待审核视频列表
func GetReviewList(reviewListReq dto.ReviewListReq) (total int64, videos []vo.ReviewListResp) {
	global.Mysql.Model(&model.Video{}).Where("status = ?", global.WAITING_REVIEW).Count(&total)
	global.Mysql.Model(&model.Video{}).Where("status = ?", global.WAITING_REVIEW).
		Limit(reviewListReq.PageSize).Offset((reviewListReq.Page - 1) * reviewListReq.PageSize).Scan(&videos)

	// 更新播放量和作者数据
	for i := 0; i < len(videos); i++ {
		videos[i].Author = GetUserBaseInfo(videos[i].Uid)
	}

	return
}

func CreateVideo(video *model.Video) (uint, error) {
	if err := global.Mysql.Create(video).Error; err != nil {
		utils.ErrorLog("创建视频失败", "video", err.Error())
		return 0, errors.New("创建视频失败")
	}

	return video.ID, nil
}

// 通过视频ID查询视频
func FindVideoById(id uint) (video model.Video, err error) {
	err = global.Mysql.Where("`id` = ?", id).First(&video).Error
	return
}

// 获取视频信息
func GetVideoInfo(videoId uint) (video vo.VideoResp) {
	video = cache.GetVideoInfo(videoId)
	if video.ID == 0 {
		global.Mysql.Model(&model.Video{}).Select(vo.VIDEO_FIELD).
			Where("id = ? and status = ?", videoId, global.AUDIT_APPROVED).Scan(&video)
		if video.ID == 0 {
			utils.ErrorLog("视频信息不存在", "video", "")
			return
		}

		// 获取作者信息
		video.Author = GetUserBaseInfo(video.Uid)
		// 获取视频资源
		video.Resources = GetVideoResourceByStatus(videoId, global.AUDIT_APPROVED)

		// 存到redis
		cache.SetVideoInfo(video)
	}

	return
}

// 获取视频状态
func getVideoStatus(videoId uint) int {
	var count int64
	global.Mysql.Model(&model.Resource{}).Where("vid = ? and status = ?", videoId, global.VIDEO_PROCESSING).Count(&count)
	// 如果没有转码中的视频，则更新视频为待审核
	if count == 0 {
		return global.WAITING_REVIEW
	}

	return global.SUBMIT_REVIEW
}
