package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
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
	r.Use(middleware.OperationRecord())
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
		// 视频相关路由
		CollectVideoRoutes(v1)
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
	}

	//获取静态文件
	r.GET("/api/image/:file", api.GetImgFile)

	return r
}
