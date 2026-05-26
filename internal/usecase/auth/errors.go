package auth

import "errors"

var (
	// ErrInvalidToken       = errors.New("invalid token")

	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token expired")
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionRevoked     = errors.New("session revoked")
)
