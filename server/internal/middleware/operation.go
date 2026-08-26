package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/pkg/jwt"
	"interastral-peace.com/alnitak/utils"
)

var respPool sync.Pool
var bufferSize = 1024

func init() {
	respPool.New = func() interface{} {
		return make([]byte, bufferSize)
	}
}

// 跳过操作日志记录的高频接口路径
var skipOperationPaths = map[string]bool{
	"/api/v1/upload/chunkVideo": true, // 分片上传，一个视频可能有几百个分片
}

// 敏感接口路径：body 中含密码/token/验证码等，需脱敏
var sensitivePaths = map[string]bool{
	"/api/v1/auth/login":        true,
	"/api/v1/auth/login/email":  true,
	"/api/v1/auth/register":     true,
	"/api/v1/auth/modifyPwd":    true,
	"/api/v1/auth/resetpwdCheck": true,
	"/api/v1/auth/updateToken":  true,
}

func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过高频接口，避免日志表写入压力过大
		if skipOperationPaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		var body []byte
		var userId int
		if c.Request.Method != http.MethodGet {
			var err error
			body, err = io.ReadAll(c.Request.Body)
			if err != nil {
				utils.ErrorLog("读取请求错误", "middleware", err.Error())
			} else {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		} else {
			query := c.Request.URL.RawQuery
			query, _ = url.QueryUnescape(query)
			split := strings.Split(query, "&")
			m := make(map[string]string)
			for _, v := range split {
				kv := strings.Split(v, "=")
				if len(kv) == 2 {
					m[kv[0]] = kv[1]
				}
			}
			body, _ = json.Marshal(&m)
		}

		claims, _ := jwt.GetTokenClaims(trimBearer(c.GetHeader("Authorization")))
		if claims != nil && claims.UserId != 0 {
			userId = int(claims.UserId)
		}

		bodyStr := string(body)

		// 敏感接口脱敏：不记录密码/token/验证码等
		if sensitivePaths[c.Request.URL.Path] {
			bodyStr = "[敏感数据已脱敏]"
		}

		record := model.Operate{
			Ip:     c.ClientIP(),
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
			Agent:  c.Request.UserAgent(),
			Body:   bodyStr,
			UserID: userId,
		}

		// 上传文件时候 中间件日志进行裁断操作
		if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
			if len(record.Body) > bufferSize {
				record.Body = "[文件]"
			}
		}

		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		now := time.Now()

		c.Next()

		latency := time.Since(now)
		record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		record.Status = c.Writer.Status()
		record.Latency = latency

		reg := regexp.MustCompile(`"msg"\s*:\s*"([^"]+)"`)
		match := reg.FindStringSubmatch(writer.body.String())
		if len(match) == 2 {
			record.Msg = match[1]
		}

		if err := service.AddOperate(&record); err != nil {
			utils.ErrorLog("添加操作记录失败", "middleware", err.Error())
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
