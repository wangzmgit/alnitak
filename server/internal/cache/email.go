package cache

import (
	"interastral-peace.com/alnitak/internal/global"
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
