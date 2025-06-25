package store

import (
	"context"

	"github.com/dorik33/DeNet/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnection(cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseCfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
