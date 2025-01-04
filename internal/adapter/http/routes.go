package http

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thaian1234/green_light/config"
	"github.com/thaian1234/green_light/internal/adapter/http/handlers"
)

type Routes struct {
	r   *gin.Engine
	cfg *config.Config
}

func NewRoutes(r *gin.Engine, cfg *config.Config, healthHandler *handlers.HealthHandler) (*Routes, error) {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	ginConfig := cors.DefaultConfig()
	allowedOrigins := cfg.HTTP.AllowedOrigins
	originsList := strings.Split(allowedOrigins, ",")
	ginConfig.AllowOrigins = originsList

	r.Use(cors.New(ginConfig))

	v1 := r.Group("/v1/api")
	{
		health := v1.Group("/health")
		{
			health.GET("/", healthHandler.Check)
		}
	}

	return &Routes{
		r:   r,
		cfg: cfg,
	}, nil
}
