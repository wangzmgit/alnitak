package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectUserAuthRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("auth")

	// 用户注册
	authGroup.POST("register", middleware.RateLimiter(middleware.LoginRateLimit), api.Register)
	// 用户登录(密码)
	authGroup.POST("login", middleware.RateLimiter(middleware.LoginRateLimit), api.Login)
	// 用户登录(邮箱)
	authGroup.POST("login/email", middleware.RateLimiter(middleware.LoginRateLimit), api.EmailLogin)
	// 当前会话用户（支持 Authorization / Cookie）
	authGroup.GET("me", api.Me)
	// 退出登录：凭 body 或 HttpOnly refresh_token Cookie，无需 Authorization（行业常见做法）
	authGroup.POST("logout", middleware.CsrfCheck(), api.Logout)
	// 更新token
	authGroup.POST("updateToken", middleware.CsrfCheck(), api.UpdateToken)
	// 修改密码检查
	authGroup.POST("resetpwdCheck", middleware.RateLimiter(middleware.ModifyPwdRateLimit), api.ResetPwdCheck)
	// 修改密码
	authGroup.POST("modifyPwd", middleware.RateLimiter(middleware.ModifyPwdRateLimit), api.ModifyPwd)

	// 修改密码（已登录状态，需校验旧密码）
	authGroup.POST("changePassword", middleware.Auth(), middleware.RateLimiter(middleware.ModifyPwdRateLimit), api.ChangePassword)

	// 认证类型列表（公开）
	authGroup.GET("type/list", api.GetAuthTypeList)
	// 获取用户认证列表（公开）
	authGroup.GET("user/list", api.GetUserAuthList)
	// 获取用户主要认证（公开）
	authGroup.GET("user/primary", api.GetUserPrimaryAuth)
	// 获取指定用户的认证信息
	authGroup.GET("user/:uid/auth", api.GetUserAuthByUid)

	// 管理接口（需要登录）
	authManage := authGroup.Group("")
	authManage.Use(middleware.Auth())
	{
		// 认证类型管理
		authManage.POST("type/add", api.AddAuthType)
		authManage.PUT("type/edit", api.EditAuthType)
		authManage.DELETE("type/:id", api.DeleteAuthType)
		authManage.GET("type/all", api.GetAllAuthTypeList)
		authManage.GET("type/:id", api.GetAuthTypeByID)

		// 用户认证管理
		authManage.POST("user/add", api.AddUserAuth)
		authManage.PUT("user/edit", api.EditUserAuth)
		authManage.DELETE("user", api.DeleteUserAuth)
		authManage.GET("user/all", api.GetUserAuthListWithUser)
		authManage.GET("user/detail/:id", api.GetUserAuthByID)
	}
}
