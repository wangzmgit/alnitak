package global

import (
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
	"interastral-peace.com/alnitak/internal/config"
	"interastral-peace.com/alnitak/pkg/casbin"
	"interastral-peace.com/alnitak/pkg/oss"
	"interastral-peace.com/alnitak/pkg/redis"
)

var (
	Config            *config.Config
	Mysql             *gorm.DB
	Redis             *redis.Redis
	Casbin            *casbin.Casbin
	Storage           oss.Storage
	SnowflakeNode     *snowflake.Node
	VideoPartitionMap map[uint]uint
)
