package domain

import "errors"

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionRevoked  = errors.New("session revoked")
	ErrTokenExpired    = errors.New("token expired")
)
