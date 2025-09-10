package repos

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"spyal/cache"
	"spyal/models"
)

var ErrGameNotFound = errors.New("game not found")

type GameRepository interface {
	Create(ctx context.Context, g *models.Game) error
	GetByID(ctx context.Context, id int64) (*models.Game, error)
	ListPublic(ctx context.Context, limit int) ([]*models.Game, error)
}

type gameRepo struct {
	db sqlx.ExtContext
}

func NewGameRepo(db sqlx.ExtContext) GameRepository {
	return &gameRepo{db: db}
}

func (r *gameRepo) Create(ctx context.Context, g *models.Game) error {
	query := `
		INSERT INTO games (host_id, title, private, status)
		VALUES (:host_id, :title, :private, :status)
		RETURNING id, created_at`
	rows, err := sqlx.NamedQueryContext(ctx, r.db, query, g)
	if err != nil {
		return fmt.Errorf("create game: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&g.ID, &g.CreatedAt); err != nil {
			return fmt.Errorf("scan returning: %w", err)
		}
	}

	data, _ := json.Marshal(g)
	_ = cache.Set(ctx, g.CacheKey(), string(data), time.Hour)
	return nil
}

func (r *gameRepo) GetByID(ctx context.Context, id int64) (*models.Game, error) {
	key := fmt.Sprintf("game_%d", id)

	if cached, err := cache.Get(ctx, key); err == nil && cached != "" {
		var g models.Game
		if jsonErr := json.Unmarshal([]byte(cached), &g); jsonErr == nil {
			return &g, nil
		}
	}

	var g models.Game
	err := sqlx.GetContext(ctx, r.db, &g,
		`SELECT id, host_id, title, created_at, private, status FROM games WHERE id=$1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGameNotFound
		}
		return nil, fmt.Errorf("get game by id: %w", err)
	}

	data, _ := json.Marshal(g)
	_ = cache.Set(ctx, key, string(data), time.Hour)
	return &g, nil
}

func (r *gameRepo) ListPublic(ctx context.Context, limit int) ([]*models.Game, error) {
	var gg []*models.Game
	err := sqlx.SelectContext(ctx, r.db, &gg,
		`SELECT * FROM games WHERE private=false ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("list public games: %w", err)
	}
	return gg, nil
}
