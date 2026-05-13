package dto

// ========== 认证类型相关 DTO ==========

// AddAuthTypeReq 添加认证类型
type AddAuthTypeReq struct {
	Code          string `json:"code" form:"code" binding:"required,max=32"`
	Name          string `json:"name" form:"name" binding:"required,max=32"`
	Category      string `json:"category" form:"category" binding:"required,max=32"`
	Desc          string `json:"desc" form:"desc" binding:"max=255"`
	Icon          string `json:"icon" form:"icon" binding:"max=255"`
	Color         string `json:"color" form:"color" binding:"max=16"`
	Priority      int    `json:"priority" form:"priority"`
	IsEnable      bool   `json:"isEnable"`
	RequiredFields string `json:"requiredFields" form:"requiredFields"`
}

// EditAuthTypeReq 编辑认证类型
type EditAuthTypeReq struct {
	ID            uint   `json:"id" form:"id" binding:"required"`
	Code          string `json:"code" form:"code" binding:"required,max=32"`
	Name          string `json:"name" form:"name" binding:"required,max=32"`
	Category      string `json:"category" form:"category" binding:"required,max=32"`
	Desc          string `json:"desc" form:"desc" binding:"max=255"`
	Icon          string `json:"icon" form:"icon" binding:"max=255"`
	Color         string `json:"color" form:"color" binding:"max=16"`
	Priority      int    `json:"priority" form:"priority"`
	IsEnable      bool   `json:"isEnable"`
	RequiredFields string `json:"requiredFields" form:"requiredFields"`
}

// ========== 用户认证相关 DTO ==========

// AddUserAuthReq 添加用户认证
type AddUserAuthReq struct {
	Uid        uint   `json:"uid" form:"uid" binding:"required"`
	AuthTypeID uint   `json:"authTypeId" form:"authTypeId" binding:"required"`
	Title      string `json:"title" form:"title" binding:"required,max=64"`
	Desc       string `json:"desc" form:"desc" binding:"max=255"`
	ExtraData  string `json:"extraData" form:"extraData"`
}

// EditUserAuthReq 编辑用户认证
type EditUserAuthReq struct {
	ID          uint   `json:"id" form:"id" binding:"required"`
	AuthTypeID  uint   `json:"authTypeId" form:"authTypeId" binding:"required"`
	Title       string `json:"title" form:"title" binding:"required,max=64"`
	Desc        string `json:"desc" form:"desc" binding:"max=255"`
	ExtraData   string `json:"extraData" form:"extraData"`
}

// DeleteUserAuthReq 删除用户认证
type DeleteUserAuthReq struct {
	ID  uint `json:"id" form:"id" binding:"required"`
	Uid uint `json:"uid" form:"uid" binding:"required"`
}