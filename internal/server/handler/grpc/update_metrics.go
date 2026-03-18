package grpc

import (
	"context"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
	"github.com/ashershnyov/go-metrics-gatherer/pkg/api/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateMetrics updates multiple metrics.
func (s *MetricsHandler) UpdateMetrics(ctx context.Context, req *proto.UpdateMetricsRequest) (*proto.UpdateMetricsResponse, error) {
	metrics := []model.InternalMetric{}
	for _, m := range req.Metrics {
		if m == nil {
			continue
		}

		switch m.Type {
		case proto.Metric_GAUGE:
			metrics = append(metrics, model.InternalMetric{
				Name:  m.Id,
				Type:  "gauge",
				Value: m.Value,
			})

		case proto.Metric_COUNTER:
			metrics = append(metrics, model.InternalMetric{
				Name:  m.Id,
				Type:  "counter",
				Delta: m.Delta,
			})

		default:
			return nil, status.Errorf(codes.InvalidArgument, "got unknown metric type: %v", m.Type)
		}
	}

	if err := s.service.UpdateMultipleMetrics(ctx, metrics); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update metrics: %v", err)
	}

	return &proto.UpdateMetricsResponse{}, nil
}
