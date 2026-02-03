package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/db"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
)

//go:generate mockgen -source=./db.go -destination=./../../../mocks/pg_db_mock.go -package=mocks . pgDB
type pgDB interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
}

// DB is a storage adapter for a database.
type DB struct {
	db pgDB
}

// NewDB creates a new DB storage.
func NewDB(db *db.Postgres) *DB {
	return &DB{
		db: db,
	}
}

const qGetGauges = "SELECT name, value FROM metrics WHERE type='gauge';"

// GetGauges returns all gauges stored upon calling.
func (d *DB) GetGauges(ctx context.Context) (Gauges, error) {
	rows, err := d.db.QueryContext(ctx, qGetGauges)
	if err != nil {
		return nil, fmt.Errorf("an error occurred when fetching gauges from DB: %w", err)
	}
	defer rows.Close()

	counters := make(Gauges)
	for rows.Next() {
		var (
			val  float64
			name string
		)
		err = rows.Scan(&name, &val)
		if err != nil {
			return nil, fmt.Errorf("an error occurred when fetching gauges from DB: %w", err)
		}
		counters[name] = val
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("an error occurred when fetching gauges from DB: %w", err)
	}
	return counters, nil
}

const qGetCounters = "SELECT name, delta FROM metrics WHERE type='counter';"

// GetCounters returns all counters stored upon calling.
func (d *DB) GetCounters(ctx context.Context) (Counters, error) {
	rows, err := d.db.QueryContext(ctx, qGetCounters)
	if err != nil {
		return nil, fmt.Errorf("an error occurred when fetching counters from DB: %w", err)
	}
	defer rows.Close()

	counters := make(Counters)
	for rows.Next() {
		var (
			val  int64
			name string
		)
		err = rows.Scan(&name, &val)
		if err != nil {
			return nil, fmt.Errorf("an error occurred when fetching counters from DB: %w", err)
		}
		counters[name] = val
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("an error occurred when fetching counters from DB: %w", err)
	}
	return counters, nil
}

const qGetGauge = "SELECT value FROM metrics WHERE id=$1 AND type='gauge';"

// GetGauge returns the value of Gauge by the specified name and the indication whether the metric exists.
// If metric does not exist yet, will return (0.0, false).
func (d *DB) GetGauge(ctx context.Context, name string) (float64, error) {
	row := d.db.QueryRowContext(ctx, qGetGauge, name)
	var val float64
	err := row.Scan(&val)
	if err != nil {
		return 0., fmt.Errorf("an error occurred when fetching gauge %s from DB: %w", name, err)
	}
	return val, nil
}

const qUpdateGauge = "INSERT INTO metrics (id, type, value) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET value = $3;"

// UpdateGauge replaces the vaule of the Gauge by the specified name.
// Creates a new Gauge with the specified name if it does not exist yet.
func (d *DB) UpdateGauge(ctx context.Context, name string, val float64) error {
	_, err := d.db.ExecContext(ctx, qUpdateGauge, name, "gauge", val)
	if err != nil {
		return fmt.Errorf("an error occurred when writing gauge %s", name)
	}
	return nil
}

const qGetCounter = "SELECT delta FROM metrics WHERE id=$1 AND type='counter';"

// GetCounter returns the value of Counter by the specified name and the indication whether the metric exists.
// If metric does not exist yet, will return (0.0, false).
func (d *DB) GetCounter(ctx context.Context, name string) (int64, error) {
	row := d.db.QueryRowContext(ctx, qGetCounter, name)
	var val int64
	err := row.Scan(&val)
	if err != nil {
		return 0., fmt.Errorf("an error occurred when fetching counter %s from DB: %w", name, err)
	}
	return val, nil
}

const qUpdateCounter = "INSERT INTO metrics (id, type, delta) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET delta = metrics.delta + $3;"

// UpdateCounter adds val to the Counter with the specified name.
// Creates a new Counter with the specified name if it does not exist yet.
func (d *DB) UpdateCounter(ctx context.Context, name string, val int64) error {
	_, err := d.db.ExecContext(ctx, qUpdateCounter, name, "counter", val)
	if err != nil {
		return fmt.Errorf("an error occurred when writing counter %s", name)
	}
	return nil
}

// UpdateMultipleMetrics updates multiple metrics' values at once.
func (d *DB) UpdateMultipleMetrics(ctx context.Context, metrics []model.InternalMetric) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("an error occurred when writing multiple metrics %s", err)
	}
	for _, metric := range metrics {
		switch metric.Type {
		case model.Gauge:
			_, err = tx.ExecContext(ctx, qUpdateGauge, metric.Name, metric.Type, metric.Value)
		case model.Counter:
			_, err = tx.ExecContext(ctx, qUpdateCounter, metric.Name, metric.Type, metric.Delta)
		}
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("an error occurred when writing multiple metrics %s", err)
		}
	}
	return tx.Commit()
}
