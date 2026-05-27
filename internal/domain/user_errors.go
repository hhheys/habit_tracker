package domain

import "errors"

var (
	//ErrEmailAlreadyExists = errors.New("user already exists")
	//ErrWrongPassword      = errors.New("wrong password")

	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidRole       = errors.New("invalid role")
	ErrNoPermissions     = errors.New("no permissions")
	ErrUnauthorized      = errors.New("unauthorized")
)
