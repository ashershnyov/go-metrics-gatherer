package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ashershnyov/go-metrics-gatherer/internal/retrier"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func markUnretriable(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return errors.Join(retrier.ErrUnretriable, err)
	}

	switch pgErr.Code {
	case pgerrcode.ConnectionException, pgerrcode.ConnectionDoesNotExist, pgerrcode.ConnectionFailure,
		pgerrcode.SQLClientUnableToEstablishSQLConnection, pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection, pgerrcode.TransactionResolutionUnknown,
		pgerrcode.ProtocolViolation:

		return err
	}

	return errors.Join(retrier.ErrUnretriable, err)
}

// Postgres is a simple wrapper over the sql.DB connector.
type Postgres struct {
	*sql.DB
	maxRetries int
}

// NewPostgres creates a Postgres wrapper.
func NewPostgres(ctx context.Context, address string, maxRetries int) (*Postgres, error) {
	db, err := sql.Open("pgx", address)
	if err != nil {
		return nil, err
	}
	return &Postgres{
		DB:         db,
		maxRetries: maxRetries,
	}, nil
}

func (p *Postgres) SQLDB() *sql.DB {
	return p.DB
}

func (p *Postgres) PingContext(ctx context.Context) error {
	return retrier.WithRetry(p.maxRetries, func() error {
		return markUnretriable(p.DB.PingContext(ctx))
	})
}

func (p *Postgres) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	var (
		r   *sql.Rows
		err error
	)
	err = retrier.WithRetry(p.maxRetries, func() error {
		var err error
		r, err = p.DB.QueryContext(ctx, query, args...)
		if r != nil && r.Err() != nil {
			return markUnretriable(r.Err())
		}
		return markUnretriable(err)
	})
	return r, err
}

func (p *Postgres) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	var (
		r   sql.Result
		err error
	)
	err = retrier.WithRetry(p.maxRetries, func() error {
		var err error
		r, err = p.DB.ExecContext(ctx, query, args...)
		return markUnretriable(err)
	})
	return r, err
}

func (p *Postgres) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	var (
		t   *sql.Tx
		err error
	)
	err = retrier.WithRetry(p.maxRetries, func() error {
		var err error
		t, err = p.DB.BeginTx(ctx, opts)
		return markUnretriable(err)
	})
	return t, err
}
