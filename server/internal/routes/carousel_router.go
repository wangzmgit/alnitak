package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	"interastral-peace.com/alnitak/internal/middleware"
)

func CollectCarouselRoutes(r *gin.RouterGroup) {
	carouselGroup := r.Group("carousel")

	carouselAuth := carouselGroup.Group("")
	carouselAuth.Use(middleware.Auth())
	{
		// 新增轮播图
		carouselAuth.POST("addCarousel", api.AddCarousel)
		// 获取轮播图列表
		carouselAuth.POST("getCarouselList", api.GetCarouselList)
		// 编辑轮播图信息
		carouselAuth.PUT("editCarousel", api.EditCarousel)
		// 删除轮播图
		carouselAuth.DELETE("deleteCarousel/:id", api.DeleteCarousel)
	}

	carouselGroup.GET("getCarousel", api.GetCarousel)
}
