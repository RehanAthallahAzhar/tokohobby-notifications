package configs

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func NewDatabaseConnection(cfg *DatabaseConfig, log *logrus.Logger) (*pgxpool.Pool, error) {
	dsn := cfg.DSN()

	log.WithField("dsn", maskPassword(dsn)).Info("Connecting to database...")

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Database connection established")
	return pool, nil
}

func maskPassword(dsn string) string {
	// Simple masking for logging
	return "postgres://***:***@..."
}
