package session

import (
	"context"
	"database/sql"
	"errors"
	"habit-tracker/internal/domain"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Repository struct {
	db *sql.DB

	log *zap.Logger
}

func NewRepository(db *sql.DB, logger *zap.Logger) *Repository {
	return &Repository{db: db, log: logger}
}

func (r *Repository) Create(ctx context.Context, session *domain.RefreshSession) error {
	if session.ID == uuid.Nil {
		session.ID = uuid.New()
	}
	return r.db.QueryRowContext(ctx, `
		INSERT INTO refresh_session
			(id, user_id, token_hash, expires_at, revoked, user_agent, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at`,
		session.ID, session.UserID, session.TokenHash, session.ExpiresAt,
		session.Revoked, session.UserAgent, session.IPAddress,
	).Scan(&session.CreatedAt)
}

func (r *Repository) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshSession, error) {
	var session domain.RefreshSession
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, token_hash, expires_at, revoked, user_agent, ip_address, created_at
		FROM refresh_session WHERE token_hash = $1`, tokenHash,
	).Scan(
		&session.ID, &session.UserID, &session.TokenHash, &session.ExpiresAt,
		&session.Revoked, &session.UserAgent, &session.IPAddress, &session.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSessionNotFound
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *Repository) Rotate(ctx context.Context, tokenHash string, replacement *domain.RefreshSession) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			r.log.Error("failed to rollback transaction", zap.Error(err))
		}
	}(tx)

	var current domain.RefreshSession
	err = tx.QueryRowContext(ctx, `
		SELECT id, user_id, expires_at, revoked
		FROM refresh_session
		WHERE token_hash = $1
		FOR UPDATE`, tokenHash,
	).Scan(&current.ID, &current.UserID, &current.ExpiresAt, &current.Revoked)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrSessionNotFound
	}
	if err != nil {
		return err
	}
	if current.Revoked {
		return domain.ErrSessionRevoked
	}
	if time.Now().After(current.ExpiresAt) {
		return domain.ErrTokenExpired
	}
	if replacement.UserID != current.UserID {
		return domain.ErrSessionNotFound
	}
	if replacement.ID == uuid.Nil {
		replacement.ID = uuid.New()
	}
	_, err = tx.ExecContext(ctx, `UPDATE refresh_session SET revoked = TRUE WHERE id = $1`, current.ID)
	if err != nil {
		return err
	}
	err = tx.QueryRowContext(ctx, `
		INSERT INTO refresh_session
			(id, user_id, token_hash, expires_at, revoked, user_agent, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at`,
		replacement.ID, replacement.UserID, replacement.TokenHash, replacement.ExpiresAt,
		replacement.Revoked, replacement.UserAgent, replacement.IPAddress,
	).Scan(&replacement.CreatedAt)
	if err != nil {
		return err
	}
	return tx.Commit()
}
