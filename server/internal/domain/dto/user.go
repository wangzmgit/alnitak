package dto

type TokenReq struct {
	RefreshToken string
}

type UserListReq struct {
	Page     int
	PageSize int
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

// 编辑用户信息
type EditUserInfoReq struct {
	Avatar     string
	Name       string
	Gender     int
	Birthday   string
	Sign       string
	SpaceCover string
}

// 修改用户密码验证
type ModifyPwdCheckReq struct {
	Email     string
	CaptchaId string
}

// 修改用户密码
type ModifyPwdReq struct {
	Email     string
	Password  string
	Code      string
	CaptchaId string
}

// 编辑用户信息
type EditUserInfoManageReq struct {
	Uid        uint
	Avatar     string
	SpaceCover string
	Name       string
	Email      string
	Sign       string
}
