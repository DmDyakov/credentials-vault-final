// Package model содержит модели данных клиента.
package model

import "time"

// VaultItem — элемент хранилища с расшифрованными данными.
type VaultItem struct {
	ID        string
	Type      string
	Secret    map[string]string
	Metadata  map[string]string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ListVaultItem — элемент списка.
type ListVaultItem struct {
	ID        string
	Type      string
	Metadata  map[string]string
	CreatedAt time.Time
	UpdatedAt time.Time
}
