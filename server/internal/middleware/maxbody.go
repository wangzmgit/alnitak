package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MaxBodySize 限制请求体的最大字节数，超过则返回 413 Request Entity Too Large。
// 通过在 Handler 之前包装 http.MaxBytesReader 实现。
func MaxBodySize(maxBytes int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxBytes)
		ctx.Next()
	}
}

const (
	// UploadChunkMaxBody 分片上传：每个分片最大 50MB
	UploadChunkMaxBody = 50 << 20
	// UploadImgMaxBody 图片上传：最大 10MB
	UploadImgMaxBody = 10 << 20
	// UploadJSONMaxBody 上传 JSON 接口（create/add/merge/check）：最大 1MB
	UploadJSONMaxBody = 1 << 20
)
