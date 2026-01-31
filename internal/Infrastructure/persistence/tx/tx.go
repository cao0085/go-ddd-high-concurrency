package tx

import (
	"context"
	"database/sql"
)

type ctxKey string

const txKey ctxKey = "tx"

// executor 是 sql.DB 和 sql.Tx 的共用介面
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func GetConn(ctx context.Context, db *sql.DB) Executor {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}
	return db
}

func WithTx(ctx context.Context, db *sql.DB, fn func(ctx context.Context) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txKey, tx)

	if err := fn(txCtx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
