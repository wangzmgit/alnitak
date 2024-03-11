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

func AddCollection(ctx *gin.Context, addCollectionReq dto.AddCollectionReq) error {
	userId := ctx.GetUint("userId")
	// 保存到数据库
	if err := global.Mysql.Create(&model.Collection{
		Uid:  userId,
		Name: addCollectionReq.Name,
	}).Error; err != nil {
		utils.ErrorLog("新增收藏夹更新数据库失败", "collection", err.Error())
		return errors.New("新增收藏夹失败")
	}
	return nil
}

func EditCollection(ctx *gin.Context, editCollectionReq dto.EditCollectionReq) error {
	userId := ctx.GetUint("userId")
	if cache.GetUploadImage(editCollectionReq.Cover) != userId {
		return errors.New("文件链接无效")
	}

	collection, err := FindCollectionById(editCollectionReq.ID)
	if err != nil {
		utils.ErrorLog("查询收藏夹失败", "collection", err.Error())
		return errors.New("获取收藏夹失败")
	}

	if collection.ID == 0 || collection.Uid != userId {
		return errors.New("收藏夹不存在")
	}

	if err := global.Mysql.Model(model.Collection{}).Where("id = ?", editCollectionReq.ID).Updates(
		map[string]interface{}{
			"name":  editCollectionReq.Name,
			"cover": editCollectionReq.Cover,
			"desc":  editCollectionReq.Desc,
			"open":  editCollectionReq.Open,
		},
	).Error; err != nil {
		utils.ErrorLog("收藏夹更新数据库失败", "collection", err.Error())
		return errors.New("编辑失败")
	}

	return nil
}

// 删除收藏夹
func DeleteCollection(ctx *gin.Context, id uint) error {
	collection, err := FindCollectionById(id)
	if err != nil {
		utils.ErrorLog("查询收藏夹失败", "collection", err.Error())
		return errors.New("无法收藏夹评论信息")
	}

	// uid为收藏夹所有者
	userId := ctx.GetUint("userId")
	if collection.ID == 0 || collection.Uid != userId {
		return errors.New("收藏夹不存在")
	}

	if err := global.Mysql.Where("id = ?", id).Delete(&model.Collection{}).Error; err != nil {
		utils.ErrorLog("删除收藏夹失败", "collection", err.Error())
		return errors.New("删除失败")
	}

	return nil
}

// 获取收藏夹列表
func GetCollectionList(ctx *gin.Context, page, pageSize int) (total int64, collections []vo.CollectionListResp) {
	userId := ctx.GetUint("userId")

	global.Mysql.Model(&model.Collection{}).Where("uid = ?", userId).Count(&total)
	global.Mysql.Model(&model.Collection{}).Select(vo.COLLECTION_LIST_FIELD).
		Where("uid = ?", userId).Limit(pageSize).Offset((page - 1) * pageSize).Scan(&collections)

	return
}

// 获取收藏夹信息
func GetCollectionInfo(ctx *gin.Context, id uint) (collection vo.CollectionResp, err error) {
	err = global.Mysql.Model(&model.Collection{}).Select(vo.COLLECTION_FIELD).
		Where("id = ?", id).Scan(&collection).Error
	if err != nil {
		utils.ErrorLog("查询收藏夹失败", "collection", err.Error())
		return collection, errors.New("无法收藏夹评论信息")
	}

	// uid为收藏夹所有者
	userId := ctx.GetUint("userId")
	if collection.ID == 0 || (!collection.Open && collection.Uid != userId) {
		return collection, errors.New("收藏夹不存在")
	}

	// 获取收藏夹所有者信息
	collection.Author = GetUserBaseInfo(collection.Uid)

	return
}

func FindCollectionById(id uint) (collection model.Collection, err error) {
	err = global.Mysql.First(&collection, id).Error

	return
}
