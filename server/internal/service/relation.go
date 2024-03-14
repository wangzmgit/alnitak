package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func Follow(ctx *gin.Context, targetUid uint) error {
	//判断关注的是否为自己
	userId := ctx.GetUint("userId")
	if targetUid == userId {
		return errors.New("不能关注自己")
	}

	relation := FindUserRelation(userId, targetUid)
	if relation.Relation == global.FOLLOWED || relation.Relation == global.MUTUAL_FANS { // 已经关注了
		return nil
	}

	// 开始事务
	tx := global.Mysql.Begin()

	// 查询目标关系
	relation.Relation = global.FOLLOWED
	targetRelation := FindUserRelation(targetUid, userId)
	if targetRelation.Relation == global.FOLLOWED {
		// 目标用户已经关注了自己，更新状态为互粉
		targetRelation.Relation = global.MUTUAL_FANS
		if err := tx.Save(&targetRelation).Error; err != nil {
			tx.Rollback()
			return errors.New("关注失败")
		}

		relation.Relation = global.MUTUAL_FANS
	}

	if relation.ID == 0 { // 之前为关注过插入一条数据
		relation.Uid = userId
		relation.TargetUid = targetUid
		if err := tx.Create(&relation).Error; err != nil {
			tx.Rollback()
			utils.ErrorLog("关注失败", "relation", err.Error())
			return errors.New("关注失败")
		}
	} else { // 之前关注过更新状态
		if err := tx.Save(&relation).Error; err != nil {
			tx.Rollback()
			utils.ErrorLog("关注失败", "relation", err.Error())
			return errors.New("关注失败")
		}
	}

	tx.Commit()

	return nil
}

func UnFollow(ctx *gin.Context, targetUid uint) error {
	//判断关注的是否为自己
	userId := ctx.GetUint("userId")
	if targetUid == userId {
		return errors.New("不能关注自己")
	}

	relation := FindUserRelation(userId, targetUid)
	if relation.ID == 0 || relation.Relation == global.NOT_FOLLOWING { // 当前没有关注
		return nil
	}

	// 开始事务
	tx := global.Mysql.Begin()

	if relation.Relation == global.MUTUAL_FANS {
		// 更新目标用户的互相关注为已关注
		if err := tx.Model(&model.Relation{}).Where("uid = ? and target_uid = ?", targetUid, userId).Updates(
			map[string]interface{}{"relation": global.FOLLOWED},
		).Error; err != nil {
			tx.Rollback()
			utils.ErrorLog("取关失败", "relation", err.Error())
			return errors.New("取关失败")
		}
	}

	relation.Relation = global.NOT_FOLLOWING
	if err := tx.Save(&relation).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("取关失败", "relation", err.Error())
		return errors.New("取关失败")
	}

	return nil
}

func GetUserRelation(ctx *gin.Context, targetUid uint) int {
	userId := ctx.GetUint("userId")

	relation := FindUserRelation(userId, targetUid)

	return relation.Relation
}

// 查询用户关系
func FindUserRelation(userId, targetUid uint) model.Relation {
	var relation model.Relation
	global.Mysql.Where("uid = ? and target_uid = ?", userId, targetUid).First(&relation)

	return relation
}
