package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
}

func NewAdapter(cfg *config.Config, db *postgres.Adapter) *Adapter {
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

	// Handlers
	healthHandler := handlers.NewHealthHandler(healthSvc)
	movieHandler := handlers.NewMovieHandler(movieSvc)
	userHandler := handlers.NewUserHandler(userSvc)

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
	}
}

func (a *Adapter) Run(ctx context.Context) error {
	// Start server in a goroutine
	go func() {
		log.Printf("Server is running on port::%s", a.cfg.HTTP.Port)
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			msg := fmt.Sprintf("failed to run Server on port::%s", a.cfg.HTTP.Port)
			logger.Fatal(msg, err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal
	sig := <-quit
	logger.Info("Received shutdown signal", "signal", sig.String())

	// Create shutdown context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	return a.srv.Shutdown(ctxWithTimeout)
}

func (a *Adapter) Stop(ctx context.Context) {
	a.srv.Shutdown(ctx)
}
