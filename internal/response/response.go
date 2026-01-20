package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SUCCESS = 0
	ERROR   = 7
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}, msg ...string) {
	r := Response{
		Code: SUCCESS,
		Msg:  "操作成功",
		Data: data,
	}
	if len(msg) > 0 {
		r.Msg = msg[0]
	}
	c.JSON(http.StatusOK, r)
}

// Fail 失败响应
func Fail(c *gin.Context, msg string, data ...interface{}) {
	r := Response{
		Code: ERROR,
		Msg:  msg,
		Data: nil,
	}
	if len(data) > 0 {
		r.Data = data[0]
	}
	c.JSON(http.StatusOK, r)
}

// FailWithErrorCode 失败响应（带自定义错误码）
func FailWithErrorCode(c *gin.Context, code int, msg string, data ...interface{}) {
	r := Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	if len(data) > 0 {
		r.Data = data[0]
	}
	c.JSON(http.StatusOK, r)
}

// FailWithError 错误响应
func FailWithError(c *gin.Context, err error) {
	Fail(c, err.Error())
}

// TokenError token验证失败响应
func TokenError(c *gin.Context) {
	FailWithErrorCode(c, ERROR, "token验证失败")
}
