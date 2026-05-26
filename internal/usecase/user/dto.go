package user

import (
	"habit-tracker/internal/domain"
	"time"
)

type Output struct {
	ID        uint
	Username  string
	Email     string
	Role      domain.UserRole
	CreatedAt time.Time
}
