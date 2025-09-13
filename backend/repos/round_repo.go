package repos

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"spyal/cache"
	"spyal/core"
	"spyal/models"
)

var ErrRoundNotFound = errors.New("round not found")

type RoundRepository interface {
	repoInterface
	Create(ctx context.Context, r *models.Round) error
	GetLatestByGame(ctx context.Context, gameID int64) (*models.Round, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	AddParticipant(ctx context.Context, roundID, userID int64, isSpy bool) error
	GetPlayers(ctx context.Context, roundID int64) ([]*models.RoundPlayer, error)
}

type roundRepo struct {
	repo
}

func NewRoundRepo(db sqlx.ExtContext) RoundRepository {
	return &roundRepo{repo: repo{db: db}}
}

func (r *roundRepo) Create(ctx context.Context, rnd *models.Round) error {
	query := `
		INSERT INTO rounds (game_id, word_id, spy_word_id)
		VALUES (:game_id, :word_id, :spy_word_id)
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
	_ = cache.Set(ctx, fmt.Sprintf("rounds_%d",rnd.ID), string(data), time.Hour)
	return nil
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

func (r *roundRepo) AddParticipant(
	ctx context.Context,
	roundID, userID int64,
	isSpy bool,
) error {
	const q = `
		INSERT INTO game_participants (round_id, user_id, is_spy)
		VALUES ($1, $2, $3)
		ON CONFLICT (round_id, user_id) DO NOTHING`
	if _, err := r.db.ExecContext(ctx, q, roundID, userID, isSpy); err != nil {
		return fmt.Errorf("add participant: %w", err)
	}
	_ = cache.Delete(ctx, fmt.Sprintf("players:%d", roundID))
	return nil
}

func (r *roundRepo) GetPlayers(ctx context.Context, roundID int64) ([]*models.RoundPlayer, error) {
	key := fmt.Sprintf("players:%d", roundID)

	if cached, err := cache.Get(ctx, key); err == nil && cached != "" {
		var out []*models.RoundPlayer
		if err := json.Unmarshal([]byte(cached), &out); err == nil {
			return out, nil
		}
		core.Logger.Warn("players cache unmarshal failed", zap.Error(err))
	}

	var players []*models.RoundPlayer
	if err := sqlx.SelectContext(ctx, r.db, &players, `
		SELECT user_id, is_spy
		FROM game_participants
		WHERE round_id = $1
		ORDER BY joined_at`, roundID); err != nil {
		return nil, fmt.Errorf("get players: %w", err)
	}

	if data, err := json.Marshal(players); err == nil {
		_ = cache.Set(ctx, key, string(data), time.Hour)
	}
	return players, nil
}
