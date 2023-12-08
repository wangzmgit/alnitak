package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectApiRoutes(r *gin.RouterGroup) {
	apiApi := api.ApiApi{}

	api := r.Group("api")
	api.Use(middleware.Auth())
	{
		// 获取API列表
		api.POST("getApiList", apiApi.GetApiList)
		// 获取全部API列表
		api.GET("getAllApiList", apiApi.GetAllApiList)
		// 新增API
		api.POST("addApi", apiApi.AddApi)
		// 编辑API
		api.PUT("editApi", apiApi.EditApi)
		// 删除API
		api.DELETE("deleteApi/:id", apiApi.DeleteApi)
		// 获取角色API
		api.GET("getRoleApi", apiApi.GetRoleApi)
		// 编辑角色API
		api.PUT("editRoleApi", apiApi.EditRoleApi)
	}
}
