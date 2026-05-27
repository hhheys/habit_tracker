package auth

import (
	"errors"

	"habit-tracker/internal/domain"
)

var (
	// ErrInvalidToken       = errors.New("invalid token")

	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = domain.ErrTokenExpired
	ErrSessionNotFound    = domain.ErrSessionNotFound
	ErrSessionRevoked     = domain.ErrSessionRevoked
)
