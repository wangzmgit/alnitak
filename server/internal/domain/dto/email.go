package dto

type SendEmailReq struct {
	// 邮箱
	Email string
	// 验证ID
	CaptchaId string
}
