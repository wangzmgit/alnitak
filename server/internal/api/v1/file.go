package api

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/pkg/playtoken"
	"interastral-peace.com/alnitak/utils"
)

var sliceMediaFileRegexp = regexp.MustCompile(`^[A-Za-z0-9_.-]+\.(m4s|ts)$`)

func sanitizeSliceMediaFileName(file string) (string, error) {
	file = strings.TrimSpace(file)
	if file == "" {
		return "", errors.New("empty file")
	}

	// 必须是纯文件名，不允许包含目录分隔符（防路径穿越）
	if filepath.Base(file) != file {
		return "", errors.New("invalid file name")
	}

	// 限定扩展名与字符集合，避免构造任意路径/协议
	if !sliceMediaFileRegexp.MatchString(file) {
		return "", errors.New("invalid file name format")
	}

	return file, nil
}

// resolveStreamDir 支持 legacy key→Redis 映射，或 st=JWT（dir+file 绑定防篡改）。
func resolveStreamDir(ctx *gin.Context, key, file string) (dir string, ok bool) {
	if st := strings.TrimSpace(ctx.Query("st")); st != "" {
		claims, err := playtoken.ParseStreamToken(st)
		if err != nil || claims.File != file || claims.Dir == "" {
			return "", false
		}
		return claims.Dir, true
	}
	if key == "" {
		return "", false
	}
	d := service.GetVideoSliceDir(key)
	if d == "" {
		return "", false
	}
	return d, true
}

// 获取视频文件
func GetVideoFile(ctx *gin.Context) {
	quality := ctx.Query("quality")
	resourceId, parseErr := service.ParseResourceID(ctx.Query("resourceId"))
	if parseErr != nil {
		resp.FailWithMessage(ctx, parseErr.Error())
		return
	}
	format := ctx.DefaultQuery("format", "m3u8") // m3u8 / mpd / dash / m3u8video / m3u8audio

	file, err := service.GetVideoFile(ctx, resourceId, quality, format)
	if err != nil {
		resp.ForbiddenWithMessage(ctx, err.Error())
		return
	}

	// 关键：设置 Cache-Control 为 5 小时（与 OSS 签名过期时间一致）
	// 防止浏览器缓存 m3u8/mpd 播放列表，导致使用旧的签名 URL
	ctx.Header("Cache-Control", "public, max-age=18000, must-revalidate")
	ctx.Header("Pragma", "no-cache")

	switch format {
	case "mpd", "dash", "dash-unified":
		ctx.Data(http.StatusOK, "application/dash+xml; charset=utf-8", []byte(file))
	case "m3u8", "m3u8video", "m3u8audio":
		ctx.Data(http.StatusOK, "application/vnd.apple.mpegurl; charset=utf-8", []byte(file))
	default:
		ctx.Data(http.StatusOK, "application/json; charset=utf-8", []byte(file))
	}
}

// 获取视频文件(后台管理)
func GetVideoFileManage(ctx *gin.Context) {
	// 禁用缓存，确保审核时看到最新视频文件
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	quality := ctx.Query("quality")
	resourceId := utils.StringToUint(ctx.Query("resourceId"))
	format := ctx.DefaultQuery("format", "m3u8")

	file, err := service.GetVideoFileManage(ctx, resourceId, quality, format)
	if err != nil {
		resp.ForbiddenWithMessage(ctx, err.Error())
		return
	}

	switch format {
	case "mpd", "dash", "dash-unified":
		ctx.Data(http.StatusOK, "application/dash+xml; charset=utf-8", []byte(file))
	case "m3u8", "m3u8video", "m3u8audio":
		ctx.Data(http.StatusOK, "application/vnd.apple.mpegurl; charset=utf-8", []byte(file))
	default:
		ctx.Data(http.StatusOK, "application/json; charset=utf-8", []byte(file))
	}
}

