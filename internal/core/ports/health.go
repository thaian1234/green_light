package ports

import "github.com/thaian1234/green_light/internal/core/domain"

type HealthService interface {
	GetStatus() domain.HealthStatus
}
