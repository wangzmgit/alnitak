package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectMenuRoutes(r *gin.RouterGroup) {
	menuApi := api.MenuApi{}

	menu := r.Group("menu")
	menu.Use(middleware.Auth())
	{
		// 获取菜单树
		menu.GET("getMenuTree", menuApi.GetMenuTree)
		// 获取用户菜单树
		menu.GET("getUserMenu", menuApi.GetUserMenuTree)
		// 新增菜单
		menu.POST("addMenu", menuApi.AddMenu)
		// 编辑菜单
		menu.PUT("editMenu", menuApi.EditMenu)
		// 删除菜单
		menu.DELETE("deleteMenu/:id", menuApi.DeleteMenu)
		// 获取角色菜单
		menu.GET("getRoleMenu", menuApi.GetRoleMenu)
		// 编辑角色菜单
		menu.PUT("editRoleMenu", menuApi.EditRoleMenu)
	}
}
