package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdaptHandler 将标准的 http.HandlerFunc 适配为 gin.HandlerFunc
// 用于逐步迁移过程中，允许新旧处理函数共存
func AdaptHandler(handler http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c.Writer, c.Request)
	}
}

// AdaptHandlerWithNext 适配器，支持在 Gin 中间件链中使用
func AdaptHandlerWithNext(handler http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 执行适配的处理函数
		handler(c.Writer, c.Request)

		// 如果处理函数没有终止请求，继续执行下一个中间件/处理函数
		if !c.IsAborted() && c.Writer.Status() < 300 {
			c.Next()
		}
	}
}

// ResponseWrapper 包装响应，确保响应格式一致
type ResponseWrapper struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// SuccessResponse 创建成功响应
func SuccessResponse(data interface{}) ResponseWrapper {
	return ResponseWrapper{
		Success: true,
		Data:    data,
	}
}

// ErrorResponse 创建错误响应
func ErrorResponse(err error, message ...string) ResponseWrapper {
	msg := err.Error()
	if len(message) > 0 {
		msg = message[0]
	}
	return ResponseWrapper{
		Success: false,
		Error:   err.Error(),
		Message: msg,
	}
}

// JSONResponse 统一的 JSON 响应辅助函数
func JSONResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

// SuccessJSON 成功的 JSON 响应
func SuccessJSON(c *gin.Context, data interface{}) {
	JSONResponse(c, http.StatusOK, SuccessResponse(data))
}

// ErrorJSON 错误的 JSON 响应
func ErrorJSON(c *gin.Context, statusCode int, err error, message ...string) {
	JSONResponse(c, statusCode, ErrorResponse(err, message...))
}

// MethodNotAllowed 统一的 Method Not Allowed 响应
func MethodNotAllowed(c *gin.Context) {
	ErrorJSON(c, http.StatusMethodNotAllowed,
		http.ErrNotSupported, "Method not allowed")
}

// BadRequest 统一的 Bad Request 响应
func BadRequest(c *gin.Context, err error) {
	ErrorJSON(c, http.StatusBadRequest, err)
}

// InternalServerError 统一的 Internal Server Error 响应
func InternalServerError(c *gin.Context, err error) {
	ErrorJSON(c, http.StatusInternalServerError, err)
}