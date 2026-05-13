package model

import "gorm.io/gorm"

// 认证类型配置表（可后台管理）
type AuthType struct {
	gorm.Model
	Code        string `gorm:"type:varchar(32);comment:认证类型code;not null;uniqueIndex"`
	Name        string `gorm:"type:varchar(32);comment:显示名称;not null"`
	Category    string `gorm:"type:varchar(32);comment:分类:creator-创作者,identity-身份,enterprise-企业,other-其他;not null;index"`
	Desc        string `gorm:"type:varchar(255);comment:描述"`
	Icon        string `gorm:"type:varchar(255);comment:图标URL"`
	Color       string `gorm:"type:varchar(16);comment:颜色代码，如#FB7299"`
	Priority    int    `gorm:"default:0;comment:显示顺序，越大越靠前"`
	IsEnable    bool   `gorm:"default:true;comment:是否启用"`
	IsDefault   bool   `gorm:"default:false;comment:是否默认认证类型"`
	RequiredFields string `gorm:"type:text;comment:必填字段JSON，如:[\"name\",\"id_card\"]"`
}

func (table *AuthType) TableName() string {
	return "auth_type"
}

// 认证分类常量
const (
	AuthCategoryCreator   = "creator"   // 创作者
	AuthCategoryIdentity = "identity"  // 身份
	AuthCategoryEnterprise = "enterprise" // 企业
	AuthCategoryOther    = "other"     // 其他
)