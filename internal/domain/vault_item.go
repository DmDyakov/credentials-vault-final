package domain

import (
	"time"

	"github.com/google/uuid"
)

// ItemType - тип элемента хранилища
type ItemType string

const (
	ItemTypeLogin  ItemType = "LOGIN"
	ItemTypeCard   ItemType = "CARD"
	ItemTypeText   ItemType = "TEXT"
	ItemTypeBinary ItemType = "BINARY"
)

// VaultItem - элемент хранилища
type VaultItem struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Type          ItemType
	EncryptedData []byte
	Metadata      map[string]string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}
