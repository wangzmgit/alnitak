package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/middleware"
)

func InitRouter() {
	// gin 模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.GinLogger, middleware.GinRecovery(true))

	// 设置信任网络 []string nil 为不计算，避免性能消耗，上线应当设置
	_ = r.SetTrustedProxies(nil)

	//跨域中间件
	r.Use(middleware.CORS())

	// 收集添加路由
	CollectRoutes(r)

	// 运行
	r.Run(":9000")
}

func CollectRoutes(r *gin.Engine) *gin.Engine {

	v1 := r.Group("/api/v1")
	{
		// 登录注册相关路由路由
		CollectAuthRoutes(v1)
		// 用户相关路由
		CollectUserRoutes(v1)
		// 验证相关路由
		CollectVerifyRoutes(v1)
		// 角色相关路由
		CollectRoleRoutes(v1)
		// API相关路由
		CollectApiRoutes(v1)
		// 菜单相关路由
		CollectMenuRoutes(v1)
	}

	//获取静态文件
	r.StaticFS("/api/image", http.Dir("./upload/image"))
	r.StaticFS("/api/video", http.Dir("./upload/video"))
	r.StaticFS("/api/config", http.Dir("./upload/config"))

	return r
}
