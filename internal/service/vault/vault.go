// Package vault содержит бизнес-логику работы с элементами хранилища.
package vault

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"credentials-vault/internal/domain"
)

//go:generate mockgen -source=vault.go -destination=mocks/vault_repository_mock.go -package=mocks VaultRepository
type VaultRepository interface {
	Create(ctx context.Context, item *domain.VaultItem) error
	FindByUserID(ctx context.Context, userID uuid.UUID, itemType *domain.ItemType) ([]*domain.VaultItem, error)
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
