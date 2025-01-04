package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/thaian1234/green_light/internal/core/ports"
	"github.com/thaian1234/green_light/pkg/response"
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
	resp := response.Envelope{
		"status": status.Status,
		"system_info": map[string]string{
			"environment": status.Environment,
			"version":     status.Version,
		},
	}
	response.SendSuccess(ctx, resp)
}
