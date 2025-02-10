package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 获取视频文件
func GetVideoFile(ctx *gin.Context) {
	quality := ctx.Query("quality")
	resourceId := utils.StringToUint(ctx.Query("resourceId"))

	file, err := service.GetVideoFile(ctx, resourceId, quality)
	if err != nil {
		resp.ForbiddenWithMessage(ctx, err.Error())
		return
	}

	ctx.Writer.Header().Set("Content-type", "text/plain; charset=utf-8")
	resp.OkWithString(ctx, file)
}

// 获取视频文件(后台管理)
func GetVideoFileManage(ctx *gin.Context) {
	quality := ctx.Query("quality")
	resourceId := utils.StringToUint(ctx.Query("resourceId"))

	file, err := service.GetVideoFileManage(ctx, resourceId, quality)
	if err != nil {
		resp.ForbiddenWithMessage(ctx, err.Error())
		return
	}

	ctx.Writer.Header().Set("Content-type", "text/plain; charset=utf-8")
	resp.OkWithString(ctx, file)
}

// 获取视频切片
func GetVideoSlice(ctx *gin.Context) {
	key := ctx.Query("key")
	file := ctx.Param("file")

	dir := service.GetVideoSliceDir(key)
	if dir == "" {
		resp.Forbidden(ctx)
		return
	}

	// 使用本地存储
	if global.Config.Storage.OssType == "local" {
		ctx.File("./upload/video/" + dir + "/" + file)
		return
	}

	// 使用oss
	redirect := global.Storage.GetObjectUrl("video/" + dir + "/" + file)
	ctx.Redirect(http.StatusMovedPermanently, redirect)
}

// 获取图片文件
func GetImgFile(ctx *gin.Context) {
	file := ctx.Param("file")

	// 使用本地存储
	if global.Config.Storage.OssType == "local" {
		ctx.File("./upload/image/" + file)
		return
	}

	// 不使用oss
	redirect := global.Storage.GetObjectUrl("image/" + file)

	ctx.Redirect(http.StatusMovedPermanently, redirect)
}
