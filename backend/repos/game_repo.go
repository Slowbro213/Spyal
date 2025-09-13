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
	_ = cache.Set(ctx, g.TableName() + g.RoomID, string(data), time.Hour)
	return nil
}
