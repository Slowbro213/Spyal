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

var ErrRoundNotFound = errors.New("round not found")

type RoundRepository interface {
	Create(ctx context.Context, r *models.Round) error
	GetByID(ctx context.Context, id int64) (*models.Round, error)
	GetLatestByGame(ctx context.Context, gameID int64) (*models.Round, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}

type roundRepo struct {
	db sqlx.ExtContext
}

func NewRoundRepo(db sqlx.ExtContext) RoundRepository {
	return &roundRepo{db: db}
}

func (r *roundRepo) Create(ctx context.Context, rnd *models.Round) error {
	query := `
		INSERT INTO rounds (game_id, status, word, spy_word)
		VALUES (:game_id, :status, :word, :spy_word)
		RETURNING id, created_at`
	rows, err := sqlx.NamedQueryContext(ctx, r.db, query, rnd)
	if err != nil {
		return fmt.Errorf("create round: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&rnd.ID, &rnd.CreatedAt); err != nil {
			return fmt.Errorf("scan returning: %w", err)
		}
	}

	data, _ := json.Marshal(rnd)
	_ = cache.Set(ctx, rnd.CacheKey(), string(data), time.Hour)
	return nil
}

func (r *roundRepo) GetByID(ctx context.Context, id int64) (*models.Round, error) {
	key := fmt.Sprintf("round_%d", id)

	if cached, err := cache.Get(ctx, key); err == nil && cached != "" {
		var rnd models.Round
		if jsonErr := json.Unmarshal([]byte(cached), &rnd); jsonErr == nil {
			return &rnd, nil
		}
	}

	var rnd models.Round
	err := sqlx.GetContext(ctx, r.db, &rnd,
		`SELECT id, game_id, status, word, spy_word, created_at
		 FROM rounds WHERE id=$1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoundNotFound
		}
		return nil, fmt.Errorf("get round by id: %w", err)
	}

	data, _ := json.Marshal(rnd)
	_ = cache.Set(ctx, key, string(data), time.Hour)
	return &rnd, nil
}

func (r *roundRepo) GetLatestByGame(ctx context.Context, gameID int64) (*models.Round, error) {
	var rnd models.Round
	err := sqlx.GetContext(ctx, r.db, &rnd,
		`SELECT * FROM rounds
		 WHERE game_id=$1
		 ORDER BY created_at DESC
		 LIMIT 1`, gameID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoundNotFound
		}
		return nil, fmt.Errorf("get latest round: %w", err)
	}
	return &rnd, nil
}

func (r *roundRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	if _, err := r.db.ExecContext(ctx,
		`UPDATE rounds SET status=$1 WHERE id=$2`, status, id); err != nil {
		return fmt.Errorf("update round status: %w", err)
	}
	key := fmt.Sprintf("round_%d", id)
	_ = cache.Delete(ctx, key)
	return nil
}
