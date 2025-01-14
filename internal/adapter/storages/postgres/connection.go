package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thaian1234/green_light/config"
)

type Adapter struct {
	*pgxpool.Pool
}

func NewAdapter(ctx context.Context, cfg *config.DB) (*Adapter, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	dbConfig, err := pgxpool.ParseConfig(dsn)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if err != nil {
		return nil, fmt.Errorf("error parsing config: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(ctxWithTimeout, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	if err := pool.Ping(ctxWithTimeout); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return &Adapter{
		Pool: pool,
	}, nil
}
