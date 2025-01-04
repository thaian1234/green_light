package services

import (
	"github.com/thaian1234/green_light/config"
	"github.com/thaian1234/green_light/internal/core/domain"
)

type HealthService struct {
	cfg *config.Config
}

func NewHealthService(cfg *config.Config) *HealthService {
	return &HealthService{
		cfg: cfg,
	}
}

func (s *HealthService) GetStatus() domain.HealthStatus {
	return domain.HealthStatus{
		Status:      "available",
		Environment: s.cfg.App.Env,
		Version:     s.cfg.App.Version,
	}
}
