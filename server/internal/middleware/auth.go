package middleware

import (
	"errors"

	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	jwt_parse "interastral-peace.com/alnitak/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 读取验证token
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			resp.Result(ctx, 3000, nil, "TOKEN无效")
			ctx.Abort()
			return
		}
		// 验证并解析token
		_, claims, err := jwt_parse.ParseToken(tokenString)
		if err != nil {
			// accessToken 过期
			if errors.Is(err, jwt.ErrTokenExpired) && claims.TokenType == 0 {
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
		if claims.TokenType == 0 { // accessToken
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
			ctx.Set("roleCode", user.Role)
			ctx.Next()
		} else {
			// resp.FailWithMessage(ctx, "token验证失败")
			resp.Result(ctx, 2000, nil, "token验证失败")
			ctx.Abort()
		}
	}
}

func WsAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 读取验证token
		tokenString := ctx.Query("token")
		// 验证并解析token
		_, claims, err := jwt_parse.ParseToken(tokenString)
		if err != nil {
			resp.Result(ctx, 2000, nil, "token验证失败")
			ctx.Abort()
			return
		}

		ctx.Set("userId", claims.UserId)
		ctx.Next()
	}
}
