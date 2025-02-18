package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/thaian1234/green_light/config"
	"github.com/thaian1234/green_light/internal/adapter/http"
	"github.com/thaian1234/green_light/internal/adapter/storages/postgres"
	"github.com/thaian1234/green_light/pkg/logger"
)

func main() {
	ctx := context.Background()
	var wg sync.WaitGroup

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

	// Pass wg to your adapters/services that need it
	httpAdapter := http.NewAdapter(cfg, dbAdapter, &wg)
	errChan := make(chan error)

	go func() {
		errChan <- httpAdapter.Run()
	}()

	// Wait for server error or interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Printf("Server error: %v", err)
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
	case <-ctx.Done():
		log.Println("Context cancelled")
	}

	log.Println("Shutting down server...")
	httpAdapter.Stop(ctx)

	// Wait for all background tasks to complete
	log.Println("Waiting for background tasks to complete...")
	wg.Wait()
	log.Println("Server shutdown completed")
}
