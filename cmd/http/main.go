package main

import (
	"log"

	"github.com/thaian1234/green_light/config"
	"github.com/thaian1234/green_light/internal/adapter/http"
	"github.com/thaian1234/green_light/internal/adapter/storages/postgres"
	"github.com/thaian1234/green_light/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	err = logger.Initialize(cfg.Logger)
	if err != nil {
		log.Fatalf("failed to load logger: %v", err)

	}

	dbAdapter, err := postgres.NewAdapter(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbAdapter.Close()

	httpAdapter := http.NewAdapter(cfg, dbAdapter)
	httpAdapter.Run()
}
