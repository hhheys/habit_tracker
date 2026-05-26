package domain

import (
	"time"

	"github.com/google/uuid"
)

type RefreshSession struct {
	ID        uuid.UUID
	UserID    uint
	TokenHash string
	ExpiresAt time.Time
	Revoked   bool
	UserAgent string
	IPAddress string
	CreatedAt time.Time
}
