package middleware

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/global"
)

// SecurityHeaders 设置 HTTP 安全响应头
func SecurityHeaders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 防 MIME 嗅探
		ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		// 控制 Referer 头携带范围
		ctx.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		// 限制浏览器 API 权限
		ctx.Writer.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), interest-cohort=()")
		// 禁止 iframe 嵌套
		ctx.Writer.Header().Set("X-Frame-Options", "DENY")

		// HSTS：当服务端直接终止 TLS 时设置
		if ctx.Request.TLS != nil {
			ctx.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		} else if global.Config.Server.Ssl.Enabled {
			// 未直接 TLS 但配置了 SSL（代理场景），也设置 HSTS
			ctx.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		ctx.Next()
	}
}
