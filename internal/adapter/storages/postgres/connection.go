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

func NewAdapter(cfg *config.DB) (*Adapter, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s pool_max_conns=25 pool_max_conn_lifetime=1h30m pool_max_conn_idle_time=15m", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	dbConfig, err := pgxpool.ParseConfig(dsn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err != nil {
		return nil, fmt.Errorf("error parsing config: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return &Adapter{}, nil
}
