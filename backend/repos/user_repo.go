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

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
}

type userRepo struct {
	db sqlx.ExtContext
}

func NewUserRepo(db sqlx.ExtContext) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *models.User) error {
	query := `INSERT INTO users (username, password) VALUES (:username, :password) RETURNING id`
	rows, err := sqlx.NamedQueryContext(ctx, r.db, query, u)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&u.ID); err != nil {
			return fmt.Errorf("scan id: %w", err)
		}
	}

	data, _ := json.Marshal(u)
	_ = cache.Set(ctx, u.CacheKey(), string(data), time.Hour)

	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	key := fmt.Sprintf("user_id_%d", id)

	if cached, err := cache.Get(ctx, key); err == nil && cached != "" {
		var u models.User
		if jsonErr := json.Unmarshal([]byte(cached), &u); jsonErr == nil {
			return &u, nil
		}
	}

	var u models.User
	err := sqlx.GetContext(ctx, r.db, &u, `SELECT id, username, password FROM users WHERE id=$1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get by id: %w", err)
	}

	data, _ := json.Marshal(u)
	_ = cache.Set(ctx, key, string(data), time.Hour)

	return &u, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	key := "user_" + username

	if cached, err := cache.Get(ctx, key); err == nil && cached != "" {
		var u models.User
		if jsonErr := json.Unmarshal([]byte(cached), &u); jsonErr == nil {
			return &u, nil
		}
	}

	var u models.User
	err := sqlx.GetContext(ctx, r.db, &u,
		`SELECT id, username, password FROM users WHERE username=$1`, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get by username: %w", err)
	}

	data, _ := json.Marshal(u)
	_ = cache.Set(ctx, key, string(data), time.Hour)

	return &u, nil
}
