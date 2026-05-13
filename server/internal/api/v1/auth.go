package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	jwt_parse "interastral-peace.com/alnitak/pkg/jwt"
	"interastral-peace.com/alnitak/utils"
)

const refreshCookieName = "refresh_token"

// 设置 refresh_token 的 HttpOnly Cookie
func setRefreshCookie(ctx *gin.Context, refreshToken string) {
	secure := ctx.Request.TLS != nil
	ctx.SetSameSite(http.SameSiteLaxMode)
	maxAge := 7 * 24 * 60 * 60
	ctx.SetCookie(refreshCookieName, refreshToken, maxAge, "/", "", secure, true)
}

func clearRefreshCookie(ctx *gin.Context) {
	secure := ctx.Request.TLS != nil
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(refreshCookieName, "", -1, "/", "", secure, true)
}

// 注册
func Register(ctx *gin.Context) {
	var registerReq dto.RegisterReq
	if err := ctx.Bind(&registerReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

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

	resp.Ok(ctx)
}

// 登录
func Login(ctx *gin.Context) {
	var loginReq dto.LoginReq
	if err := ctx.Bind(&loginReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if !utils.VerifyEmail(loginReq.Email) {
		resp.FailWithMessage(ctx, "邮箱格式错误")
		return
	}

	if utils.VerifyStringLength(loginReq.Password, "<", 6) {
		resp.FailWithMessage(ctx, "密码长度不能小于6位")
		return
	}

	if utils.VerifyStringLength(loginReq.CaptchaId, ">", 0) {
		if cache.GetCaptchaStatus(loginReq.CaptchaId) == global.CAPTCHA_STATUS_PASS {
			cache.DelLoginTryCount(loginReq.Email)
			cache.DelCaptchaStatus(loginReq.CaptchaId)
		}
	}

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

	setRefreshCookie(ctx, refreshToken)
	resp.OkWithData(ctx, gin.H{"token": accessToken, "refreshToken": refreshToken, "userId": userId})
}

// 邮箱登录
func EmailLogin(ctx *gin.Context) {
	var loginReq dto.EmailLoginReq
	if err := ctx.Bind(&loginReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if !utils.VerifyEmail(loginReq.Email) {
		resp.FailWithMessage(ctx, "邮箱格式错误")
		return
	}

	if utils.VerifyStringLength(loginReq.CaptchaId, ">", 0) {
		if cache.GetCaptchaStatus(loginReq.CaptchaId) == global.CAPTCHA_STATUS_PASS {
			cache.DelLoginTryCount(loginReq.Email)
			cache.DelCaptchaStatus(loginReq.CaptchaId)
		}
	}

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

	setRefreshCookie(ctx, refreshToken)
	resp.OkWithData(ctx, gin.H{"token": accessToken, "refreshToken": refreshToken, "userId": userId})
}

// 刷新token
func UpdateToken(ctx *gin.Context) {
	var tokenReq dto.TokenReq
	_ = ctx.ShouldBindJSON(&tokenReq)
	if tokenReq.RefreshToken == "" {
		if v, err := ctx.Cookie(refreshCookieName); err == nil {
			tokenReq.RefreshToken = v
		}
	}
	if tokenReq.RefreshToken == "" {
		resp.Result(ctx, 2000, nil, "需要登录")
		return
	}

	accessToken, refreshToken, userId, err := service.UpdateToken(ctx, tokenReq)
	if err != nil {
		resp.Result(ctx, 2000, nil, err.Error())
		return
	}

	if refreshToken != "" {
		setRefreshCookie(ctx, refreshToken)
	}

	resp.OkWithData(ctx, gin.H{"token": accessToken, "refreshToken": refreshToken, "userId": userId})
}

// 当前会话用户信息
func Me(ctx *gin.Context) {
	if tokenString := ctx.GetHeader("Authorization"); tokenString != "" {
		token, claims, err := jwt_parse.ParseToken(tokenString)
		if err == nil && token != nil && token.Valid && claims.TokenType == 0 {
			user := service.GetUserInfo(claims.UserId)
			resp.OkWithData(ctx, gin.H{"userInfo": user})
			return
		}
	}

	rt, err := ctx.Cookie(refreshCookieName)
	if err != nil || rt == "" {
		resp.Result(ctx, 2000, nil, "需要登录")
		return
	}

	accessToken, refreshToken, userId, err := service.UpdateToken(ctx, dto.TokenReq{RefreshToken: rt})
	if err != nil {
		clearRefreshCookie(ctx)
		resp.Result(ctx, 2000, nil, "需要登录")
		return
	}
	if refreshToken != "" {
		setRefreshCookie(ctx, refreshToken)
	}

	user := service.GetUserInfo(userId)
	resp.OkWithData(ctx, gin.H{"userInfo": user, "token": accessToken, "refreshToken": refreshToken, "userId": userId})
}

// 退出登录
func Logout(ctx *gin.Context) {
	var tokenReq dto.TokenReq
	_ = ctx.ShouldBindJSON(&tokenReq)
	if tokenReq.RefreshToken == "" {
		if v, err := ctx.Cookie(refreshCookieName); err == nil {
			tokenReq.RefreshToken = v
		}
	}

	service.Logout(ctx, tokenReq)
	clearRefreshCookie(ctx)
	resp.Ok(ctx)
}

// 重置密码检查
func ResetPwdCheck(ctx *gin.Context) {
	var modifyPwdCheckReq dto.ModifyPwdCheckReq
	if err := ctx.Bind(&modifyPwdCheckReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

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
	case global.CAPTCHA_STATUS_ABSENT:
		captchaId := cache.CreateCaptchaStatus()
		resp.Result(ctx, -1, gin.H{"captchaId": captchaId}, "需要人机验证")
		return
	case global.CAPTCHA_STATUS_PASS:
		cache.DelCaptchaStatus(modifyPwdCheckReq.CaptchaId)
	}

	if err := service.ResetPwdCheck(ctx, modifyPwdCheckReq.Email); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// 修改密码
func ModifyPwd(ctx *gin.Context) {
	var modifyPwdReq dto.ModifyPwdReq
	if err := ctx.Bind(&modifyPwdReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

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

	resp.Ok(ctx)
}

// ========== 认证类型 API ==========

// AddAuthType 添加认证类型
func AddAuthType(c *gin.Context) {
	var req dto.AddAuthTypeReq
	if err := c.ShouldBind(&req); err != nil {
		resp.FailWithMessage(c, "参数有误")
		return
	}

	if err := service.AddAuthType(c, req); err != nil {
		resp.FailWithMessage(c, err.Error())
		return
	}

	resp.Ok(c)
}

// EditAuthType 编辑认证类型
func EditAuthType(c *gin.Context) {
	var req dto.EditAuthTypeReq
	if err := c.ShouldBind(&req); err != nil {
		resp.FailWithMessage(c, "参数有误")
		return
	}

	if err := service.EditAuthType(c, req); err != nil {
		resp.FailWithMessage(c, err.Error())
		return
	}

	resp.Ok(c)
}

// DeleteAuthType 删除认证类型
func DeleteAuthType(c *gin.Context) {
	id := utils.StringToUint(c.Param("id"))
	if id == 0 {
		resp.FailWithMessage(c, "参数有误")
		return
	}

	if err := service.DeleteAuthType(c, id); err != nil {
		resp.FailWithMessage(c, err.Error())
		return
	}

	resp.Ok(c)
}

// GetAuthTypeList 获取认证类型列表（公开）
func GetAuthTypeList(c *gin.Context) {
	category := c.Query("category")
	list := service.GetAuthTypeList(category)
	resp.OkWithData(c, gin.H{"list": list})
}

// GetAllAuthTypeList 获取所有认证类型列表（管理用）
func GetAllAuthTypeList(c *gin.Context) {
	page := utils.StringToInt(c.Query("page"))
	pageSize := utils.StringToInt(c.Query("pageSize"))
	if pageSize > 30 {
		pageSize = 30
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if page <= 0 {
		page = 1
	}

	total, list := service.GetAllAuthTypeList(page, pageSize)
	resp.OkWithData(c, gin.H{
		"total": total,
		"list":  list,
	})
}

// GetAuthTypeByID 获取认证类型详情
func GetAuthTypeByID(c *gin.Context) {
	id := utils.StringToUint(c.Param("id"))
	if id == 0 {
		resp.FailWithMessage(c, "参数有误")
		return
	}

	authType, err := service.GetAuthTypeByID(id)
	if err != nil {
		resp.FailWithMessage(c, err.Error())
		return
	}

	resp.OkWithData(c, authType)
}

// ========== 用户认证 API ==========

// AddUserAuth 添加用户认证
func AddUserAuth(c *gin.Context) {
	var req dto.AddUserAuthReq
	if err := c.ShouldBind(&req); err != nil {
		resp.FailWithMessage(c, "参数有误")
		return
	}

	if err := service.AddUserAuth(c, req); err != nil {
		resp.FailWithMessage(c, err.Error())
		return
	}

	resp.Ok(c)
}

// EditUserAuth 编辑用户认证
func EditUserAuth(c *gin.Context) {
	var req dto.EditUserAuthReq
	if err := c.ShouldBind(&req); err != nil {
		resp.FailWithMessage(c, "参数有误")
		return
	}

	if err := service.EditUserAuth(c, req); err != nil {
		resp.FailWithMessage(c, err.Error())
		return
	}

	resp.Ok(c)
}

// DeleteUserAuth 删除用户认证
func DeleteUserAuth(c *gin.Context) {
	var req dto.DeleteUserAuthReq
	if err := c.ShouldBind(&req); err != nil {
		resp.FailWithMessage(c, "参数有误")
		return
	}

	if err := service.DeleteUserAuth(c, req); err != nil {
		resp.FailWithMessage(c, err.Error())
		return
	}

	resp.Ok(c)
}

// GetUserAuthList 获取用户认证列表（公开，用户自己查看）
func GetUserAuthList(c *gin.Context) {
	uid := utils.StringToUint(c.Query("uid"))
	if uid == 0 {
		resp.FailWithMessage(c, "参数有误")
		return
	}

	list := service.GetUserAuthList(uid)
	resp.OkWithData(c, gin.H{"list": list})
}

// GetUserAuthListWithUser 获取用户认证列表（管理用，带用户信息）
func GetUserAuthListWithUser(c *gin.Context) {
	page := utils.StringToInt(c.Query("page"))
	pageSize := utils.StringToInt(c.Query("pageSize"))
	authTypeID := utils.StringToUint(c.Query("authTypeId"))

	if pageSize > 30 {
		pageSize = 30
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if page <= 0 {
		page = 1
	}

	total, list := service.GetUserAuthListWithUser(page, pageSize, authTypeID)
	resp.OkWithData(c, gin.H{
		"total": total,
		"list":  list,
	})
}

// GetUserAuthByID 获取用户认证详情
func GetUserAuthByID(c *gin.Context) {
	id := utils.StringToUint(c.Param("id"))
	if id == 0 {
		resp.FailWithMessage(c, "参数有误")
		return
	}

	userAuth, err := service.GetUserAuthByID(id)
	if err != nil {
		resp.FailWithMessage(c, err.Error())
		return
	}

	resp.OkWithData(c, userAuth)
}

// GetUserAuthByUid 获取指定用户的认证信息
func GetUserAuthByUid(c *gin.Context) {
	uid := utils.StringToUint(c.Param("uid"))
	if uid == 0 {
		resp.FailWithMessage(c, "参数有误")
		return
	}

	list, err := service.GetUserAuthByUid(uid)
	if err != nil {
		resp.FailWithMessage(c, "获取失败")
		return
	}

	resp.OkWithData(c, gin.H{"list": list})
}

// GetUserPrimaryAuth 获取用户主要认证（用于前端展示）
func GetUserPrimaryAuth(c *gin.Context) {
	uid := utils.StringToUint(c.Query("uid"))
	if uid == 0 {
		resp.FailWithMessage(c, "参数有误")
		return
	}

	auth := service.GetUserPrimaryAuth(uid)
	if auth == nil {
		resp.OkWithData(c, gin.H{"auth": nil})
		return
	}
	resp.OkWithData(c, gin.H{"auth": auth})
}