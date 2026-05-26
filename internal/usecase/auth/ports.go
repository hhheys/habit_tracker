package auth

import (
	"context"
	"habit-tracker/internal/domain"
)

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePasswordHash(hashedPassword, password string) bool
}

type TokenHasher interface {
	HashToken(token string) (string, error)
	CompareToken(hashedToken, token string) bool
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, id uint) (*domain.User, error)
}

type RefreshSessionRepository interface {
	Create(ctx context.Context, session *domain.RefreshSession) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshSession, error)
	Rotate(ctx context.Context, tokenHash string, replacement *domain.RefreshSession) error
}

type AccessTokenGenerator interface {
	GenerateToken(user *TokenSubject) (string, error)
}

type RefreshTokenGenerator interface {
	GenerateToken(user *TokenSubject) (string, error)
}
