package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR   = 500
	SUCCESS = 200
)

func Result(c *gin.Context, code int, data interface{}, msg string) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(c, SUCCESS, map[string]interface{}{}, "ok")
}

func OkWithMessage(c *gin.Context, message string) {
	Result(c, SUCCESS, map[string]interface{}{}, message)
}

func OkWithData(c *gin.Context, data interface{}) {
	Result(c, SUCCESS, data, "ok")
}

func OkWithDetailed(c *gin.Context, data interface{}, message string) {
	Result(c, SUCCESS, data, message)
}

func OkWithString(c *gin.Context, message string) {
	c.String(SUCCESS, message)
}

func Fail(c *gin.Context) {
	Result(c, ERROR, map[string]interface{}{}, "ok")
}

func FailWithMessage(c *gin.Context, message string) {
	Result(c, ERROR, map[string]interface{}{}, message)
}

func FailWithDetailed(c *gin.Context, data interface{}, message string) {
	Result(c, ERROR, data, message)
}

func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{Code: http.StatusForbidden, Data: nil, Msg: ""})
}

func ForbiddenWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{Code: http.StatusForbidden, Data: nil, Msg: message})
}