// 获取视频切片（兼容模式：SegmentList）
func GetVideoSlice(ctx *gin.Context) {
	key := ctx.Query("key")
	file, err := sanitizeSliceMediaFileName(ctx.Param("file"))
	if err != nil {
		resp.Forbidden(ctx)
		return
	}

	dir, ok := resolveStreamDir(ctx, key, file)
	if !ok {
		resp.Forbidden(ctx)
		return
	}

	// 使用本地存储
	if global.Config.Storage.OssType == "local" {
		ctx.File("./upload/video/" + dir + "/" + file)
		return
	}

	// OSS 存储：重定向到公开URL
	redirect := global.GetOssUrl("video/" + dir + "/" + file)
	ctx.Redirect(http.StatusFound, redirect)
}

// GetVideoStream 获取视频流（B站风格：支持字节范围请求）
func GetVideoStream(ctx *gin.Context) {
	key := ctx.Query("key")
	file, err := sanitizeSliceMediaFileName(ctx.Param("file"))
	if err != nil {
		resp.Forbidden(ctx)
		return
	}

	dir, ok := resolveStreamDir(ctx, key, file)
	if !ok {
		resp.Forbidden(ctx)
		return
	}

	filePath := "./upload/video/" + dir + "/" + file

	// 使用本地存储时，支持 Range 请求
	if global.Config.Storage.OssType == "local" {
		// 设置 CORS 头
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Range")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Content-Range, Accept-Ranges")

		// 处理 OPTIONS 预检请求
		if ctx.Request.Method == "OPTIONS" {
			ctx.Status(http.StatusOK)
			return
		}

		// .m4s 文件 Go 标准库不识别 MIME 类型，需要手动设置
		// 否则 iOS 原生 HLS 播放器可能无法识别
		if strings.HasSuffix(file, ".m4s") {
			ctx.Header("Content-Type", "video/iso.segment")
		}

		// 使用 http.ServeFile，它会自动处理：
		// - Range 请求（返回 206）
		// - Content-Type 检测（对于已知类型）
		// - 缓存头（Last-Modified, ETag）
		// - 高效的文件传输
		http.ServeFile(ctx.Writer, ctx.Request, filePath)
		return
	}

	// OSS 存储：重定向到公开URL
	redirect := global.GetOssUrl("video/" + dir + "/" + file)
	ctx.Redirect(http.StatusFound, redirect)
}

// 获取图片文件
func GetImgFile(ctx *gin.Context) {
	file := ctx.Param("file")

	// 使用本地存储
	if viper.GetString("storage.oss_type") == "local" {
		// 设置缓存头，告知浏览器缓存一天
		ctx.Header("Cache-Control", "public, max-age=86400, must-revalidate")
		ctx.File("./upload/image/" + file)
		return
	}

	// OSS 存储：重定向到公开URL
	redirect := global.GetOssUrl("image/" + file)

	// 开发模式下打印重定向信息
	if global.Config.Log.Mode == "dev" {
		fmt.Println("redirect", redirect, "image/"+file)
	}

	// 设置缓存头，告知浏览器缓存5小时（与OSS签名过期时间一致）
	ctx.Header("Cache-Control", "public, max-age=18000, must-revalidate")
	ctx.Redirect(http.StatusFound, redirect)
}

// GetSubtitleFile GET /api/subtitle/:file（file 为 snowflake.vtt，须已在 subtitle_track 中登记；策略与 GetImgFile 一致）
func GetSubtitleFile(ctx *gin.Context) {
	file := ctx.Param("file")
	localPath, objectKey, ok := service.GetSubtitleTrackForFileServe(ctx, file)
	if !ok {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
	if ctx.Request.Method == http.MethodOptions {
		ctx.Status(http.StatusOK)
		return
	}

	if global.Config.Storage.OssType == "local" {
		ctx.Header("Cache-Control", "public, max-age=1800, must-revalidate")
		ctx.Header("Content-Type", "text/vtt; charset=utf-8")
		ctx.File(localPath)
		return
	}

	// OSS：与图片路由相同，302 到公开 URL / 预签名 URL（播放端须在 <video crossorigin> + 存储 CORS）
	redirect := global.GetOssUrl(objectKey)
	if global.Config.Log.Mode == "dev" {
		fmt.Println("redirect subtitle", redirect, objectKey)
	}
	ctx.Header("Cache-Control", "public, max-age=18000, must-revalidate")
	ctx.Redirect(http.StatusFound, redirect)
}
