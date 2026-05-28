package txmanager

import (
	"context"
	"database/sql"
)

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Exec(query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryRow(query string, args ...any) *sql.Row
}

func ExecutorFromContext(ctx context.Context, db *sql.DB) Executor {
	tx, ok := ctx.Value(TXKey{}).(*sql.Tx)
	if ok {
		return tx
	}

	return db
}
