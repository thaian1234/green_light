package handlers

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
}

var errorStatusMsg = map[error]string{
	domain.ErrInternalServer:     "Internal Server Error",
	domain.ErrDataNotFound:       "Data Not Found",
	domain.ErrConflictingData:    "Conflicting Data",
	domain.ErrInvalidCredentials: "Invalid Credentials",
	domain.ErrUnauthorized:       "Unauthorized",
	domain.ErrInvalidToken:       "Invalid Token",
	domain.ErrExpiredToken:       "Token Expired",
	domain.ErrForbidden:          "Forbidden",
	domain.ErrNoUpdatedData:      "No Updated Data",
	domain.ErrorValidation:       "Validation Error",
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
	statusMsg := errorStatusMsg[domain.ErrorValidation]
	errRsp := newErrorResponse(statusMsg, errMsg)
	ctx.JSON(http.StatusBadRequest, errRsp)
}

func HandleError(ctx *gin.Context, err error) {
	errMsg := util.ParseError(err)
	statusMsg := errorStatusMsg[err]

	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}
	errResponse := newErrorResponse(statusMsg, errMsg)
	ctx.JSON(statusCode, errResponse)
}

func HandleAbort(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	statusMsg := errorStatusMsg[err]

	if !ok {
		statusCode = http.StatusInternalServerError
	}
	errResponse := newErrorResponse(statusMsg, err)
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
