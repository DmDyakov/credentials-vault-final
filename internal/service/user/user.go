// Package user содержит бизнес-логику работы с пользователями.
package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"credentials-vault/internal/domain"
)

//go:generate mockgen -source=auth.go -destination=mocks/users_mock.go -package=mocks UserRepository
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Register(ctx context.Context, username, password string) (*domain.User, error) {
	if err := validateRegisterCredentials(username, password); err != nil {
		return nil, err
	}

	_, err := s.repo.FindByUsername(ctx, username)
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
	case err == nil:
		return nil, domain.ErrUserAlreadyExists
	default:
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Username: username,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (*domain.User, error) {
	if err := validateLoginCredentials(username, password); err != nil {
		return nil, err
	}

	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}
