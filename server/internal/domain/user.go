package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Username  string
	Password  string
	Salt      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}
