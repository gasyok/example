package pgxpool

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool = pgxpool.Pool

type Config struct {
	Host     string
	Port     string
	UserName string
	Password string
	DBName   string
}

func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)
	return NewPoolFromDSN(ctx, dsn)
}

func NewPoolFromDSN(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pool.Ping: %w", err)
	}

	return pool, nil
}
