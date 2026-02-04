package config

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(databaseURL string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("failed to parse database url: %v", err)
	}

	// Pool settings
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 1 * time.Minute
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	return pool
}
