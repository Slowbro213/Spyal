package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
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
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	var u models.User
	err := sqlx.GetContext(ctx, r.db, &u, `SELECT id, username, password FROM users WHERE id=$1`, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	return &u, nil
}


func (r *userRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var u models.User
	err := sqlx.GetContext(ctx, r.db, &u,
		`SELECT id, username, password FROM users WHERE username=$1`, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get by username: %w", err)
	}
	return &u, nil
}

