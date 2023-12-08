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
				service.AddFailOperation(ctx, "权限", "用户Token过期", gin.H{"用户ID": claims.UserId}, nil)
				resp.Result(ctx, 3000, nil, "TOKEN过期")
				ctx.Abort()
				return
			}

			service.AddFailOperation(ctx, "权限", "Token验证失败", gin.H{"用户ID": claims.UserId}, err)
			resp.FailWithMessage(ctx, "验证失败")
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
				service.AddFailOperation(ctx, "权限", "权限不足", gin.H{
					"用户ID": claims.UserId,
					"请求方式": act,
					"请求路径": obj,
				}, nil)
				resp.FailWithMessage(ctx, "权限不足")
				ctx.Abort()
				return
			}

			ctx.Set("userId", claims.UserId)
			ctx.Set("roleCode", user.Role)
			ctx.Next()
		} else {
			service.AddFailOperation(ctx, "权限", "Token类型错误", gin.H{"用户ID": claims.UserId}, nil)
			resp.FailWithMessage(ctx, "token验证失败")
			ctx.Abort()
		}
	}
}
