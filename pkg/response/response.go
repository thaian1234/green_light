package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thaian1234/green_light/pkg/logger"
	"github.com/thaian1234/green_light/pkg/util"
)

type (
	Response struct {
		Success bool   `json:"success"`
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    any    `json:"data,omitempty"`
	}
	ErrorCode string
	Envelope  map[string]any
)

const (
	Success          ErrorCode = "SUCCESS"
	ValidationError  ErrorCode = "VALIDATION_ERROR"
	NotFoundError    ErrorCode = "NOT_FOUND"
	ServerError      ErrorCode = "SERVER_ERROR"
	MethodNotAllowed ErrorCode = "METHOD_NOT_ALLOWED"
)

func NewResponse(code ErrorCode, message string, data any) Response {
	return Response{
		Success: code == Success,
		Code:    string(code),
		Message: message,
		Data:    data,
	}
}

func SendSuccess(ctx *gin.Context, data any) {
	response := NewResponse(Success, "Operation successful", data)
	ctx.JSON(http.StatusOK, response)
}

func SendValidationError(ctx *gin.Context, err error) {
	errMsgs := util.ParseError(err)
	response := NewResponse(ValidationError, "Validation failed", errMsgs)
	logger.Error("Validation Error", "error", errMsgs)
	ctx.JSON(http.StatusUnprocessableEntity, response)
}

func SendServerError(ctx *gin.Context, err error) {
	errMsgs := util.ParseError(err)
	response := NewResponse(ServerError, "Internal server error", nil)
	logger.Error("Server Error", "error", errMsgs)
	ctx.JSON(http.StatusInternalServerError, response)
}

func SendNotFound(ctx *gin.Context, err error) {
	errMsgs := util.ParseError(err)
	response := NewResponse(NotFoundError, "Resource not found", nil)
	logger.Error("Not Found", "error", errMsgs)
	ctx.JSON(http.StatusNotFound, response)
}

func SendMethodNotAllowed(ctx *gin.Context) {
	response := NewResponse(MethodNotAllowed, "Method not allowed", nil)
	logger.Error("Method Not Allowed", "method", ctx.Request.Method)
	ctx.JSON(http.StatusMethodNotAllowed, response)
}
