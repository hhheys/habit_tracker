package postgres

import (
	"database/sql"

	"go.uber.org/zap"
)

// Repository provides access to the repository.
//
//go:generate mockery --name=Repository --output=../../../mocks --outpkg=mocks
type Repository interface {
	UserRepository
	HabitRepository
}

// RepositoryImpl implements Repository.
type RepositoryImpl struct {
	db  *sql.DB
	log *zap.Logger

	UserRepository
	HabitRepository
}

// NewRepository creates a new Repository.
func NewRepository(db *sql.DB, logger *zap.Logger) Repository {
	return &RepositoryImpl{
		db:              db,
		log:             logger,
		UserRepository:  NewUserRepository(db, logger),
		HabitRepository: NewHabitRepository(db, logger),
	}
}
