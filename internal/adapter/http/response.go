package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thaian1234/green_light/internal/core/domain"
	"github.com/thaian1234/green_light/pkg/util"
)

type (
	Response struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
		Data    any    `json:"data,omitempty"`
	}
	ErrorResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
		Errors  any    `json:"errors,omitempty"`
	}
	Envelope map[string]any
)

var errorStatusMap = map[error]int{
	domain.ErrInternalServer:     http.StatusInternalServerError,
	domain.ErrDataNotFound:       http.StatusNotFound,
	domain.ErrConflictingData:    http.StatusConflict,
	domain.ErrInvalidCredentials: http.StatusUnauthorized,
	domain.ErrUnauthorized:       http.StatusUnauthorized,
	domain.ErrInvalidToken:       http.StatusUnauthorized,
	domain.ErrExpiredToken:       http.StatusUnauthorized,
	domain.ErrForbidden:          http.StatusForbidden,
	domain.ErrNoUpdatedData:      http.StatusBadRequest,
	domain.ErrorValidation:       http.StatusUnprocessableEntity,
	domain.ErrConflictingData:    http.StatusConflict,
}

func newResponse(message string, data any) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func newErrorResponse(message string, err any) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Message: message,
		Errors:  err,
	}
}

func HandleValidationError(ctx *gin.Context, err error) {
	errMsg := util.ParseError(err)
	errRsp := newErrorResponse("validation error", errMsg)
	ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, errRsp)
}

func HandleError(ctx *gin.Context, err error) {
	errMsg := util.ParseError(err)
	msg := err.Error()
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
		msg = "Internal server error"
	}
	errResponse := newErrorResponse(msg, errMsg)
	ctx.JSON(statusCode, errResponse)
}

func HandleAbort(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	msg := err.Error()
	if !ok {
		statusCode = http.StatusInternalServerError
		msg = "Internal server error"
	}
	errResponse := newErrorResponse(msg, err)
	ctx.AbortWithStatusJSON(statusCode, errResponse)
}

func SendSuccess(ctx *gin.Context, data any) {
	response := newResponse("Operation successful", data)
	ctx.JSON(http.StatusOK, response)
}

func SendCreatedSuccess(ctx *gin.Context, data any) {
	response := newResponse("Resource created successfully", data)
	ctx.JSON(http.StatusCreated, response)
}

func SendDeletedSuccess(ctx *gin.Context) {
	response := newResponse("Resource deleted successfully", nil)
	ctx.JSON(http.StatusOK, response)
}

func SendUpdatedSuccess(ctx *gin.Context, data any) {
	response := newResponse("Resource updated successfully", data)
	ctx.JSON(http.StatusOK, response)
}
