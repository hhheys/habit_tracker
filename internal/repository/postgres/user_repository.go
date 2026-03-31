package postgres

import (
	"database/sql"
	"errors"
	"habit-tracker/internal/models"

	appErrors "habit-tracker/internal/errors"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

// UserRepository defines the interface for user repository
//
//go:generate mockery --name=UserRepository --output=../../../mocks --outpkg=mocks
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
}

// UserRepositoryImpl implements the UserRepository interface
type userRepositoryImpl struct {
	DB     *sql.DB
	logger *zap.Logger
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB, logger *zap.Logger) UserRepository {
	return &userRepositoryImpl{
		DB:     db,
		logger: logger,
	}
}

func (r *userRepositoryImpl) CreateUser(user *models.User) error {
	err := r.DB.QueryRow(
		`INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id, created_at, is_active`,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.IsActive)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // unique_violation
				r.logger.Warn("User already exists", zap.String("email", user.Email))
				return appErrors.ErrUserAlreadyExists
			}
		}
		r.logger.Error("Failed to insert user", zap.Error(err))
		return err
	}

	return nil
}

// GetUserByEmail retrieves a user by email
func (r *userRepositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow(
		`SELECT id, username, email, password, role, created_at FROM users WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErrors.ErrUserNotFound
		}
		r.logger.Error("Failed to get user by email", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow(
		`SELECT id, username, email, password, role, created_at FROM users WHERE id = $1`,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErrors.ErrUserNotFound
		}
		r.logger.Error("Failed to get user by id", zap.Error(err))
		return nil, err
	}
	return &user, nil
}
