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
	ErrorResponse struct {
		Success bool   `json:"success"`
		Code    string `json:"code"`
		Message string `json:"message"`
		Errors  any    `json:"errors,omitempty"`
	}
	Code     string
	Envelope map[string]any
)

const (
	Success          Code = "SUCCESS"
	Created          Code = "CREATED"
	Deleted          Code = "DELETED"
	Updated          Code = "UPDATED"
	BadRequest       Code = "BAD_REQUEST"
	Unauthorized     Code = "UNAUTHORIZED"
	Unauthenticated  Code = "UNAUTHENTICATED"
	Forbidden        Code = "FORBIDDEN"
	Internal         Code = "INTERNAL"
	ValidationError  Code = "VALIDATION_ERROR"
	NotFoundError    Code = "NOT_FOUND"
	ServerError      Code = "SERVER_ERROR"
	MethodNotAllowed Code = "METHOD_NOT_ALLOWED"
)

func newResponse(code Code, message string, data any) Response {
	return Response{
		Success: code == Success,
		Code:    string(code),
		Message: message,
		Data:    data,
	}
}

func newErrorResponse(code Code, message string, err any) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Code:    string(code),
		Message: message,
		Errors:  err,
	}
}

func SendSuccess(ctx *gin.Context, data any) {
	response := newResponse(Success, "Operation successful", data)
	ctx.JSON(http.StatusOK, response)
}

func SendCreatedSuccess(ctx *gin.Context, data any) {
	response := newResponse(Created, "Resource created successfully", data)
	ctx.JSON(http.StatusCreated, response)
}

func SendDeletedSuccess(ctx *gin.Context) {
	response := newResponse(Deleted, "Resource deleted successfully", nil)
	ctx.JSON(http.StatusOK, response)
}

func SendUpdatedSuccess(ctx *gin.Context) {
	response := newResponse(Updated, "Resource updated successfully", nil)
	ctx.JSON(http.StatusOK, response)
}

func SendBadRequest(ctx *gin.Context, err error) {
	errMsgs := util.ParseError(err)
	response := newErrorResponse(BadRequest, "Bad request", errMsgs)
	logger.Error("Bad Request", "error", errMsgs)
	ctx.JSON(http.StatusBadRequest, response)
}

func SendAuthorized(ctx *gin.Context, err error) {
	errMsgs := util.ParseError(err)
	response := newErrorResponse(Unauthorized, "Unauthorized", errMsgs)
	logger.Error("Unauthorized", "error", errMsgs)
	ctx.JSON(http.StatusUnauthorized, response)
}

func SendUnauthenticated(ctx *gin.Context, err error) {
	errMsgs := util.ParseError(err)
	response := newErrorResponse(Unauthenticated, "Unauthenticated", errMsgs)
	logger.Error("Unauthenticated", "error", errMsgs)
	ctx.JSON(http.StatusUnauthorized, response)
}

func SendValidationError(ctx *gin.Context, err error) {
	errMsgs := util.ParseError(err)
	response := newErrorResponse(ValidationError, "Validation failed", errMsgs)
	logger.Error("Validation Error", "error", errMsgs)
	ctx.JSON(http.StatusUnprocessableEntity, response)
}

func SendServerError(ctx *gin.Context, err error) {
	errMsgs := util.ParseError(err)
	response := newErrorResponse(ServerError, "Internal server error", nil)
	logger.Error("Server Error", "error", errMsgs)
	ctx.JSON(http.StatusInternalServerError, response)
}

func SendNotFound(ctx *gin.Context, err error) {
	errMsgs := util.ParseError(err)
	response := newErrorResponse(NotFoundError, "Resource not found", nil)
	logger.Error("Not Found", "error", errMsgs)
	ctx.JSON(http.StatusNotFound, response)
}

func SendMethodNotAllowed(ctx *gin.Context) {
	response := newErrorResponse(MethodNotAllowed, "Method not allowed", nil)
	logger.Error("Method Not Allowed", "method", ctx.Request.Method)
	ctx.JSON(http.StatusMethodNotAllowed, response)
}
