package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Redis struct {
	redisClient *redis.Client
	ctx         context.Context
}

func Init() *Redis {
	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")
	password := viper.GetString("redis.password")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0, // use default DB
	})
	zap.L().Info("redis连接成功", zap.String("module", "db"))

	return &Redis{
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}

func (r *Redis) Set(key string, value interface{}, expiration time.Duration) {
	r.redisClient.Set(r.ctx, key, value, expiration)
}

func (r *Redis) Get(key string) string {
	return r.redisClient.Get(r.ctx, key).Val()
}

func (r *Redis) Del(key string) {
	r.redisClient.Del(r.ctx, key)
}

func (r *Redis) Incr(key string) {
	r.redisClient.Incr(r.ctx, key)
}

func (r *Redis) Keys(key string) []string {
	return r.redisClient.Keys(r.ctx, key).Val()
}

func (r *Redis) Expire(key string, expiration time.Duration) bool {
	return r.redisClient.Expire(r.ctx, key, expiration).Val()
}

func (r *Redis) TTL(key string) time.Duration {
	return r.redisClient.TTL(r.ctx, key).Val()
}

// 向有序集合插入数据
func (r *Redis) ZAdd(key string, score float64, member interface{}) {
	r.redisClient.ZAdd(r.ctx, key, redis.Z{Score: score, Member: member})
}

// 有序集合中的数量
func (r *Redis) ZCard(key string) int64 {
	return r.redisClient.ZCard(r.ctx, key).Val()
}

// 有序集中成员的分数值
func (r *Redis) ZScore(key string, member string) float64 {
	return r.redisClient.ZScore(r.ctx, key, member).Val()
}

// 移除有序集中指定排名区间内的所有成员
func (r *Redis) ZRemRangeByRank(key string, start, stop int64) {
	r.redisClient.ZRemRangeByRank(r.ctx, key, start, stop)
}

// 向有序集合插入数据
func (r *Redis) ZRem(key string, member ...interface{}) {
	r.redisClient.ZRem(r.ctx, key, member...)
}

// 向集合插入数据
func (r *Redis) SAdd(key string, member interface{}) {
	r.redisClient.SAdd(r.ctx, key, member)
}

// 向集合移除数据
func (r *Redis) SRem(key string, member interface{}) {
	r.redisClient.SRem(r.ctx, key, member)
}

// 随机从集合中获取n个数据
func (r *Redis) SRandMemberN(key string, count int64) []string {
	return r.redisClient.SRandMemberN(r.ctx, key, count).Val()
}
