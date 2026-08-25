package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"credentials-vault/server/internal/domain"
)

type VaultItemRepository struct {
	db *pgxpool.Pool
}

func NewVaultItemRepository(db *pgxpool.Pool) *VaultItemRepository {
	return &VaultItemRepository{db: db}
}

// Create создаёт новый элемент хранилища.
func (r *VaultItemRepository) Create(ctx context.Context, item *domain.VaultItem) error {
	metadataJSON, err := json.Marshal(item.Metadata)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	query := `
		INSERT INTO vault_items (user_id, type, encrypted_data, metadata)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err = r.db.QueryRow(
		ctx,
		query,
		item.UserID,
		item.Type,
		item.EncryptedData,
		metadataJSON,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		return fmt.Errorf("create vault item: %w", err)
	}

	return nil
}

// FindByID находит элемент по ID и user_id.
func (r *VaultItemRepository) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.VaultItem, error) {
	item := &domain.VaultItem{}
	var metadataJSON []byte

	query := `
		SELECT id, user_id, type, encrypted_data, metadata, created_at, updated_at, deleted_at
		FROM vault_items
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`

	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&item.ID,
		&item.UserID,
		&item.Type,
		&item.EncryptedData,
		&metadataJSON,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrVaultItemNotFound
		}
		return nil, fmt.Errorf("find vault item by id: %w", err)
	}

	if err := json.Unmarshal(metadataJSON, &item.Metadata); err != nil {
		return nil, fmt.Errorf("unmarshal metadata: %w", err)
	}

	return item, nil
}

// FindByUserID находит все элементы пользователя.
func (r *VaultItemRepository) FindByUserID(ctx context.Context, userID uuid.UUID, itemType *domain.ItemType) ([]*domain.VaultItem, error) {
	query := `
		SELECT id, user_id, type, encrypted_data, metadata, created_at, updated_at, deleted_at
		FROM vault_items
		WHERE user_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{userID}

	if itemType != nil {
		query += " AND type = $2"
		args = append(args, *itemType)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("find vault items by user id: %w", err)
	}
	defer rows.Close()

	items := make([]*domain.VaultItem, 0)
	for rows.Next() {
		item := &domain.VaultItem{}
		var metadataJSON []byte

		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Type,
			&item.EncryptedData,
			&metadataJSON,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scan vault item: %w", err)
		}

		if err := json.Unmarshal(metadataJSON, &item.Metadata); err != nil {
			return nil, fmt.Errorf("unmarshal metadata: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate vault items: %w", err)
	}

	return items, nil
}

// Update обновляет encrypted_data и metadata элемента.
func (r *VaultItemRepository) Update(ctx context.Context, item *domain.VaultItem) error {
	metadataJSON, err := json.Marshal(item.Metadata)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	query := `
		UPDATE vault_items
		SET encrypted_data = $1, metadata = $2, updated_at = NOW()
		WHERE id = $3 AND user_id = $4 AND deleted_at IS NULL
		RETURNING updated_at
	`

	err = r.db.QueryRow(
		ctx,
		query,
		item.EncryptedData,
		metadataJSON,
		item.ID,
		item.UserID,
	).Scan(&item.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrVaultItemNotFound
		}
		return fmt.Errorf("update vault item: %w", err)
	}

	return nil
}

// SoftDelete мягко удаляет элемент.
func (r *VaultItemRepository) SoftDelete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `
		UPDATE vault_items
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("soft delete vault item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrVaultItemNotFound
	}

	return nil
}
