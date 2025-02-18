package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thaian1234/green_light/config"
	"github.com/thaian1234/green_light/internal/adapter/http/handlers"
	"github.com/thaian1234/green_light/internal/adapter/http/middlewares"
	"github.com/thaian1234/green_light/internal/adapter/storages/postgres"
	"github.com/thaian1234/green_light/internal/adapter/storages/postgres/repository"
	"github.com/thaian1234/green_light/internal/core/services"
	"github.com/thaian1234/green_light/pkg/logger"
	"github.com/thaian1234/green_light/pkg/util"
)

type Adapter struct {
	cfg *config.Config
	srv *http.Server
	db  *postgres.Adapter
	wg  *sync.WaitGroup
}

func NewAdapter(cfg *config.Config, db *postgres.Adapter, wg *sync.WaitGroup) *Adapter {
	router := gin.Default()

	// Middlewares
	router.Use(middlewares.RateLimit(cfg.Limiter))

	// Custom Validator
	validator := util.NewValidator()
	validator.SetupValidator()

	// repositories
	movieRepo := repository.NewMovieRepository(db.Pool)
	userRepo := repository.NewUserRepository(db.Pool)

	// services
	healthSvc := services.NewHealthService(cfg)
	movieSvc := services.NewMovieService(movieRepo)
	userSvc := services.NewUserService(userRepo)
	mailerSvc := services.NewMailerService(cfg.Smtp)

	// Handlers
	healthHandler := handlers.NewHealthHandler(healthSvc)
	movieHandler := handlers.NewMovieHandler(movieSvc)
	userHandler := handlers.NewUserHandler(wg, userSvc, mailerSvc)

	// Routes
	_, err := NewRoutes(
		router,
		cfg,
		healthHandler,
		movieHandler,
		userHandler,
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
		db:  db,
		wg:  wg,
	}
}

func (a *Adapter) Run() error {
	logger.Info("Server is running on port::" + a.cfg.HTTP.Port)
	if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Failed to run server", err)
		return err
	}
	return nil
}

func (a *Adapter) Stop(ctx context.Context) error {
	return a.srv.Shutdown(ctx)
}
