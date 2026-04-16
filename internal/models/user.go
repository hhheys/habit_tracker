package models

import (
	"habit-tracker/internal/dto/request"
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
	Role      string
	CreatedAt time.Time
	IsActive  bool
}

func NewUser(req *request.UserRegisterRequest) *User {
	return &User{
		ID:        0,
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		Role:      "user",
		CreatedAt: time.Time{},
		IsActive:  false,
	}
}
