package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

type UserAuthApi struct{}

// 注册
func (u *UserAuthApi) Register(ctx *gin.Context) {
	// 获取参数
	var registerReq dto.RegisterReq
	if err := ctx.Bind(&registerReq); err != nil {
		service.AddFailOperation(ctx, "注册", "请求参数有误", nil, nil)
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if !utils.VerifyEmail(registerReq.Email) {
		service.AddFailOperation(ctx, "注册", "邮箱格式错误", gin.H{"邮箱": registerReq.Email}, nil)
		resp.FailWithMessage(ctx, "邮箱格式错误")
		return
	}

	if utils.VerifyStringLength(registerReq.Password, "<", 6) {
		service.AddFailOperation(ctx, "注册", "密码长度小于6位", gin.H{"邮箱": registerReq.Email}, nil)
		resp.FailWithMessage(ctx, "密码长度不能小于6位")
		return
	}

	if !utils.VerifyStringLength(registerReq.Code, "=", 6) {
		service.AddFailOperation(ctx, "注册", "验证码长度不为6位", gin.H{"邮箱": registerReq.Email}, nil)
		resp.FailWithMessage(ctx, "验证码长度为6位")
		return
	}

	if err := service.UserRegister(ctx, registerReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	service.AddOkOperation(ctx, "注册", "注册成功", gin.H{"邮箱": registerReq.Email})
	resp.Ok(ctx)
}

// 登录
func (u *UserAuthApi) Login(ctx *gin.Context) {
	// 获取参数
	var loginReq dto.LoginReq
	if err := ctx.Bind(&loginReq); err != nil {
		service.AddFailOperation(ctx, "登录", "请求参数有误", nil, nil)
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if !utils.VerifyEmail(loginReq.Email) {
		service.AddFailOperation(ctx, "登录", "邮箱格式错误", gin.H{"邮箱": loginReq.Email}, nil)
		resp.FailWithMessage(ctx, "邮箱格式错误")
		return
	}

	if utils.VerifyStringLength(loginReq.Password, "<", 6) {
		service.AddFailOperation(ctx, "登录", "密码长度小于6位", gin.H{"邮箱": loginReq.Email}, nil)
		resp.FailWithMessage(ctx, "密码长度不能小于6位")
		return
	}

	// 如果人机验证通过，删除登录尝试次数
	if cache.GetCaptchaStatus(loginReq.CaptchaId) == global.CAPTCHA_STATUS_PASS {
		// 删除登录尝试次数
		cache.DelLoginTryCount(loginReq.Email)
		// 删除人机验证状态
		cache.DelCaptchaStatus(loginReq.CaptchaId)
	}

	// 读取登录尝试次数，超过3次进行滑块验证
	loginTryCount := cache.GetLoginTryCount(loginReq.Email)
	if loginTryCount >= 3 {
		captchaId := cache.CreateCaptchaStatus()

		service.AddFailOperation(ctx, "登录", "需要人机验证", gin.H{"邮箱": loginReq.Email, "滑块ID": captchaId}, nil)
		resp.Result(ctx, -1, gin.H{"captchaId": captchaId}, "需要人机验证")
		return
	}

	accessToken, refreshToken, err := service.UserLogin(ctx, loginReq)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	service.AddOkOperation(ctx, "登录", "登录成功", nil)
	resp.OkWithData(ctx, gin.H{"token": accessToken, "refreshToken": refreshToken})
}

// 刷新token
func (u *UserAuthApi) UpdateToken(ctx *gin.Context) {
	var tokenReq dto.TokenReq
	if err := ctx.Bind(&tokenReq); err != nil {
		service.AddFailOperation(ctx, "刷新TOKEN", "请求参数有误", nil, nil)
		return
	}

	accessToken, refreshToken, err := service.UpdateToken(ctx, tokenReq)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"token": accessToken, "refreshToken": refreshToken})
}
