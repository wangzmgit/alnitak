package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/resp"

	"github.com/gin-gonic/gin"
)

const (
	rateLimitPrefix = "rate_limit:"
)

// RateLimitConfig 频率限制配置
type RateLimitConfig struct {
	Window time.Duration // 时间窗口
	Limit  int           // 窗口内最大请求数
	KeyFn  func(*gin.Context) string // Redis key 生成函数
}

// 默认按 IP+路径 限流
func ipEndpointKey(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	path := ctx.FullPath()
	return rateLimitPrefix + ip + ":" + path
}

// RateLimiter 返回 Gin 频率限制中间件
func RateLimiter(cfg RateLimitConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := cfg.KeyFn(ctx)
		rdb := global.Redis.RawClient()
		count, err := rdb.Incr(context.Background(), key).Result()
		if err != nil {
			// Redis 异常时放行，避免误杀
			ctx.Next()
			return
		}
		// 第一次访问（count == 1）时设置过期时间
		if count == 1 {
			rdb.Expire(context.Background(), key, cfg.Window)
		}
		if count > int64(cfg.Limit) {
			resp.Result(ctx, http.StatusTooManyRequests, nil, "请求过于频繁，请稍后再试")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// 按用户ID+路径限流（用于已认证的写接口）
func userEndpointKey(ctx *gin.Context) string {
	userID := ctx.GetInt("userId")
	path := ctx.FullPath()
	return rateLimitPrefix + "user:" + itoa(userID) + ":" + path
}

// strconv.Itoa 别名，避免引入新 import
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	return fmt.Sprintf("%d", i)
}

// 预定义敏感接口的限流配置
var (
	// 登录/注册：每 IP 每分钟 10 次
	LoginRateLimit = RateLimitConfig{
		Window: 1 * time.Minute,
		Limit:  10,
		KeyFn:  ipEndpointKey,
	}

	// 发送验证码：每 IP 每分钟 3 次（防止短信/邮件轰炸）
	EmailCodeRateLimit = RateLimitConfig{
		Window: 1 * time.Minute,
		Limit:  3,
		KeyFn:  ipEndpointKey,
	}

	// 修改密码：每 IP 每分钟 5 次
	ModifyPwdRateLimit = RateLimitConfig{
		Window: 1 * time.Minute,
		Limit:  5,
		KeyFn:  ipEndpointKey,
	}

	// 弹幕：每用户每分钟 30 条
	DanmakuRateLimit = RateLimitConfig{
		Window: 1 * time.Minute,
		Limit:  30,
		KeyFn:  userEndpointKey,
	}

	// 评论：每用户每分钟 5 条
	CommentRateLimit = RateLimitConfig{
		Window: 1 * time.Minute,
		Limit:  5,
		KeyFn:  userEndpointKey,
	}

	// 私信：每用户每分钟 10 条
	WhisperRateLimit = RateLimitConfig{
		Window: 1 * time.Minute,
		Limit:  10,
		KeyFn:  userEndpointKey,
	}

	// 创建视频/新增分P（写数据库+OSS）：每用户每分钟 3 次
	UploadCreateRateLimit = RateLimitConfig{
		Window: 1 * time.Minute,
		Limit:  3,
		KeyFn:  userEndpointKey,
	}

	// 上传视频分片+合并：每用户每分钟 60 次
	UploadChunkRateLimit = RateLimitConfig{
		Window: 1 * time.Minute,
		Limit:  60,
		KeyFn:  userEndpointKey,
	}

	// 上传图片：每用户每分钟 10 次
	UploadImgRateLimit = RateLimitConfig{
		Window: 1 * time.Minute,
		Limit:  10,
		KeyFn:  userEndpointKey,
	}

	// 查询分片（checkVideo）：每用户每分钟 30 次
	UploadCheckRateLimit = RateLimitConfig{
		Window: 1 * time.Minute,
		Limit:  30,
		KeyFn:  userEndpointKey,
	}
)

// IsSensitiveAuthPath 判断是否属于敏感认证路径（与 operation.go 脱敏逻辑保持一致）
func IsSensitiveAuthPath(path string) bool {
	sensitivePaths := []string{
		"/api/v1/auth/register",
		"/api/v1/auth/login",
		"/api/v1/auth/emailLogin",
		"/api/v1/auth/modifyPwd",
		"/api/v1/email/sendEmailCode",
	}
	for _, p := range sensitivePaths {
		if strings.EqualFold(path, p) {
			return true
		}
	}
	return false
}
