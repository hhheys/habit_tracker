package domain

import (
	"time"
)

type UserRole string

const (
	UserRoleDefault UserRole = "user"
	UserRoleAdmin   UserRole = "admin"
)

type User struct {
	ID        uint
	Username  string
	Email     string
	Password  string
	Role      UserRole
	CreatedAt time.Time
	IsActive  bool
}
