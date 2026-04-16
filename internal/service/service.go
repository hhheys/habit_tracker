package service

import (
	"habit-tracker/internal/auth"
	"habit-tracker/internal/repository/postgres"

	"go.uber.org/zap"
)

// Service provides access to the service.
//
//go:generate mockery --name=Service --output=../../../mocks--outpkg=mocks
type Service interface {
	UserService
	HabitService
	StreakService
}

type serviceImpl struct {
	log        *zap.Logger
	repository postgres.Repository

	UserService
	HabitService
	StreakService
}

// NewService создаёт новый базовый сервис
func NewService(log *zap.Logger, repository postgres.Repository, jwt auth.JWTService) Service {
	return &serviceImpl{
		log:        log,
		repository: repository,

		UserService:   NewUserService(log, jwt, repository),
		HabitService:  NewHabitService(log, repository),
		StreakService: NewStreakService(log, repository),
	}
}
