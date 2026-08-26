package api

import (
	"net/http"
	"strings"

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
func setRefreshCookie(ctx *gin.Context, refreshToken string, maxAge int) {
	secure := ctx.Request.TLS != nil
	ctx.SetSameSite(http.SameSiteLaxMode)
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

	// 人机验证
	if utils.VerifyStringLength(registerReq.CaptchaId, "=", 0) {
		captchaId := cache.CreateCaptchaStatus()
		resp.Result(ctx, -1, gin.H{"captchaId": captchaId}, "需要人机验证")
		return
	}

	switch cache.GetCaptchaStatus(registerReq.CaptchaId) {
	case global.CAPTCHA_STATUS_ABSENT:
		captchaId := cache.CreateCaptchaStatus()
		resp.Result(ctx, -1, gin.H{"captchaId": captchaId}, "需要人机验证")
		return
	case global.CAPTCHA_STATUS_PASS:
		cache.DelCaptchaStatus(registerReq.CaptchaId)
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

	maxAge := 7 * 24 * 60 * 60 // 默认 7 天
	if loginReq.RememberMe {
		maxAge = 30 * 24 * 60 * 60 // 勾选"记住我" 30 天
	}
	setRefreshCookie(ctx, refreshToken, maxAge)
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

	if !utils.VerifyStringLength(loginReq.Code, "=", 6) {
		resp.FailWithMessage(ctx, "验证码长度为6位")
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

	maxAge := 7 * 24 * 60 * 60
	if loginReq.RememberMe {
		maxAge = 30 * 24 * 60 * 60
	}
	setRefreshCookie(ctx, refreshToken, maxAge)
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
		// 活跃用户刷新 token 用长过期时间
		setRefreshCookie(ctx, refreshToken, 30*24*60*60)
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
		setRefreshCookie(ctx, refreshToken, 30*24*60*60)
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

	// 将当前 accessToken 加入黑名单，使其立即失效
	if accessToken := ctx.GetHeader("Authorization"); accessToken != "" {
		if _, claims, err := jwt_parse.ParseToken(accessToken); err == nil && claims != nil {
			if claims.TokenType == 0 && claims.ExpiresAt != nil {
				// ParseToken 内部使用了 trimBearer，但原始 accessToken 可能含 "Bearer " 前缀，
				// 而 Auth/OptionalAuth 中间件传入 IsAccessTokenBlacklisted 的是 trimBearer 之后的值。
				// 为保证 hash key 一致，此处统一去除前缀（大小写不敏感）。
				cleanToken := accessToken
				if len(cleanToken) > 7 && strings.EqualFold(cleanToken[:7], "Bearer ") {
					cleanToken = cleanToken[7:]
				}
				cache.SetBlacklistedAccessToken(cleanToken, claims.ExpiresAt.Time)
			}
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

// 修改密码（已登录用户，需校验旧密码）
func ChangePassword(ctx *gin.Context) {
	var req dto.ChangePasswordReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if req.NewPassword == "" || req.OldPassword == "" {
		resp.FailWithMessage(ctx, "参数不能为空")
		return
	}

	userId, _ := ctx.Get("userId")
	if err := service.ChangePassword(ctx, userId.(uint), req); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 改密后清除本地凭证，弹回登录页（前端拦截器处理）
	resp.OkWithMessage(ctx, "密码修改成功，请重新登录")
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