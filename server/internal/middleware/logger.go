package middleware

import (
	"net"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger 用于替换gin框架的Logger中间件，不传参数，直接这样写
func GinLogger(c *gin.Context) {
	// 已记录操作日志，这里不再记录
	c.Next() // 执行视图函数
}

// GinRecovery 用于替换gin框架的Recovery中间件，因为传入参数，再包一层
func GinRecovery(stack bool) gin.HandlerFunc {
	logger := zap.L()
	return func(c *gin.Context) {
		defer func() {
			// defer 延迟调用，出了异常，处理并恢复异常，记录日志
			if err := recover(); err != nil {
				//  这个不必须，检查是否存在断开的连接(broken pipe或者connection reset by peer)---------开始--------
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				//httputil包预先准备好的DumpRequest方法
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// 如果连接已断开，我们无法向其写入状态
					c.Error(err.(error))
					c.Abort()
					return
				}
				//  这个不必须，检查是否存在断开的连接(broken pipe或者connection reset by peer)---------结束--------

				// 是否打印堆栈信息，使用的是debug.Stack()，传入false，在日志中就没有堆栈信息
				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				// 有错误，直接返回给前端错误，前端直接报错
				//c.AbortWithStatus(http.StatusInternalServerError)
				// 该方式前端不报错
				c.JSON(200, gin.H{"code": 500, "msg": "服务器异常"})
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}
