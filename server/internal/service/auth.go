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

// ========== 认证类型 Service ==========

// AddAuthType 添加认证类型
func AddAuthType(ctx *gin.Context, req dto.AddAuthTypeReq) error {
	// 检查 code 是否已存在
	var count int64
	global.Mysql.Model(&model.AuthType{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		return errors.New("认证类型code已存在")
	}

	authType := model.AuthType{
		Code:            req.Code,
		Name:            req.Name,
		Category:        req.Category,
		Desc:            req.Desc,
		Icon:            req.Icon,
		Color:           req.Color,
		Priority:        req.Priority,
		IsEnable:        req.IsEnable,
		RequiredFields: req.RequiredFields,
	}

	if err := global.Mysql.Create(&authType).Error; err != nil {
		utils.ErrorLog("添加认证类型失败", "auth", err.Error())
		return errors.New("添加失败")
	}

	return nil
}

// EditAuthType 编辑认证类型
func EditAuthType(ctx *gin.Context, req dto.EditAuthTypeReq) error {
	var authType model.AuthType
	if err := global.Mysql.First(&authType, req.ID).Error; err != nil {
		return errors.New("认证类型不存在")
	}

	// 检查 code 是否与其他重复
	var count int64
	global.Mysql.Model(&model.AuthType{}).Where("code = ? AND id != ?", req.Code, req.ID).Count(&count)
	if count > 0 {
		return errors.New("认证类型code已存在")
	}

	authType.Code = req.Code
	authType.Name = req.Name
	authType.Category = req.Category
	authType.Desc = req.Desc
	authType.Icon = req.Icon
	authType.Color = req.Color
	authType.Priority = req.Priority
	authType.IsEnable = req.IsEnable
	authType.RequiredFields = req.RequiredFields

	if err := global.Mysql.Save(&authType).Error; err != nil {
		utils.ErrorLog("编辑认证类型失败", "auth", err.Error())
		return errors.New("编辑失败")
	}

	return nil
}

// DeleteAuthType 删除认证类型
func DeleteAuthType(ctx *gin.Context, id uint) error {
	// 检查是否有用户正在使用该认证类型
	var count int64
	global.Mysql.Model(&model.UserAuth{}).Where("auth_type_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("该认证类型已有用户使用，无法删除")
	}

	if err := global.Mysql.Delete(&model.AuthType{}, id).Error; err != nil {
		utils.ErrorLog("删除认证类型失败", "auth", err.Error())
		return errors.New("删除失败")
	}

	return nil
}

// GetAuthTypeList 获取认证类型列表
func GetAuthTypeList(category string) (list []vo.AuthTypeListResp) {
	query := global.Mysql.Model(&model.AuthType{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	query = query.Where("is_enable = ?", true).Order("priority desc, id asc")

	query.Scan(&list)
	return
}

// GetAllAuthTypeList 获取所有认证类型列表（管理用）
func GetAllAuthTypeList(page, pageSize int) (total int64, list []vo.AuthTypeResp) {
	global.Mysql.Model(&model.AuthType{}).Count(&total)
	global.Mysql.Model(&model.AuthType{}).Order("priority desc, id asc").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Scan(&list)
	return
}

// GetAuthTypeByID 根据ID获取认证类型
func GetAuthTypeByID(id uint) (*model.AuthType, error) {
	var authType model.AuthType
	if err := global.Mysql.First(&authType, id).Error; err != nil {
		return nil, errors.New("认证类型不存在")
	}
	return &authType, nil
}

// GetAuthTypeByCode 根据code获取认证类型
func GetAuthTypeByCode(code string) (*model.AuthType, error) {
	var authType model.AuthType
	if err := global.Mysql.Where("code = ? AND is_enable = ?", code, true).First(&authType).Error; err != nil {
		return nil, errors.New("认证类型不存在")
	}
	return &authType, nil
}

// ========== 用户认证 Service ==========

// AddUserAuth 添加用户认证
func AddUserAuth(ctx *gin.Context, req dto.AddUserAuthReq) error {
	// 验证认证类型存在且启用
	authType, err := GetAuthTypeByID(req.AuthTypeID)
	if err != nil {
		return errors.New("认证类型不存在")
	}
	if !authType.IsEnable {
		return errors.New("该认证类型已禁用")
	}

	// 检查用户是否已有该类型认证
	var count int64
	global.Mysql.Model(&model.UserAuth{}).Where("uid = ? AND auth_type_id = ?", req.Uid, req.AuthTypeID).Count(&count)
	if count > 0 {
		return errors.New("用户已有该类型认证")
	}

	userAuth := model.UserAuth{
		Uid:          req.Uid,
		AuthTypeID:   req.AuthTypeID,
		AuthTypeCode: authType.Code,
		Title:        req.Title,
		Desc:         req.Desc,
		ExtraData:    req.ExtraData,
	}

	if err := global.Mysql.Create(&userAuth).Error; err != nil {
		utils.ErrorLog("添加用户认证失败", "auth", err.Error())
		return errors.New("添加失败")
	}

	return nil
}

// EditUserAuth 编辑用户认证
func EditUserAuth(ctx *gin.Context, req dto.EditUserAuthReq) error {
	var userAuth model.UserAuth
	if err := global.Mysql.First(&userAuth, req.ID).Error; err != nil {
		return errors.New("用户认证记录不存在")
	}

	// 验证认证类型存在且启用
	authType, err := GetAuthTypeByID(req.AuthTypeID)
	if err != nil {
		return errors.New("认证类型不存在")
	}
	if !authType.IsEnable {
		return errors.New("该认证类型已禁用")
	}

	userAuth.AuthTypeID = req.AuthTypeID
	userAuth.AuthTypeCode = authType.Code
	userAuth.Title = req.Title
	userAuth.Desc = req.Desc
	userAuth.ExtraData = req.ExtraData

	if err := global.Mysql.Save(&userAuth).Error; err != nil {
		utils.ErrorLog("编辑用户认证失败", "auth", err.Error())
		return errors.New("编辑失败")
	}

	return nil
}

// DeleteUserAuth 删除用户认证
func DeleteUserAuth(ctx *gin.Context, req dto.DeleteUserAuthReq) error {
	var userAuth model.UserAuth
	if err := global.Mysql.Where("id = ? AND uid = ?", req.ID, req.Uid).First(&userAuth).Error; err != nil {
		return errors.New("用户认证记录不存在")
	}

	if err := global.Mysql.Delete(&userAuth).Error; err != nil {
		utils.ErrorLog("删除用户认证失败", "auth", err.Error())
		return errors.New("删除失败")
	}

	return nil
}

// GetUserAuthList 获取用户认证列表（公开，不含 ExtraData）
func GetUserAuthList(uid uint) (list []vo.PublicUserAuthResp) {
	global.Mysql.Model(&model.UserAuth{}).
		Select("user_auth.id, user_auth.uid, user_auth.auth_type_id, user_auth.auth_type_code, user_auth.title, user_auth.desc, auth_type.name as auth_type_name, auth_type.icon as auth_type_icon, auth_type.color as auth_type_color, user_auth.created_at, user_auth.updated_at").
		Joins("LEFT JOIN auth_type ON user_auth.auth_type_id = auth_type.id").
		Where("user_auth.uid = ?", uid).
		Scan(&list)
	return
}

// GetUserAuthListWithUser 获取用户认证列表（带用户信息，管理用）
func GetUserAuthListWithUser(page, pageSize int, authTypeID uint) (total int64, list []vo.UserAuthWithUserResp) {
	query := global.Mysql.Model(&model.UserAuth{}).Joins("LEFT JOIN auth_type ON user_auth.auth_type_id = auth_type.id").Joins("LEFT JOIN user ON user_auth.uid = user.id")

	if authTypeID > 0 {
		query = query.Where("user_auth.auth_type_id = ?", authTypeID)
	}

	query.Count(&total)
	query.Select("user_auth.id, user_auth.uid, user.username, user.nickname, user.avatar, user_auth.auth_type_id, user_auth.auth_type_code, auth_type.name as auth_type_name, auth_type.icon as auth_type_icon, auth_type.color as auth_type_color, user_auth.title, user_auth.desc, user_auth.extra_data, user_auth.created_at, user_auth.updated_at").
		Order("user_auth.created_at desc").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Scan(&list)
	return
}

// GetUserAuthByID 根据ID获取用户认证
func GetUserAuthByID(id uint) (*vo.UserAuthResp, error) {
	var userAuth vo.UserAuthResp
	err := global.Mysql.Model(&model.UserAuth{}).
		Select("user_auth.id, user_auth.uid, user_auth.auth_type_id, user_auth.auth_type_code, user_auth.title, user_auth.desc, user_auth.extra_data, auth_type.name as auth_type_name, auth_type.icon as auth_type_icon, auth_type.color as auth_type_color, user_auth.created_at, user_auth.updated_at").
		Joins("LEFT JOIN auth_type ON user_auth.auth_type_id = auth_type.id").
		Where("user_auth.id = ?", id).First(&userAuth).Error
	if err != nil {
		return nil, errors.New("用户认证记录不存在")
	}
	return &userAuth, nil
}

// GetUserAuthByUid 根据用户ID获取认证信息（公开，不含 ExtraData）
func GetUserAuthByUid(uid uint) ([]vo.PublicUserAuthResp, error) {
	var list []vo.PublicUserAuthResp
	err := global.Mysql.Model(&model.UserAuth{}).
		Select("user_auth.id, user_auth.uid, user_auth.auth_type_id, user_auth.auth_type_code, user_auth.title, user_auth.desc, auth_type.name as auth_type_name, auth_type.icon as auth_type_icon, auth_type.color as auth_type_color, user_auth.created_at, user_auth.updated_at").
		Joins("LEFT JOIN auth_type ON user_auth.auth_type_id = auth_type.id").
		Where("user_auth.uid = ?", uid).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetUserPrimaryAuth 获取用户主要认证（用于展示，公开不含 ExtraData）
func GetUserPrimaryAuth(uid uint) *vo.PublicUserAuthResp {
	list := GetUserAuthList(uid)
	if len(list) > 0 {
		return &list[0]
	}
	return nil
}

// InitDefaultAuthTypes 初始化默认认证类型
func InitDefaultAuthTypes() {
	defaultTypes := []model.AuthType{
		{Code: "up", Name: "UP主认证", Category: "creator", Desc: "知名创作者", Icon: "", Color: "#FB7299", Priority: 100, IsEnable: true},
		{Code: "musician", Name: "音乐人认证", Category: "creator", Desc: "音乐领域创作者", Icon: "", Color: "#23ADE5", Priority: 90, IsEnable: true},
		{Code: "artist", Name: "画师认证", Category: "creator", Desc: "绘画领域创作者", Icon: "", Color: "#F25F8D", Priority: 90, IsEnable: true},
		{Code: "vlog", Name: "Vlog认证", Category: "creator", Desc: "Vlog领域创作者", Icon: "", Color: "#FF6B6B", Priority: 90, IsEnable: true},
		{Code: "enterprise", Name: "企业认证", Category: "enterprise", Desc: "企业机构认证", Icon: "", Color: "#FFD700", Priority: 80, IsEnable: true},
		{Code: "media", Name: "媒体认证", Category: "other", Desc: "媒体机构认证", Icon: "", Color: "#4A90D9", Priority: 70, IsEnable: true},
		{Code: "government", Name: "政府认证", Category: "other", Desc: "政府机构认证", Icon: "", Color: "#7CB342", Priority: 70, IsEnable: true},
	}

	for _, t := range defaultTypes {
		var count int64
		global.Mysql.Model(&model.AuthType{}).Where("code = ?", t.Code).Count(&count)
		if count == 0 {
			global.Mysql.Create(&t)
		}
	}
}