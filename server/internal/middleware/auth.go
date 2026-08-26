package middleware

import (
	"errors"
	"net/http"
	"strings"

	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	jwt_parse "interastral-peace.com/alnitak/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func trimBearer(token string) string {
	if len(token) > 7 && strings.EqualFold(token[:7], "Bearer ") {
		return token[7:]
	}
	return token
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 读取验证token
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			resp.Result(ctx, 3000, nil, "TOKEN无效")
			ctx.Abort()
			return
		}
		tokenString = trimBearer(tokenString)
		// 验证并解析token
		_, claims, err := jwt_parse.ParseToken(tokenString)
		if err != nil {
			// accessToken 过期（claims 非 nil 时才判断类型，防御性 nil 检查）
			if claims != nil && errors.Is(err, jwt.ErrTokenExpired) && claims.TokenType == 0 {
				if ctx.FullPath() == "/api/v1/token/update" {
					ctx.Next()
					return
				}

				// 提示需要刷新token
				resp.Result(ctx, 3000, nil, "TOKEN过期")
				ctx.Abort()
				return
			}

			// resp.FailWithMessage(ctx, "验证失败")
			resp.Result(ctx, 2000, nil, "token验证失败")
			ctx.Abort()
			return
		}

		// 验证token存在 -> 判断token类型
		// ParseToken 成功时 claims 保证非 nil
		if claims.TokenType == 0 { // accessToken
			// 检查 accessToken 是否已被吊销（登出后立即失效）
			if cache.IsAccessTokenBlacklisted(tokenString) {
				resp.Result(ctx, 3000, nil, "TOKEN已失效")
				ctx.Abort()
				return
			}

			user, _ := service.FindUserById(claims.UserId)

			// 获得用户的全部角色
			sub := user.Role
			// 获取请求方式
			act := ctx.Request.Method
			// 获得请求路径URL
			obj := ctx.FullPath()
			isPass := global.Casbin.CasbinCheck(sub, obj, act)
			if !isPass {
				resp.FailWithMessage(ctx, "权限不足")
				ctx.Abort()
				return
			}

			ctx.Set("userId", claims.UserId)
			ctx.Set("status", user.Status)
			ctx.Set("roleCode", user.Role)
			ctx.Next()
		} else {
			resp.Result(ctx, 2000, nil, "token验证失败")
			ctx.Abort()
		}
	}
}

// OptionalAuth 可选认证中间件，有token则解析设置userId，无token则跳过
func OptionalAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.Next()
			return
		}
		tokenString = trimBearer(tokenString)
		_, claims, err := jwt_parse.ParseToken(tokenString)
		if err == nil && claims != nil && claims.TokenType == 0 {
			// 跳过已吊销的 token
			if !cache.IsAccessTokenBlacklisted(tokenString) {
				ctx.Set("userId", claims.UserId)
			}
		}
		ctx.Next()
	}
}

// CsrfCheck 校验 X-Requested-With 头，防止 Cookie 类接口遭受简单 CSRF 攻击。
// 浏览器跨域请求无法手动设置 X-Requested-With，该中间件配合 SameSite=Lax 提供深度防御。
func CsrfCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetHeader("X-Requested-With") != "XMLHttpRequest" {
			resp.Result(ctx, http.StatusForbidden, nil, "非法请求")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func WsAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 优先从 Authorization header 读取（非浏览器客户端），其次从 query param 读取（浏览器 WebSocket API）
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			tokenString = ctx.Query("token")
		}
		tokenString = trimBearer(tokenString)
		if tokenString == "" {
			resp.Result(ctx, 2000, nil, "token验证失败")
			ctx.Abort()
			return
		}

		// 验证并解析token
		_, claims, err := jwt_parse.ParseToken(tokenString)
		if err != nil {
			resp.Result(ctx, 2000, nil, "token验证失败")
			ctx.Abort()
			return
		}
		if claims.TokenType != 0 {
			resp.Result(ctx, 2000, nil, "token验证失败")
			ctx.Abort()
			return
		}

		// 检查 accessToken 是否已被吊销
		if cache.IsAccessTokenBlacklisted(tokenString) {
			resp.Result(ctx, 3000, nil, "TOKEN已失效")
			ctx.Abort()
			return
		}

		ctx.Set("userId", claims.UserId)
		ctx.Next()
	}
}
