package repos

import (
	"errors"
	"fmt"
	"time"
	"encoding/json"
	"context"

	"github.com/jmoiron/sqlx"
	"spyal/models"
	"spyal/cache"
)

var ErrGameNotFound = errors.New("game not found")

type GameRepository interface {
	repoInterface
	Create(context.Context,*models.Game) error
	GetPublicActive(context.Context,string) ([]*models.Game, error)
}

type gameRepo struct {
	repo
}

func NewGameRepo(db sqlx.ExtContext) GameRepository {
	return &gameRepo{repo: repo{db: db}}
}

func (r *gameRepo) Create(ctx context.Context, g *models.Game) error {
	query := `
		INSERT INTO games (host_id, room_id, spy_number, max_players, name, private)
		VALUES (:host_id, :room_id, :spy_number, :max_players, :name, :private)
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
	_ = cache.Set(ctx, g.TableName()+g.RoomID, string(data), time.Hour)
	_ = cache.Delete(ctx, "public_active")
	return nil
}

func (r *gameRepo) GetPublicActive(ctx context.Context, searchTerm string) ([]*models.Game, error) {
	key := "public_active"
	if searchTerm != "" {
		key += ":search:" + searchTerm
	}

	if val, err := cache.Get(ctx, key); err == nil && val != "" {
		var cached []*models.Game
		if json.Unmarshal([]byte(val), &cached) == nil {
			return cached, nil
		}
	}

	var games []*models.Game
	sb := `SELECT * FROM games WHERE private = false AND active = true`
	if searchTerm != "" {
		sb += ` AND (name ILIKE '%' || $1 || '%' OR room_id ILIKE '%' || $1 || '%')`
	}
	sb += ` ORDER BY created_at DESC`

	if searchTerm != "" {
		if err := sqlx.SelectContext(ctx, r.db, &games, sb, searchTerm); err != nil {
			return nil, fmt.Errorf("get public active games: %w", err)
		}
	} else {
		if err := sqlx.SelectContext(ctx, r.db, &games, sb); err != nil {
			return nil, fmt.Errorf("get public active games: %w", err)
		}
	}

	data, _ := json.Marshal(games)
	//nolint
	_ = cache.Set(ctx, key, string(data), 5*time.Minute)
	return games, nil
}
