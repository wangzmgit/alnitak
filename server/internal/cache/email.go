package cache

import (
	"strconv"

	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

// 获取邮箱验证码
func GetEmailCode(email string) string {
	return global.Redis.Get(EMAIL_CODE_KEY + email)
}

// 保存邮箱验证码
func SetEmailCode(email string, code string) {
	global.Redis.Set(EMAIL_CODE_KEY+email, code, EMIAL_CODE_EXPIRATION_TIME)
}

// 删除邮箱验证码
func DelEmailCode(email string) {
	global.Redis.Del(EMAIL_CODE_KEY + email)
}

// 获取邮箱验证码发送冷却剩余时间（秒），0表示无冷却
func GetEmailCodeCooldown(email string) int {
	ttl := global.Redis.TTL(EMAIL_CODE_COOLDOWN_KEY + email)
	if ttl > 0 {
		return int(ttl.Seconds())
	}
	return 0
}

// 设置邮箱验证码发送冷却
func SetEmailCodeCooldown(email string) {
	global.Redis.Set(EMAIL_CODE_COOLDOWN_KEY+email, "1", EMAIL_CODE_COOLDOWN_TIME)
}

// ========== 验证码防暴力枚举 ==========

const EMAIL_CODE_TRY_COUNT_KEY = "email_code_try_count_key:"

// 验证码尝试次数过期时间（与验证码本身一致，验证码过期后计数自动清空）
const EMAIL_CODE_TRY_COUNT_EXPIRATION = EMIAL_CODE_EXPIRATION_TIME

func GetEmailCodeTryCount(email string) int {
	s := global.Redis.Get(EMAIL_CODE_TRY_COUNT_KEY + email)
	if s == "" {
		return 0
	}
	count, err := strconv.Atoi(s)
	if err != nil {
		utils.ErrorLog("邮箱验证码尝试次数转换int失败", "cache", err.Error())
	}
	return count
}

func IncrEmailCodeTryCount(email string) {
	global.Redis.Incr(EMAIL_CODE_TRY_COUNT_KEY + email)
	global.Redis.Expire(EMAIL_CODE_TRY_COUNT_KEY+email, EMAIL_CODE_TRY_COUNT_EXPIRATION)
}

func DelEmailCodeTryCount(email string) {
	global.Redis.Del(EMAIL_CODE_TRY_COUNT_KEY + email)
}
