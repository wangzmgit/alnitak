package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func Collect(ctx *gin.Context, addCollectReq dto.AddCollectReq) error {
	videoId := addCollectReq.Vid
	video := GetVideoInfo(videoId)
	if video.ID == 0 {
		return errors.New("视频不存在")
	}

	//处理收藏部分
	userId := ctx.GetUint("userId")

	// 取用户要添加的收藏夹id数组和已经存在收藏夹id数组的差集
	ids := utils.DifferenceSet(addCollectReq.AddList, GetVideoCollection(userId, videoId))

	if length := len(ids); length > 0 {
		//处理要写入的数据
		newCollects := make([]model.CollectVideo, length)
		for i := 0; i < length; i++ {
			newCollects[i] = model.CollectVideo{
				Uid:          userId,
				Vid:          videoId,
				CollectionID: ids[i],
			}
		}

		if err := global.Mysql.Create(&newCollects).Error; err != nil {
			utils.ErrorLog("收藏失败", "collect", err.Error())
			return errors.New("收藏失败")
		}
	}

	//处理取消收藏部分
	if len(addCollectReq.CancelList) > 0 {
		if err := global.Mysql.Where("uid = ? and vid = ? and collection_id in ?", userId, videoId, addCollectReq.CancelList).
			Delete(&model.CollectVideo{}).Error; err != nil {
			utils.ErrorLog("收藏失败", "collect", err.Error())
			return errors.New("收藏失败")
		}
	}

	return nil
}

func HasCollect(ctx *gin.Context, videoId uint) (bool, error) {
	userId := ctx.GetUint("userId")
	collect, err := FindCollectByUid(videoId, userId)
	if err != nil {
		utils.ErrorLog("获取收藏信息失败", "collect", err.Error())
		return false, errors.New("获取失败")
	}

	return (collect.ID != 0), nil
}

func GetCollectedInfo(ctx *gin.Context, videoId uint) []uint {
	userId := ctx.GetUint("userId")

	return GetVideoCollection(userId, videoId)
}

func GetCollectVideo(ctx *gin.Context, collectionId uint, page, pageSize int) (total int64, videos []vo.VideoResp, err error) {
	var collection model.Collection
	err = global.Mysql.First(&collection, collectionId).Error
	if err != nil {
		utils.ErrorLog("查询收藏夹失败", "collect", err.Error())
		return total, videos, errors.New("无法收藏夹评论信息")
	}

	// uid为收藏夹所有者
	userId := ctx.GetUint("userId")
	if collection.ID == 0 || (!collection.Open && collection.Uid != userId) {
		return total, videos, errors.New("收藏夹不存在")
	}

	var videoIds []uint
	global.Mysql.Model(&model.CollectVideo{}).Where("collection_id = ? and uid = ?", collectionId, userId).Count(&total)
	global.Mysql.Model(&model.CollectVideo{}).Where("collection_id = ? and uid = ?", collectionId, userId).
		Limit(pageSize).Offset((page-1)*pageSize).Pluck("vid", &videoIds)
	if err := global.Mysql.Model(&model.Video{}).Select(vo.VIDEO_FIELD).
		Where("id in (?)", videoIds).Scan(&videos).Error; err != nil {
		utils.ErrorLog("获取收藏视频失败", "collect", err.Error())
		return total, videos, errors.New("获取失败")
	}

	for i := 0; i < len(videos); i++ {
		videos[i].Author = GetUserBaseInfo(videos[i].Uid)
		videos[i].Clicks += GetVideoClicks(videos[i].ID)
	}

	return
}

func FindCollectByUid(videoId, userId uint) (collect model.CollectVideo, err error) {
	err = global.Mysql.Where("`uid` = ? and vid = ?", userId, videoId).First(&collect).Error

	return
}

// 获取视频所在的收藏夹
func GetVideoCollection(userId, videoId uint) (id []uint) {
	global.Mysql.Model(&model.CollectVideo{}).Where("uid = ? and vid = ?", userId, videoId).Pluck("collection_id", &id)

	return
}
