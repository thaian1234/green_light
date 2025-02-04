package handlers

import (
	"time"

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
	resp := Envelope{
		"status": status.Status,
		"system_info": map[string]string{
			"environment": status.Environment,
			"version":     status.Version,
		},
	}
	time.Sleep(4 * time.Second)
	SendSuccess(ctx, resp)
}
