package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thaian1234/green_light/internal/adapter/logger"
	"github.com/thaian1234/green_light/internal/core/util"
)

type response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type errorResponse struct {
	Success  bool     `json:"success"`
	Messages []string `json:"messages"`
}

func newResponse(success bool, message string, data any) response {
	return response{
		Success: success,
		Message: message,
		Data:    data,
	}
}

func newErrorResponse(errMsgs []string) errorResponse {
	return errorResponse{
		Success:  false,
		Messages: errMsgs,
	}
}

func validationError(ctx *gin.Context, err error) {
	errMsgs := util.ParseError(err)
	errRsp := newErrorResponse(errMsgs)
	logger.Error("Validation Error", "error", errMsgs)
	ctx.JSON(http.StatusUnprocessableEntity, errRsp)
}

func handleSuccess(ctx *gin.Context, httpStatus int, data any) {
	rsp := newResponse(true, "Success", data)
	ctx.JSON(httpStatus, rsp)
}

func handleError(ctx *gin.Context, err error) {
	errMsg := util.ParseError(err)
	errResp := newErrorResponse(errMsg)
	logger.Error("HandleError Failed", "error", err)
	ctx.JSON(http.StatusBadRequest, errResp)
}
