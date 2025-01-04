package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thaian1234/green_light/internal/core/ports"
)

type HealthHandler struct {
	healthSvc ports.HealthService
}

func NewHealthHandler(healSvc ports.HealthService) *HealthHandler {
	return &HealthHandler{
		healthSvc: healSvc,
	}
}

func (h *HealthHandler) Check(ctx *gin.Context) {
	status := h.healthSvc.GetStatus()
	handleSuccess(ctx, http.StatusOK, status)
}
