package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectPGCRoutes(r *gin.RouterGroup) {
	pgcGroup := r.Group("pgc")
	{
		pgcGroup.GET("list", api.GetPGCList)
		pgcGroup.GET("detail", api.GetPGCDetail)
		pgcGroup.GET("/:pgc_id/episodes", api.GetPGCEpisodes)
		pgcGroup.GET("search", api.SearchPGC)
		pgcGroup.GET("type/:type", api.GetPGCByType)
		pgcGroup.GET("ongoing", api.GetOngoingPGC)
		pgcGroup.GET("recommended", api.GetRecommendedPGC)
		// 参考 B 站：PGC 推荐（按 seed/type 召回 + 可播过滤）
		pgcGroup.GET("recommend", api.RecommendPGC)
		// 播放页：按当前 vid 推荐同类 PGC
		pgcGroup.GET("recommend-by-video", api.RecommendPGCByVideo)
		pgcGroup.GET("detail-with-episodes", api.GetPGCDetailWithEpisodes)
		pgcGroup.GET("play-panel-by-video", api.GetPGCPlayPanelByVideo)
		pgcGroup.GET("episode-detail", api.GetPGCEpisodeDetail)

		// 后台审核（仅需登录 + Casbin，勿放在 :pgc_id 动态路由之后）
		pgcAdmin := pgcGroup.Group("")
		pgcAdmin.Use(middleware.Auth())
		{
			pgcAdmin.POST("getManageList", api.GetPGCManageList)
			pgcAdmin.POST("getReviewList", api.GetPGCReviewList)
			pgcAdmin.POST("reviewApproved", api.ReviewPGCApproved)
			pgcAdmin.POST("reviewFailed", api.ReviewPGCFailed)
			pgcAdmin.POST("adminUpdateStatus", api.AdminUpdatePGCStatus)
			pgcAdmin.DELETE("adminDelete/:pgc_id", api.AdminDeletePGC)
		}

		pgcAuth := pgcGroup.Group("")
		pgcAuth.Use(middleware.Auth(), middleware.RequirePGC())
		{
			pgcAuth.POST("create", api.CreatePGC)
			pgcAuth.PUT("update", api.UpdatePGC)
			pgcAuth.PUT("/:pgc_id/status", api.UpdatePGCStatus)
			pgcAuth.DELETE("/:pgc_id", api.DeletePGC)
			pgcAuth.POST("/:pgc_id/episodes/add", api.AddPGCEpisode)
			pgcAuth.DELETE("/:pgc_id/episodes/:id", api.DeletePGCEpisode)
			pgcAuth.PUT("/:pgc_id/episodes/:id/bind", api.BindPGCEpisodeVideo)
			pgcAuth.PUT("/:pgc_id/episodes/:id", api.UpdatePGCEpisode)
			pgcAuth.PUT("/:pgc_id/episodes/:id/status", api.UpdatePGCEpisodeStatus)
		}
	}
}
