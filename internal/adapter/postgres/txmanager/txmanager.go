package txmanager

import (
	"context"
	"database/sql"
)

type TXKey struct{}

type TXManager struct {
	db *sql.DB
}

func NewTXManager(db *sql.DB) *TXManager {
	return &TXManager{db: db}
}

func (m *TXManager) WithTx(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, TXKey{}, tx)

	err = fn(ctx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
