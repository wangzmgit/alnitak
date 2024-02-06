package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectApiRoutes(r *gin.RouterGroup) {

	apiGroup := r.Group("api")
	apiGroup.Use(middleware.Auth())
	{
		// 获取API列表
		apiGroup.POST("getApiList", api.GetApiList)
		// 获取全部API列表
		apiGroup.GET("getAllApiList", api.GetAllApiList)
		// 新增API
		apiGroup.POST("addApi", api.AddApi)
		// 编辑API
		apiGroup.PUT("editApi", api.EditApi)
		// 删除API
		apiGroup.DELETE("deleteApi/:id", api.DeleteApi)
		// 获取角色API
		apiGroup.GET("getRoleApi", api.GetRoleApi)
		// 编辑角色API
		apiGroup.PUT("editRoleApi", api.EditRoleApi)
	}
}
