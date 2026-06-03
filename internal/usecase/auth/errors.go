package auth

import (
	"errors"
)

var (
	// ErrInvalidToken       = errors.New("invalid token")

	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidTimezone    = errors.New("invalid timezone")
	//ErrTokenExpired       = domain.ErrTokenExpired
	//ErrSessionNotFound    = domain.ErrSessionNotFound
	//ErrSessionRevoked     = domain.ErrSessionRevoked
)
