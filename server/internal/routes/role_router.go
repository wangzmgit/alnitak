package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectRoleRoutes(r *gin.RouterGroup) {
	roleApi := api.RoleApi{}

	role := r.Group("role")
	role.Use(middleware.Auth())
	{
		// 获取个人角色信息
		role.GET("getRoleInfo", roleApi.GetRoleInfo)
		// 新增角色
		role.POST("addRole", roleApi.AddRole)
		// 获取角色列表
		role.POST("getRoleList", roleApi.GetRoleList)
		// 编辑角色
		role.PUT("editRole", roleApi.EditRole)
		// 删除角色
		role.DELETE("deleteRole/:id", roleApi.DeleteRole)
		// 编辑角色首页
		role.PUT("editRoleHome", roleApi.EditRoleHome)
	}
}
