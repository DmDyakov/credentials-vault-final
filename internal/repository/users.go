// Package repository реализует доступ к данным.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"credentials-vault/internal/domain"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
        INSERT INTO users (username, password_hash)
        VALUES ($1, $2)
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	user := &domain.User{}
	query := `
        SELECT id, username, password_hash, created_at, updated_at
        FROM users
        WHERE username = $1
    `

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by username: %w", err)
	}

	return user, nil
}
