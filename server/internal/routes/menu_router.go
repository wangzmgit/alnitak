package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectMenuRoutes(r *gin.RouterGroup) {

	menuGroup := r.Group("menu")
	menuGroup.Use(middleware.Auth())
	{
		// 获取菜单树
		menuGroup.GET("getMenuTree", api.GetMenuTree)
		// 获取用户菜单树
		menuGroup.GET("getUserMenu", api.GetUserMenuTree)
		// 新增菜单
		menuGroup.POST("addMenu", api.AddMenu)
		// 编辑菜单
		menuGroup.PUT("editMenu", api.EditMenu)
		// 删除菜单
		menuGroup.DELETE("deleteMenu/:id", api.DeleteMenu)
		// 获取角色菜单
		menuGroup.GET("getRoleMenu", api.GetRoleMenu)
		// 编辑角色菜单
		menuGroup.PUT("editRoleMenu", api.EditRoleMenu)
	}
}
