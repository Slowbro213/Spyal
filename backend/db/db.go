package db

import (
	"context"
	"fmt"
	"time"
	"errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
)

type DB struct {
	*sqlx.DB
}

type Config struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func Connect(cfg Config) (*DB, error) {
	if cfg.DSN == "" {
		return nil, errors.New("empty DSN")
	}

	sqlxDB, err := sqlx.Connect("postgres", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	if cfg.MaxOpenConns > 0 {
		sqlxDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlxDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlxDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		sqlxDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}

	return &DB{DB: sqlxDB}, nil
}

func (d *DB) Close() error {
	if d == nil || d.DB == nil {
		return nil
	}
	return d.DB.Close()
}

func (d *DB) Health(ctx context.Context) error {
	return d.PingContext(ctx)
}

func (d *DB) WithTx(ctx context.Context, fn func(tx *sqlx.Tx) error) (err error) {
	tx, err := d.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
