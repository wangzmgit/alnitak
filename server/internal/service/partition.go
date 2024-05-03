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

func GetPartitionList() (partitions []vo.PartitionResp) {
	partitions = cache.GetPartition()
	if len(partitions) == 0 {
		global.Mysql.Select("id,name,parent_id").Model(&model.Partition{}).Scan(&partitions)

		// 存入缓存
		cache.SetPartition(partitions)
	}

	return
}

func AddPartition(ctx *gin.Context, addPartitionReq dto.AddPartitionReq) error {
	if addPartitionReq.ParentId != 0 && !IsParentPartitionExist(addPartitionReq.ParentId) {
		return errors.New("所属分区不存在")
	}

	// 保存到数据库
	partition := model.Partition{
		Name:     addPartitionReq.Name,
		ParentId: addPartitionReq.ParentId,
	}

	if err := global.Mysql.Create(&partition).Error; err != nil {
		utils.ErrorLog("创建分区失败", "partition", err.Error())
		return errors.New("创建分区失败")
	}

	global.PartitionMap = GetPartitionMap()

	// 更新缓存
	var partitions []vo.PartitionResp
	global.Mysql.Select("id,name,parent_id").Model(&model.Partition{}).Scan(&partitions)
	cache.SetPartition(partitions)

	return nil
}

// 删除分区
func DeletePartition(ctx *gin.Context, id uint) error {
	currentPartition, err := FindPartitionById(id)
	if err != nil {
		return errors.New("获取分区信息失败")
	}

	if currentPartition.ParentId == 0 { // 判断是否存在二级分区
		var partition model.Partition
		global.Mysql.Where("parent_id = ?", currentPartition.ID).First(&partition)
		if partition.ID != 0 {
			return errors.New("当前分区下存在二级分区")
		}
	} else { // 判断分区下是否存在视频
		var video model.Video
		global.Mysql.Where("partition_id = ?", currentPartition.ID).First(&video)
		if video.ID != 0 {
			return errors.New("当前分区下存在视频")
		}
	}

	if err := global.Mysql.Where("id = ?", id).Delete(&model.Partition{}).Error; err != nil {
		utils.ErrorLog("删除分区失败", "partition", err.Error())
		return errors.New("删除分区失败")
	}

	global.PartitionMap = GetPartitionMap()

	// 更新缓存
	var partitions []vo.PartitionResp
	global.Mysql.Select("id,name,parent_id").Model(&model.Partition{}).Scan(&partitions)
	cache.SetPartition(partitions)

	return nil
}

// 所属分区是否存在
func IsParentPartitionExist(id uint) bool {
	var partition model.Partition
	global.Mysql.First(&partition, id)
	if partition.ID != 0 && partition.ParentId == 0 {
		return true
	}
	return false
}

// 是否为子分区
func IsSubpartition(id uint) bool {
	var partition model.Partition
	global.Mysql.First(&partition, id)
	if partition.ID != 0 && partition.ParentId != 0 {
		return true
	}
	return false
}

func GetPartitionMap() map[uint]uint {
	var partitions []model.Partition
	global.Mysql.Where("parent_id != 0 ").Find(&partitions)

	// 生成map，key为parentId不为0的id，value为parentId
	partitionMap := make(map[uint]uint)
	for _, partition := range partitions {
		partitionMap[partition.ID] = partition.ParentId
	}

	return partitionMap
}

// 通过分区ID查询分区
func FindPartitionById(id uint) (partition model.Partition, err error) {
	err = global.Mysql.Where("`id` = ?", id).First(&partition).Error
	return
}
