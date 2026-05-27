package user

import (
	"context"
	"database/sql"
	"errors"
	"habit-tracker/internal/domain"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

type Repository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewRepository(db *sql.DB, log *zap.Logger) *Repository {
	return &Repository{db: db, log: log}
}

func (r *Repository) Create(ctx context.Context, user *domain.User) error {
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO users (username, email, password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, is_active`,
		user.Username, user.Email, user.Password, user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.IsActive)
	if isUniqueViolation(err) {
		return domain.ErrUserAlreadyExists
	}
	return err
}

func (r *Repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.exists(ctx, "email", email)
}

func (r *Repository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return r.exists(ctx, "username", username)
}

func (r *Repository) exists(ctx context.Context, column, value string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE `+column+` = $1)`, value).Scan(&exists)
	return exists, err
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.get(ctx, `SELECT id, username, email, password, role, created_at, is_active FROM users WHERE email = $1`, email)
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return r.get(ctx, `SELECT id, username, email, password, role, created_at, is_active FROM users WHERE username = $1`, username)
}

func (r *Repository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	return r.get(ctx, `SELECT id, username, email, password, role, created_at, is_active FROM users WHERE id = $1`, id)
}

func (r *Repository) get(ctx context.Context, query string, arg any) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.IsActive,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		r.log.Error("postgres user query failed", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}
