package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func UploadVideoInfo(ctx *gin.Context, uploadVideoReq dto.UploadVideoReq) (uint, error) {

	userId := ctx.GetUint("userId")
	if cache.GetUploadImage(uploadVideoReq.Cover) != userId {
		return 0, errors.New("文件链接无效")
	}

	if !IsSubpartition(uploadVideoReq.PartitionId) {
		return 0, errors.New("分区不存在")
	}

	// 保存到数据库
	video := model.Video{
		Uid:         userId,
		Title:       uploadVideoReq.Title,
		Cover:       uploadVideoReq.Cover,
		Desc:        uploadVideoReq.Desc,
		Copyright:   uploadVideoReq.Copyright,
		PartitionId: uploadVideoReq.PartitionId,
		Status:      global.CREATED_VIDEO,
	}

	if err := global.Mysql.Create(&video).Error; err != nil {
		utils.ErrorLog("创建视频失败", "video", err.Error())
		return 0, errors.New("创建视频失败")
	}

	return video.ID, nil
}

func GetVideoStatus(ctx *gin.Context, vid uint) (video vo.VideoStatusResp, err error) {
	userId := ctx.GetUint("userId")
	global.Mysql.Model(&model.Video{}).Select(vo.VIDEO_STATUS_FIELD).Where("id = ? and uid = ?", vid, userId).Scan(&video)
	if video.ID == 0 {
		return video, errors.New("视频不存在")
	}

	//TODO: 查询分区下的视频资源

	return video, nil
}

// 通过视频ID查询视频
func FindVideoById(id uint) (video model.Video, err error) {
	err = global.Mysql.Where("`id` = ?", id).First(&video).Error
	return
}
