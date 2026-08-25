// Package vault содержит бизнес-логику работы с элементами хранилища.
package vault

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"credentials-vault/server/internal/domain"
)

//go:generate mockgen -source=vault.go -destination=mocks/vault_repository_mock.gen.go -package=mocks VaultRepository
type VaultRepository interface {
	Create(ctx context.Context, item *domain.VaultItem) error
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.VaultItem, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, itemType *domain.ItemType) ([]*domain.VaultItem, error)
	Update(ctx context.Context, item *domain.VaultItem) error
	SoftDelete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type VaultService struct {
	repo VaultRepository
}

func NewService(repo VaultRepository) *VaultService {
	return &VaultService{
		repo: repo,
	}
}

// CreateItem создаёт новый элемент хранилища
func (s *VaultService) CreateItem(ctx context.Context, userID uuid.UUID, itemType domain.ItemType, encryptedData []byte, metadata map[string]string) (*domain.VaultItem, error) {
	if err := validateItemType(itemType); err != nil {
		return nil, err
	}

	if len(encryptedData) == 0 {
		return nil, domain.ErrEncryptedDataRequired
	}

	if userID == uuid.Nil {
		return nil, domain.ErrUserIDRequired
	}

	if metadata == nil {
		metadata = make(map[string]string)
	}

	item := &domain.VaultItem{
		UserID:        userID,
		Type:          itemType,
		EncryptedData: encryptedData,
		Metadata:      metadata,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to create vault item: %w", err)
	}

	return item, nil
}

// GetItem возвращает элемент по ID.
func (s *VaultService) GetItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.VaultItem, error) {
	if id == uuid.Nil {
		return nil, domain.ErrVaultItemIDRequired
	}

	if userID == uuid.Nil {
		return nil, domain.ErrUserIDRequired
	}

	item, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, domain.ErrVaultItemNotFound) {
			return nil, domain.ErrVaultItemNotFound
		}
		return nil, fmt.Errorf("failed to get vault item: %w", err)
	}

	return item, nil
}

// UpdateItem обновляет элемент хранилища.
func (s *VaultService) UpdateItem(ctx context.Context, id uuid.UUID, userID uuid.UUID, encryptedData []byte, metadata map[string]string) (*domain.VaultItem, error) {
	if id == uuid.Nil {
		return nil, domain.ErrVaultItemIDRequired
	}

	if userID == uuid.Nil {
		return nil, domain.ErrUserIDRequired
	}

	if len(encryptedData) == 0 {
		return nil, domain.ErrEncryptedDataRequired
	}

	if metadata == nil {
		metadata = make(map[string]string)
	}

	item := &domain.VaultItem{
		ID:            id,
		UserID:        userID,
		EncryptedData: encryptedData,
		Metadata:      metadata,
	}

	if err := s.repo.Update(ctx, item); err != nil {
		if errors.Is(err, domain.ErrVaultItemNotFound) {
			return nil, domain.ErrVaultItemNotFound
		}
		return nil, fmt.Errorf("failed to update vault item: %w", err)
	}

	return item, nil
}

// DeleteItem мягко удаляет элемент хранилища.
func (s *VaultService) DeleteItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	if id == uuid.Nil {
		return domain.ErrVaultItemIDRequired
	}

	if userID == uuid.Nil {
		return domain.ErrUserIDRequired
	}

	if err := s.repo.SoftDelete(ctx, id, userID); err != nil {
		if errors.Is(err, domain.ErrVaultItemNotFound) {
			return domain.ErrVaultItemNotFound
		}
		return fmt.Errorf("failed to delete vault item: %w", err)
	}

	return nil
}

// ListItems возвращает список элементов пользователя.
func (s *VaultService) ListItems(ctx context.Context, userID uuid.UUID, itemType *domain.ItemType) ([]*domain.VaultItem, error) {
	if userID == uuid.Nil {
		return nil, domain.ErrUserIDRequired
	}

	if itemType != nil {
		if err := validateItemType(*itemType); err != nil {
			return nil, err
		}
	}

	items, err := s.repo.FindByUserID(ctx, userID, itemType)
	if err != nil {
		return nil, fmt.Errorf("failed to list vault items: %w", err)
	}

	return items, nil
}
