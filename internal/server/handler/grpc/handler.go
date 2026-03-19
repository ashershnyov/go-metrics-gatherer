package grpc

import (
	"context"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/db"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
	"github.com/ashershnyov/go-metrics-gatherer/pkg/api/proto"
)

type metricService interface {
	UpdateMultipleMetrics(ctx context.Context, metrics []model.InternalMetric) error
}

// MetricsHandler is a handler for updating and getting metrics.
type MetricsHandler struct {
	proto.UnimplementedMetricsServer
	service metricService
	db      db.DB
}

// NewMetricsHandler returns an empty metrics update handler.
func NewMetricsHandler(s metricService, db *db.Postgres) *MetricsHandler {
	return &MetricsHandler{
		service: s,
		db:      db,
	}
}
