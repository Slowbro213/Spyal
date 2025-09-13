package repos

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"spyal/cache"
	"spyal/core"
	"spyal/models"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var ErrNotFound = errors.New("game not found")

type repoInterface interface {
	GetBy(context.Context, models.Model, string, any) error
}

type repo struct {
	db sqlx.ExtContext
}


func (b repo) GetBy(
	ctx context.Context,
	dest models.Model,
	column string,
	value any,
) error {
	key := fmt.Sprintf("%s_%s_%v", dest.TableName(), column, value)

	if cached, err := cache.Get(ctx, key); err == nil && cached != "" {
		if err := json.Unmarshal([]byte(cached), dest); err == nil {
			return nil
		}
		core.Logger.Warn("cache unmarshal failed, falling back to DB",  zap.Error(err))
	}

	q := fmt.Sprintf(`SELECT * FROM %s WHERE %s = $1 LIMIT 1`, dest.TableName(), column)
	if err := sqlx.GetContext(ctx, b.db, dest, q, value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return fmt.Errorf("db query failed: %w", err)
	}

	if data, err := json.Marshal(dest); err == nil {
		if err := cache.Set(ctx, key, string(data), time.Hour); err != nil {
			core.Logger.Warn("failed to set cache", zap.Error(err))
		}
	} else {
		core.Logger.Warn("failed to marshal for cache", zap.Error(err))
	}

	return nil
}


func (b repo) Exists(ctx context.Context, table models.Model, column string, value any) (bool, error) {
	q := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE %s = $1)`, table.TableName(), column)
	var ok bool
	err := sqlx.GetContext(ctx, b.db, &ok, q, value)
	return ok, err
}

func (b repo) Count(ctx context.Context, table models.Model, column string, value any) (int64, error) {
	q := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s = $1`, table.TableName(), column)
	var n int64
	err := sqlx.GetContext(ctx, b.db, &n, q, value)
	return n, err
}

func (b repo) DeleteBy(ctx context.Context, table models.Model, column string, value any) (int64, error) {
	q := fmt.Sprintf(`DELETE FROM %s WHERE %s = $1`, table.TableName(), column)
	res, err := b.db.ExecContext(ctx, q, value)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (b repo) UpdateColumns(
	ctx context.Context,
	table models.Model,
	idCol string, idValue any,
	cols map[string]any,
) error {
	if len(cols) == 0 {
		return nil
	}
	set := make([]string, 0, len(cols))
	args := []any{idValue}
	i := 2
	for k, v := range cols {
		set = append(set, fmt.Sprintf("%s = $%d", k, i))
		args = append(args, v)
		i++
	}
	q := fmt.Sprintf(`UPDATE %s SET %s WHERE %s = $1`,
		table.TableName(), strings.Join(set, ", "), idCol)
	_, err := b.db.ExecContext(ctx, q, args...)
	return err
}
