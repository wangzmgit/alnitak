package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectAuthRoutes(r *gin.RouterGroup) {

	authGroup := r.Group("auth")

	authAuth := authGroup.Group("")
	authAuth.Use(middleware.Auth())
	{
		authAuth.POST("logout", api.Logout)
	}

	// 用户注册
	authGroup.POST("register", api.Register)
	// 用户登录(密码)
	authGroup.POST("login", api.Login)
	// 用户登录(邮箱)
	authGroup.POST("login/email", api.EmailLogin)
	// 更新token
	authGroup.POST("updateToken", api.UpdateToken)
	// 清除Cookie
	authGroup.POST("clearCookie", api.ClearCookie)
	// 修改密码检查
	authGroup.POST("resetpwdCheck", api.ResetPwdCheck)
	// 修改密码
	authGroup.POST("modifyPwd", api.ModifyPwd)
}
