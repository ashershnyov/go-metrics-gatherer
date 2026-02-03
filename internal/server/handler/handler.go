package handler

import (
	"context"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/db"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
)

type metricService interface {
	ListMetrics(ctx context.Context) ([]model.InternalMetric, error)
	UpdateMetric(ctx context.Context, metric model.InternalMetric) error
	UpdateMultipleMetrics(ctx context.Context, metrics []model.InternalMetric) error
	GetMetric(ctx context.Context, name string, typ model.MetricType) (model.InternalMetric, error)
}

//go:generate mockgen -source=./handler.go -destination=./../../../mocks/audit_logger_mock.go -package=mocks . auditLogger
type auditLogger interface {
	Log(context.Context, []string, string) error
}

// MetricsHandler is a handler for updating and getting metrics.
type MetricsHandler struct {
	service metricService
	db      *db.Postgres
	audit   auditLogger
}

// NewMetricsHandler returns an empty metrics update handler.
func NewMetricsHandler(s metricService, db *db.Postgres, auditLogger auditLogger) *MetricsHandler {
	return &MetricsHandler{
		service: s,
		db:      db,
		audit:   auditLogger,
	}
}
