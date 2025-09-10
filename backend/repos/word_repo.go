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

var ErrWordNotFound = errors.New("word not found")

type WordRepository interface {
	Create(ctx context.Context, w *models.Word) error
	GetByID(ctx context.Context, id int64) (*models.Word, error)
	GetByWord(ctx context.Context, word string) (*models.Word, error)
	AddRelated(ctx context.Context, wordID1, wordID2 int64) error
	GetRelated(ctx context.Context, wordID int64) ([]*models.Word, error)
}

type wordRepo struct {
	db sqlx.ExtContext
}

func NewWordRepo(db sqlx.ExtContext) WordRepository {
	return &wordRepo{db: db}
}

func (r *wordRepo) Create(ctx context.Context, w *models.Word) error {
	query := `INSERT INTO words (word) VALUES (:word) RETURNING id`
	rows, err := sqlx.NamedQueryContext(ctx, r.db, query, w)
	if err != nil {
		return fmt.Errorf("create word: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&w.ID); err != nil {
			return fmt.Errorf("scan id: %w", err)
		}
	}
	data, _ := json.Marshal(w)
	_ = cache.Set(ctx, w.CacheKey(), string(data), time.Hour)
	return nil
}

func (r *wordRepo) GetByID(ctx context.Context, id int64) (*models.Word, error) {
	key := fmt.Sprintf("word_%d", id)
	if cached, err := cache.Get(ctx, key); err == nil && cached != "" {
		var w models.Word
		if json.Unmarshal([]byte(cached), &w) == nil {
			return &w, nil
		}
	}
	var w models.Word
	err := sqlx.GetContext(ctx, r.db, &w, `SELECT id, word FROM words WHERE id=$1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWordNotFound
		}
		return nil, fmt.Errorf("get word by id: %w", err)
	}
	data, _ := json.Marshal(w)
	_ = cache.Set(ctx, key, string(data), time.Hour)
	return &w, nil
}

func (r *wordRepo) GetByWord(ctx context.Context, word string) (*models.Word, error) {
	key := "word_txt_" + word
	if cached, err := cache.Get(ctx, key); err == nil && cached != "" {
		var w models.Word
		if json.Unmarshal([]byte(cached), &w) == nil {
			return &w, nil
		}
	}
	var w models.Word
	err := sqlx.GetContext(ctx, r.db, &w, `SELECT id, word FROM words WHERE word=$1`, word)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWordNotFound
		}
		return nil, fmt.Errorf("get word by text: %w", err)
	}
	data, _ := json.Marshal(w)
	_ = cache.Set(ctx, key, string(data), time.Hour)
	return &w, nil
}

func (r *wordRepo) AddRelated(ctx context.Context, wordID1, wordID2 int64) error {
	if wordID1 == wordID2 {
		return errors.New("cannot relate word to itself")
	}
	if wordID1 > wordID2 {
		wordID1, wordID2 = wordID2, wordID1
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO word_related (word_id_1, word_id_2) VALUES ($1, $2)
		 ON CONFLICT DO NOTHING`, wordID1, wordID2)
	if err != nil {
		return fmt.Errorf("add related: %w", err)
	}
	return nil
}

func (r *wordRepo) GetRelated(ctx context.Context, wordID int64) ([]*models.Word, error) {
	var ww []*models.Word
	err := sqlx.SelectContext(ctx, r.db, &ww, `
		SELECT w.id, w.word
		FROM word_related wr
		JOIN words w ON w.id = CASE
			WHEN wr.word_id_1 = $1 THEN wr.word_id_2
			ELSE wr.word_id_1
		END
		WHERE $1 IN (wr.word_id_1, wr.word_id_2)`, wordID)
	if err != nil {
		return nil, fmt.Errorf("get related: %w", err)
	}
	return ww, nil
}
