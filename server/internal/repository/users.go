// Package repository реализует доступ к данным.
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"credentials-vault/server/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
        INSERT INTO users (username, password_hash, salt)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRow(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Salt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrUserAlreadyExists
		}
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	user := &domain.User{}
	query := `
        SELECT id, username, password_hash, salt, created_at, updated_at
        FROM users
        WHERE username = $1
    `

	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Salt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by username: %w", err)
	}

	return user, nil
}
