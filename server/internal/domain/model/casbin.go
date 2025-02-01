package model

// 角色权限规则
type CasbinRule struct {
	ID    uint   `gorm:"primarykey"`
	Ptype string `gorm:"type:varchar(100);"`
	V0    string `gorm:"type:varchar(100);"`
	V1    string `gorm:"type:varchar(100);"`
	V2    string `gorm:"type:varchar(100);"`
	V3    string `gorm:"type:varchar(100);"`
	V4    string `gorm:"type:varchar(100);"`
	V5    string `gorm:"type:varchar(100);"`
}

func (table *CasbinRule) TableName() string {
	return "casbin_rule"
}
