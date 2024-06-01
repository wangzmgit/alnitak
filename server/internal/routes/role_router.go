package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectRoleRoutes(r *gin.RouterGroup) {

	roleGroup := r.Group("role")
	roleGroup.Use(middleware.Auth())
	{
		// 获取个人角色信息
		roleGroup.GET("getRoleInfo", api.GetRoleInfo)
		// 新增角色
		roleGroup.POST("addRole", api.AddRole)
		// 获取角色列表
		roleGroup.POST("getRoleList", api.GetRoleList)
		// 获取全部角色
		roleGroup.GET("getAllRoleList", api.GetAllRoleList)
		// 编辑角色
		roleGroup.PUT("editRole", api.EditRole)
		// 删除角色
		roleGroup.DELETE("deleteRole/:id", api.DeleteRole)
		// 编辑角色首页
		roleGroup.PUT("editRoleHome", api.EditRoleHome)
	}
}
