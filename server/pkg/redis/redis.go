package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"interastral-peace.com/alnitak/internal/config"
)

type Redis struct {
	redisClient *redis.Client
	ctx         context.Context
}

func Init(c config.Redis) *Redis {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password: c.Password,
		DB:       0, // use default DB
	})
	zap.L().Info("redis连接成功", zap.String("module", "db"))

	return &Redis{
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}

func (r *Redis) Set(key string, value interface{}, expiration time.Duration) {
	if err := r.redisClient.Set(r.ctx, key, value, expiration).Err(); err != nil {
		zap.L().Error("Redis Set 操作失败", zap.String("key", key), zap.Error(err))
	}
}

func (r *Redis) Get(key string) string {
	val, err := r.redisClient.Get(r.ctx, key).Result()
	if err != nil && err != redis.Nil {
		zap.L().Error("Redis Get 操作失败", zap.String("key", key), zap.Error(err))
		return ""
	}
	return val
}

func (r *Redis) Del(key string) {
	if err := r.redisClient.Del(r.ctx, key).Err(); err != nil {
		zap.L().Error("Redis Del 操作失败", zap.String("key", key), zap.Error(err))
	}
}

func (r *Redis) Incr(key string) {
	if err := r.redisClient.Incr(r.ctx, key).Err(); err != nil {
		zap.L().Error("Redis Incr 操作失败", zap.String("key", key), zap.Error(err))
	}
}

// SetNX 原子性地"仅当 key 不存在时设置并带过期时间"，返回 true 表示成功占位，常用于限流/防刷。
// 底层 Redis 错误时返回 false，宁可放行也不误伤用户。
func (r *Redis) SetNX(key string, value interface{}, expiration time.Duration) bool {
	ok, err := r.redisClient.SetNX(r.ctx, key, value, expiration).Result()
	if err != nil {
		zap.L().Error("Redis SetNX 操作失败", zap.String("key", key), zap.Error(err))
		return false
	}
	return ok
}

func (r *Redis) Keys(key string) []string {
	val, err := r.redisClient.Keys(r.ctx, key).Result()
	if err != nil {
		zap.L().Error("Redis Keys 操作失败", zap.String("pattern", key), zap.Error(err))
		return []string{}
	}
	return val
}

func (r *Redis) Expire(key string, expiration time.Duration) bool {
	val, err := r.redisClient.Expire(r.ctx, key, expiration).Result()
	if err != nil {
		zap.L().Error("Redis Expire 操作失败", zap.String("key", key), zap.Error(err))
		return false
	}
	return val
}

func (r *Redis) TTL(key string) time.Duration {
	val, err := r.redisClient.TTL(r.ctx, key).Result()
	if err != nil {
		zap.L().Error("Redis TTL 操作失败", zap.String("key", key), zap.Error(err))
		return -1
	}
	return val
}

// 向有序集合插入数据
func (r *Redis) ZAdd(key string, score float64, member interface{}) {
	if err := r.redisClient.ZAdd(r.ctx, key, redis.Z{Score: score, Member: member}).Err(); err != nil {
		zap.L().Error("Redis ZAdd 操作失败", zap.String("key", key), zap.Error(err))
	}
}

// 有序集合中的数量
func (r *Redis) ZCard(key string) int64 {
	val, err := r.redisClient.ZCard(r.ctx, key).Result()
	if err != nil {
		zap.L().Error("Redis ZCard 操作失败", zap.String("key", key), zap.Error(err))
		return 0
	}
	return val
}

// 有序集中成员的分数值
func (r *Redis) ZScore(key string, member string) float64 {
	val, err := r.redisClient.ZScore(r.ctx, key, member).Result()
	if err != nil && err != redis.Nil {
		zap.L().Error("Redis ZScore 操作失败", zap.String("key", key), zap.Error(err))
		return 0
	}
	return val
}

// 移除有序集中指定排名区间内的所有成员
func (r *Redis) ZRemRangeByRank(key string, start, stop int64) {
	if err := r.redisClient.ZRemRangeByRank(r.ctx, key, start, stop).Err(); err != nil {
		zap.L().Error("Redis ZRemRangeByRank 操作失败", zap.String("key", key), zap.Error(err))
	}
}

// 向有序集合插入数据
func (r *Redis) ZRem(key string, member ...interface{}) {
	if err := r.redisClient.ZRem(r.ctx, key, member...).Err(); err != nil {
		zap.L().Error("Redis ZRem 操作失败", zap.String("key", key), zap.Error(err))
	}
}

// 向集合插入数据
func (r *Redis) SAdd(key string, member interface{}) {
	if err := r.redisClient.SAdd(r.ctx, key, member).Err(); err != nil {
		zap.L().Error("Redis SAdd 操作失败", zap.String("key", key), zap.Error(err))
	}
}

// 向集合移除数据
func (r *Redis) SRem(key string, member interface{}) {
	if err := r.redisClient.SRem(r.ctx, key, member).Err(); err != nil {
		zap.L().Error("Redis SRem 操作失败", zap.String("key", key), zap.Error(err))
	}
}

// 随机从集合中获取n个数据
func (r *Redis) SRandMemberN(key string, count int64) []string {
	val, err := r.redisClient.SRandMemberN(r.ctx, key, count).Result()
	if err != nil {
		zap.L().Error("Redis SRandMemberN 操作失败", zap.String("key", key), zap.Error(err))
		return []string{}
	}
	return val
}

// 从集合中获取数据
func (r *Redis) SMembers(key string) []string {
	val, err := r.redisClient.SMembers(r.ctx, key).Result()
	if err != nil {
		zap.L().Error("Redis SMembers 操作失败", zap.String("key", key), zap.Error(err))
		return []string{}
	}
	return val
}

// 清空当前数据库的所有key
func (r *Redis) FlushDB() {
	if err := r.redisClient.FlushDB(r.ctx).Err(); err != nil {
		zap.L().Error("Redis FlushDB 操作失败", zap.Error(err))
	}
}
