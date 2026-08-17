package vault

import (
	"credentials-vault/server/internal/domain"
)

func validateItemType(itemType domain.ItemType) error {
	switch itemType {
	case domain.ItemTypeLogin,
		domain.ItemTypeCard,
		domain.ItemTypeText,
		domain.ItemTypeBinary:
		return nil
	default:
		return domain.ErrInvalidItemType
	}
}
