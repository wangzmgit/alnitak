package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// 获取视频文件
func GetVideoFile(ctx *gin.Context) {
	quality := ctx.Query("quality")
	resourceId := utils.StringToUint(ctx.Query("resource_id"))

	file, err := service.GetVideoFile(ctx, resourceId, quality)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	ctx.Writer.Header().Set("Content-type", "text/plain; charset=utf-8")
	resp.OkWithString(ctx, file)
}

// 获取视频切片
func GetVideoSlice(ctx *gin.Context) {
	key := ctx.Query("key")
	dir := ctx.Query("dir")
	file := ctx.Query("file")

	if !service.GetVideoSliceAvailableStatus(key) {
		resp.Result(ctx, http.StatusForbidden, map[string]interface{}{}, "ok")
		return
	}

	// 使用本地存储
	if viper.GetString("storage.oss_type") == "local" {
		ctx.File("./upload/video/" + dir + "/" + file)
		return
	}

	// 使用oss且但不使用cdn
	redirect := global.Storage.GetObjectUrl("video/" + dir + "/" + file)
	ctx.Redirect(http.StatusMovedPermanently, redirect)
}

// cdn远程鉴权
func VideoRemoteAuth(ctx *gin.Context) {
	key := ctx.Query("key")

	if !service.GetVideoSliceAvailableStatus(key) {
		resp.Result(ctx, http.StatusForbidden, map[string]interface{}{}, "ok")
		return
	} else {
		resp.Ok(ctx)
	}
}

// 获取图片文件
func GetImgFile(ctx *gin.Context) {
	file := ctx.Param("file")

	// 使用本地存储
	if viper.GetString("storage.oss_type") == "local" {
		ctx.File("./upload/image/" + file)
		return
	}

	// 不使用oss
	redirect := global.Storage.GetObjectUrl("image/" + file)
	fmt.Println("redirect", redirect, "image/"+file)
	ctx.Redirect(http.StatusMovedPermanently, redirect)
}
