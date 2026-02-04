package db

import (
	"context"
	"database/sql"
)

// DB defines a database interface.
//
//go:generate mockgen -source=./db.go -destination=./../../../mocks/pg_db_mock.go -package=mocks . DB
type DB interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	PingContext(context.Context) error
}
