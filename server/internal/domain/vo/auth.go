package vo

import "time"

// ========== 认证类型 VO ==========

// AuthTypeResp 认证类型响应
type AuthTypeResp struct {
	ID             uint      `json:"id"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	Category       string    `json:"category"`
	Desc           string    `json:"desc"`
	Icon           string    `json:"icon"`
	Color          string    `json:"color"`
	Priority       int       `json:"priority"`
	IsEnable       bool      `json:"isEnable"`
	RequiredFields string    `json:"requiredFields"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// AuthTypeListResp 认证类型列表响应
type AuthTypeListResp struct {
	ID             uint   `json:"id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	Category       string `json:"category"`
	Desc           string `json:"desc"`
	Icon           string `json:"icon"`
	Color          string `json:"color"`
	Priority       int    `json:"priority"`
	IsEnable       bool   `json:"isEnable"`
	RequiredFields string `json:"requiredFields"`
}

// ========== 用户认证 VO ==========

// UserAuthResp 用户认证响应（管理用，含 ExtraData）
type UserAuthResp struct {
	ID             uint      `json:"id"`
	Uid            uint      `json:"uid"`
	AuthTypeID     uint      `json:"authTypeId"`
	AuthTypeCode   string    `json:"authTypeCode"`
	Title          string    `json:"title"`
	Desc           string    `json:"desc"`
	ExtraData      string    `json:"extraData"`
	AuthTypeName   string    `json:"authTypeName"`   // 认证类型名称
	AuthTypeIcon   string    `json:"authTypeIcon"`   // 认证类型图标
	AuthTypeColor  string    `json:"authTypeColor"`  // 认证类型颜色
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// PublicUserAuthResp 公开接口用户认证响应（不含 ExtraData，避免敏感材料泄露）
type PublicUserAuthResp struct {
	ID            uint      `json:"id"`
	Uid           uint      `json:"uid"`
	AuthTypeID    uint      `json:"authTypeId"`
	AuthTypeCode  string    `json:"authTypeCode"`
	Title         string    `json:"title"`
	Desc          string    `json:"desc"`
	AuthTypeName  string    `json:"authTypeName"`
	AuthTypeIcon  string    `json:"authTypeIcon"`
	AuthTypeColor string    `json:"authTypeColor"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// UserAuthWithUserResp 用户认证（带用户信息）
type UserAuthWithUserResp struct {
	ID             uint      `json:"id"`
	Uid            uint      `json:"uid"`
	Username       string    `json:"username"`       // 用户名
	Nickname       string    `json:"nickname"`       // 昵称
	Avatar         string    `json:"avatar"`         // 头像
	AuthTypeID     uint      `json:"authTypeId"`
	AuthTypeCode   string    `json:"authTypeCode"`
	AuthTypeName   string    `json:"authTypeName"`
	AuthTypeIcon   string    `json:"authTypeIcon"`
	AuthTypeColor  string    `json:"authTypeColor"`
	Title          string    `json:"title"`
	Desc           string    `json:"desc"`
	ExtraData      string    `json:"extraData"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}