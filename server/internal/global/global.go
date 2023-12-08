package global

import (
	"gorm.io/gorm"
	"interastral-peace.com/alnitak/pkg/casbin"
	"interastral-peace.com/alnitak/pkg/redis"
)

var (
	Mysql  *gorm.DB
	Redis  *redis.Redis
	Casbin *casbin.Casbin
)
