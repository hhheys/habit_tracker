package errors

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrWrongPassword     = errors.New("wrong password")
	ErrInvalidRole       = errors.New("invalid role")
	ErrNoPermissions     = errors.New("no permissions")
	ErrUnauthorized      = errors.New("unauthorized")
)
