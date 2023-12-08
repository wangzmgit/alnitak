package dto

type TokenReq struct {
	RefreshToken string
}

type LoginReq struct {
	// 用户名
	Email string
	// 密码
	Password string
	// 验证ID
	CaptchaId string
}

type EmailLoginReq struct {
	// 邮箱
	Email string
	// 验证码
	Code string
	// 验证ID
	CaptchaId string
}

type RegisterReq struct {
	// 邮箱
	Email string
	// 密码
	Password string
	// 邮箱验证码
	Code string
	// 验证ID
	CaptchaId string
}
