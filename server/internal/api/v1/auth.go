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

// 注册
func Register(ctx *gin.Context) {
	// 获取参数
	var registerReq dto.RegisterReq
	if err := ctx.Bind(&registerReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if !utils.VerifyEmail(registerReq.Email) {
		resp.FailWithMessage(ctx, "邮箱格式错误")
		return
	}

	if utils.VerifyStringLength(registerReq.Password, "<", 6) {
		resp.FailWithMessage(ctx, "密码长度不能小于6位")
		return
	}

	if !utils.VerifyStringLength(registerReq.Code, "=", 6) {
		resp.FailWithMessage(ctx, "验证码长度为6位")
		return
	}

	if err := service.UserRegister(ctx, registerReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 登录
func Login(ctx *gin.Context) {
	// 获取参数
	var loginReq dto.LoginReq
	if err := ctx.Bind(&loginReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if !utils.VerifyEmail(loginReq.Email) {
		resp.FailWithMessage(ctx, "邮箱格式错误")
		return
	}

	if utils.VerifyStringLength(loginReq.Password, "<", 6) {
		resp.FailWithMessage(ctx, "密码长度不能小于6位")
		return
	}
	// 如果人机验证通过，删除登录尝试次数
	if utils.VerifyStringLength(loginReq.CaptchaId, ">", 0) {
		if cache.GetCaptchaStatus(loginReq.CaptchaId) == global.CAPTCHA_STATUS_PASS {
			// 删除登录尝试次数
			cache.DelLoginTryCount(loginReq.Email)
			// 删除人机验证状态
			cache.DelCaptchaStatus(loginReq.CaptchaId)
		}
	}

	// 读取登录尝试次数，超过3次进行滑块验证
	loginTryCount := cache.GetLoginTryCount(loginReq.Email)
	if loginTryCount >= 3 {
		captchaId := cache.CreateCaptchaStatus()

		resp.Result(ctx, -1, gin.H{"captchaId": captchaId}, "需要人机验证")
		return
	}

	accessToken, refreshToken, userId, err := service.UserLogin(ctx, loginReq)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"token": accessToken, "refreshToken": refreshToken, "userId": userId})
}

// 邮箱登录
func EmailLogin(ctx *gin.Context) {
	// 获取参数
	var loginReq dto.EmailLoginReq
	if err := ctx.Bind(&loginReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if !utils.VerifyEmail(loginReq.Email) {
		resp.FailWithMessage(ctx, "邮箱格式错误")
		return
	}

	// 如果人机验证通过，删除登录尝试次数
	if utils.VerifyStringLength(loginReq.CaptchaId, ">", 0) {
		if cache.GetCaptchaStatus(loginReq.CaptchaId) == global.CAPTCHA_STATUS_PASS {
			// 删除登录尝试次数
			cache.DelLoginTryCount(loginReq.Email)
			// 删除人机验证状态
			cache.DelCaptchaStatus(loginReq.CaptchaId)
		}
	}

	// 读取登录尝试次数，超过3次进行滑块验证
	loginTryCount := cache.GetLoginTryCount(loginReq.Email)
	if loginTryCount >= 3 {
		captchaId := cache.CreateCaptchaStatus()

		resp.Result(ctx, -1, gin.H{"captchaId": captchaId}, "需要人机验证")
		return
	}

	accessToken, refreshToken, userId, err := service.EmailLogin(ctx, loginReq)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"token": accessToken, "refreshToken": refreshToken, "userId": userId})
}

// 刷新token
func UpdateToken(ctx *gin.Context) {
	var tokenReq dto.TokenReq
	if err := ctx.Bind(&tokenReq); err != nil {
		return
	}

	accessToken, refreshToken, userId, err := service.UpdateToken(ctx, tokenReq)
	if err != nil {
		resp.Result(ctx, 2000, nil, err.Error())
		return
	}

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"token": accessToken, "refreshToken": refreshToken, "userId": userId})
}

// 退出登录
func Logout(ctx *gin.Context) {
	var tokenReq dto.TokenReq
	if err := ctx.Bind(&tokenReq); err != nil {
		return
	}

	service.Logout(ctx, tokenReq)

	// 返回给前端
	resp.Ok(ctx)
}

// 重置密码检查
func ResetPwdCheck(ctx *gin.Context) {
	// 获取参数
	var modifyPwdCheckReq dto.ModifyPwdCheckReq
	if err := ctx.Bind(&modifyPwdCheckReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if !utils.VerifyEmail(modifyPwdCheckReq.Email) {
		resp.FailWithMessage(ctx, "邮箱格式错误")
		return
	}

	if utils.VerifyStringLength(modifyPwdCheckReq.CaptchaId, "=", 0) {
		captchaId := cache.CreateCaptchaStatus()
		resp.Result(ctx, -1, gin.H{"captchaId": captchaId}, "需要人机验证")
		return
	}

	switch cache.GetCaptchaStatus(modifyPwdCheckReq.CaptchaId) {
	case global.CAPTCHA_STATUS_ABSENT: // 如果未进行人机验证
		captchaId := cache.CreateCaptchaStatus()
		resp.Result(ctx, -1, gin.H{"captchaId": captchaId}, "需要人机验证")
		return
	case global.CAPTCHA_STATUS_PASS:
		cache.DelCaptchaStatus(modifyPwdCheckReq.CaptchaId) // 删除人机验证状态
	}

	if err := service.ResetPwdCheck(ctx, modifyPwdCheckReq.Email); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

// 修改密码
func ModifyPwd(ctx *gin.Context) {
	var modifyPwdReq dto.ModifyPwdReq
	if err := ctx.Bind(&modifyPwdReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if !utils.VerifyEmail(modifyPwdReq.Email) {
		resp.FailWithMessage(ctx, "邮箱格式错误")
		return
	}

	if utils.VerifyStringLength(modifyPwdReq.Password, "<", 6) {
		resp.FailWithMessage(ctx, "密码长度不能小于6位")
		return
	}

	if !utils.VerifyStringLength(modifyPwdReq.Code, "=", 6) {
		resp.FailWithMessage(ctx, "验证码长度为6位")
		return
	}

	if err := service.ModifyPwd(ctx, modifyPwdReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}
