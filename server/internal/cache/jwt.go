package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

const ACCESS_TOKEN_BLACKLIST_KEY = "access_token_blacklist:"

func IsRefreshTokenExist(userId uint, token string) bool {
	return global.Redis.ZScore(REFRESH_TOKEN_KEY+utils.UintToString(userId), token) != 0
}

func SetRefreshToken(id uint, token string, ttl ...time.Duration) {
	key := REFRESH_TOKEN_KEY + utils.UintToString(id)
	if global.Redis.ZCard(key) >= MAX_LOGIN_LIMIT {
		// 移除最旧的，保留最近 MAX_LOGIN_LIMIT - 1 个（加上即将添加的新 token 共 MAX_LOGIN_LIMIT 个）
		global.Redis.ZRemRangeByRank(key, 0, -MAX_LOGIN_LIMIT)
	}

	duration := REFRESH_TOKEN_EXPIRATION_TIME
	if len(ttl) > 0 && ttl[0] > 0 {
		duration = ttl[0]
	}
	global.Redis.ZAdd(key, float64(time.Now().Add(duration).Unix()), token)
}

func DelRefreshToken(id uint, token string) {
	global.Redis.ZRem(REFRESH_TOKEN_KEY+utils.UintToString(id), token)
}

// DelAllRefreshToken 删除用户的所有 refreshToken（改密码/强制下线时使用）
func DelAllRefreshToken(id uint) {
	global.Redis.Del(REFRESH_TOKEN_KEY + utils.UintToString(id))
}

// ========== AccessToken 黑名单（登出后立即失效） ==========

// SetBlacklistedAccessToken 将 accessToken 加入黑名单，TTL 为 token 剩余有效期
func SetBlacklistedAccessToken(tokenString string, expiresAt time.Time) {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return
	}
	// 用 token 的 SHA256 摘要做 key，避免原始 token 过长
	hash := sha256.Sum256([]byte(tokenString))
	key := ACCESS_TOKEN_BLACKLIST_KEY + hex.EncodeToString(hash[:])
	global.Redis.Set(key, "1", ttl)
}

// IsAccessTokenBlacklisted 检查 accessToken 是否已被吊销
func IsAccessTokenBlacklisted(tokenString string) bool {
	hash := sha256.Sum256([]byte(tokenString))
	key := ACCESS_TOKEN_BLACKLIST_KEY + hex.EncodeToString(hash[:])
	return global.Redis.Get(key) != ""
}
