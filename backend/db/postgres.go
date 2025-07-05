package db

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)


//nolint:gochecknoglobals
var (
	pgInstance *pgxpool.Pool
	pgOnce     sync.Once
)

// GetPostgres returns the singleton Postgres connection pool.
func GetPostgres() (*pgxpool.Pool, error) {
	var err error
	pgOnce.Do(func() {
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			err = errors.New("DATABASE_URL not set")
			return
		}

		pgInstance, err = pgxpool.New(context.Background(), dsn)
	})

	return pgInstance, err
}
