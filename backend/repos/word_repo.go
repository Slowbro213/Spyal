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
	repoInterface
	Create(ctx context.Context, w *models.Word) error
	AddRelated(ctx context.Context, wordID1, wordID2 int64) error
	RandomPair(ctx context.Context) (main, related *models.Word, err error)
}

type wordRepo struct {
	repo
}

func NewWordRepo(db sqlx.ExtContext) WordRepository {
	return &wordRepo{repo: repo{db: db}}
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

type pair struct {
	MainID      int64  `db:"main_id"`
	MainWord    string `db:"main_word"`
	RelatedID   int64  `db:"related_id"`
	RelatedWord string `db:"related_word"`
}

func (r *wordRepo) RandomPair(ctx context.Context) (main, related *models.Word, err error) {
	var p pair
	err = sqlx.GetContext(ctx, r.db, &p, `
		SELECT w1.id            AS main_id,
		       w1.word          AS main_word,
		       COALESCE(w2.id, w1.id)   AS related_id,
		       COALESCE(w2.word, w1.word) AS related_word
		FROM  (SELECT id, word FROM words ORDER BY RANDOM() LIMIT 1) w1
		LEFT  JOIN LATERAL (
		        SELECT w.id, w.word
		        FROM   word_related wr
		        JOIN   words w ON w.id = CASE WHEN wr.word_id_1 = w1.id
		                                      THEN wr.word_id_2
		                                      ELSE wr.word_id_1 END
		        WHERE  w1.id IN (wr.word_id_1, wr.word_id_2)
		        ORDER  BY RANDOM()
		        LIMIT   1
		      ) w2 ON true`)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, ErrWordNotFound
		}
		return nil, nil, fmt.Errorf("random pair: %w", err)
	}

	word := &models.Word{ID: p.MainID, Word: p.MainWord}
	related = &models.Word{ID: p.RelatedID, Word: p.RelatedWord}
	return word, related, nil
}
