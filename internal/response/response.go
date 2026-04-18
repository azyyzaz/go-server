package response

import "github.com/gin-gonic/gin"

const (
	CodeOK = "OK"
)

type Body struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	TraceID   string `json:"trace_id,omitempty"`
}

type AppError struct {
	Status  int
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

func Success(c *gin.Context, data any) {
	c.JSON(200, Body{
		Code:      CodeOK,
		Message:   "success",
		Data:      data,
		RequestID: RequestIDFromContext(c),
		TraceID:   TraceIDFromContext(c),
	})
}

func Created(c *gin.Context, data any) {
	c.JSON(201, Body{
		Code:      CodeOK,
		Message:   "created",
		Data:      data,
		RequestID: RequestIDFromContext(c),
		TraceID:   TraceIDFromContext(c),
	})
}

func Fail(c *gin.Context, status int, code, message string) {
	c.JSON(status, Body{
		Code:      code,
		Message:   message,
		RequestID: RequestIDFromContext(c),
		TraceID:   TraceIDFromContext(c),
	})
}

func RequestIDFromContext(c *gin.Context) string {
	if rid, ok := c.Get("request_id"); ok {
		if v, ok := rid.(string); ok {
			return v
		}
	}
	return ""
}

func TraceIDFromContext(c *gin.Context) string {
	if tid, ok := c.Get("trace_id"); ok {
		if v, ok := tid.(string); ok {
			return v
		}
	}
	return ""
}
