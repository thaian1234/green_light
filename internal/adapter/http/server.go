package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thaian1234/green_light/config"
	"github.com/thaian1234/green_light/internal/adapter/http/handlers"
	"github.com/thaian1234/green_light/internal/adapter/logger"
	"github.com/thaian1234/green_light/internal/core/services"
)

type Adapter struct {
	cfg *config.Config
	srv *http.Server
}

func NewAdapter(cfg *config.Config) *Adapter {
	router := gin.Default()

	// services
	healthSvc := services.NewHealthService(cfg)
	movieSvc := services.NewMovieService()

	// Handlers
	healthHandler := handlers.NewHealthHandler(healthSvc)
	movieHandler := handlers.NewMovieHandler(movieSvc)

	// Routes
	_, err := NewRoutes(
		router,
		cfg,
		healthHandler,
		movieHandler,
	)

	if err != nil {
		logger.Fatal("failed to setup routes ", err)
	}

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.HTTP.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &Adapter{
		cfg: cfg,
		srv: srv,
	}
}

func (a *Adapter) Run() {
	log.Printf("Server is running on port::%s", a.cfg.HTTP.Port)
	if err := a.srv.ListenAndServe(); err != nil {
		msg := fmt.Sprintf("failed to run Server on port::%s", a.cfg.HTTP.Port)
		logger.Fatal(msg, err)
	}
}

func (a *Adapter) Stop(ctx context.Context) {
	a.srv.Shutdown(ctx)
}
