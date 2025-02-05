package main

import (
	"context"
	"log"

	"github.com/thaian1234/green_light/config"
	"github.com/thaian1234/green_light/internal/adapter/http"
	"github.com/thaian1234/green_light/internal/adapter/storages/postgres"
	"github.com/thaian1234/green_light/pkg/logger"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	err = logger.Initialize(cfg.Logger)
	if err != nil {
		log.Fatalf("failed to load logger: %v", err)

	}

	dbAdapter, err := postgres.NewAdapter(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbAdapter.Close()

	httpAdapter := http.NewAdapter(cfg, dbAdapter)
	// Create a channel to listen for errors
	errChan := make(chan error)

	go func() {
		errChan <- httpAdapter.Run(ctx)
	}()

	// Wait for server error or interrupt signal
	select {
	case err := <-errChan:
		log.Printf("Server error: %v", err)
	case <-ctx.Done():
		log.Println("Shutting down server...")
		httpAdapter.Stop(ctx)
	}
}
