package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/middleware"
)

const defaultPort = "9000"

func InitRouter() {
	// gin 模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.GinLogger, middleware.GinRecovery(true))

	// 设置信任网络 []string nil 为不计算，避免性能消耗，上线应当设置
	_ = r.SetTrustedProxies(nil)

	//跨域中间件
	r.Use(middleware.CORS())
	//安全响应头中间件
	r.Use(middleware.SecurityHeaders())
	if !global.Config.Security.CloseRecordUserOperation {
		r.Use(middleware.OperationRecord())
	}

	// 收集添加路由
	CollectRoutes(r)

	// API / Casbin 增量同步已在 initialize.InitDefaultData() -> SyncApiData() 中执行

	// 获取HTTP端口
	httpPort := global.Config.Server.Port
	if httpPort == "" {
		httpPort = defaultPort
	}

	// 检查是否启用HTTPS
	sslConfig := global.Config.Server.Ssl
	if sslConfig.Enabled {
		httpsPort := sslConfig.Port
		if httpsPort == "" {
			httpsPort = "443"
		}

		// HTTPS 与 HTTP 需不同端口，默认 HTTP 9000、HTTPS 9001
		if httpsPort == httpPort {
			httpsPort = "9001"
		}

		zap.L().Info("启动HTTPS+HTTP双模式",
			zap.String("https_port", httpsPort),
			zap.String("http_port", httpPort),
			zap.String("cert", sslConfig.CertFile))

		// 启动 HTTPS（goroutine）
		go func() {
			zap.L().Info("HTTPS服务器已启动", zap.String("port", httpsPort))
			if err := r.RunTLS(":"+httpsPort, sslConfig.CertFile, sslConfig.KeyFile); err != nil && err != http.ErrServerClosed {
				zap.L().Error("HTTPS服务器异常退出", zap.Error(err))
			}
		}()

		// 启动 HTTP（主协程，保持 9000 便于 Admin 等 HTTP 客户端）
		zap.L().Info("HTTP服务器已启动", zap.String("port", httpPort))
		if err := r.Run(":" + httpPort); err != nil {
			zap.L().Fatal("HTTP服务器启动失败", zap.Error(err))
		}
	} else {
		// 仅 HTTP 模式
		zap.L().Info("启动HTTP服务器", zap.String("port", httpPort))
		r.Run(":" + httpPort)
	}
}

func CollectRoutes(r *gin.Engine) *gin.Engine {

	v1 := r.Group("/api/v1")
	{
		// 客户端日志上报（需登录，避免被滥用）
		clientGroup := v1.Group("client")
		clientGroup.Use(middleware.Auth())
		clientGroup.POST("log", api.ClientLog)
		// 用户认证相关路由
		CollectUserAuthRoutes(v1)
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
		// 视频相关路由
		CollectVideoRoutes(v1)
		CollectPlayRoutes(v1)
		// 资源相关路由
		CollectResourceRoutes(v1)
		// 分区相关路由
		CollectPartitionRoutes(v1)
		// 上传相关路由
		CollectUploadRoutes(v1)
		// 评论回复相关路由
		CollectCommentRoutes(v1)
		// 点赞收藏相关路由
		CollectArchiveRoutes(v1)
		// 收藏夹相关路由
		CollectCollectionRoutes(v1)
		// 用户关系相关路由
		CollectRelationRoutes(v1)
		// 弹幕相关路由
		CollectDanmakuRoutes(v1)
		// 历史记录相关路由
		CollectHistoryRoutes(v1)
		// 消息相关接口
		CollectMessageRoutes(v1)
		// 审核相关接口
		CollectReviewRoutes(v1)
		// 轮播图相关接口
		CollectCarouselRoutes(v1)
		// 文章相关接口
		CollectArticleRoutes(v1)
		// 合集相关接口
		CollectPlaylistRoutes(v1)
		// 实时在线
		CollectOnlineRoutes(v1)
		// 配置相关接口
		CollectConfigRoutes(v1)
		// 备用 OSS 重试相关接口
		CollectBackupRoutes(v1)
		// PGC相关接口
		CollectPGCRoutes(v1)
	}

	// 获取静态文件
	r.GET("/api/image/:file", api.GetImgFile)
	r.GET("/api/subtitle/:file", middleware.OptionalAuth(), api.GetSubtitleFile)
	// 后台静态页（如 PGC 管理）
	r.Static("/admin", "./static/admin")

	// 管理员：查看远程转码 Worker 状态
	r.GET("/api/v1/admin/workers", middleware.Auth(), api.ListWorkers)

	return r
}
