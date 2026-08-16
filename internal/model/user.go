// Package model содержит доменные модели.
package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
