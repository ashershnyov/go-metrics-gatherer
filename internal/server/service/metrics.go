package service

import (
	"context"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
)

type MetricStorage interface {
	GetGauges(ctx context.Context) (map[string]float64, error)
	GetCounters(ctx context.Context) (map[string]int64, error)
	UpdateGauge(ctx context.Context, name string, val float64) error
	UpdateCounter(ctx context.Context, name string, val int64) error
	GetGauge(ctx context.Context, name string) (float64, error)
	GetCounter(ctx context.Context, name string) (int64, error)
	UpdateMultipleMetrics(ctx context.Context, metrics []model.InternalMetric) error
}

// Service defines the sevice layer.
type Service struct {
	storage MetricStorage
}

// NewService returns a new service with an underlying storage s.
func NewService(s MetricStorage) *Service {
	return &Service{
		storage: s,
	}
}

// ListMetrics returns all registered metrics.
func (s *Service) ListMetrics(ctx context.Context) ([]model.InternalMetric, error) {
	gauges, err := s.storage.GetGauges(ctx)
	if err != nil {
		return nil, err
	}
	counters, err := s.storage.GetCounters(ctx)
	if err != nil {
		return nil, err
	}
	metrics := make([]model.InternalMetric, 0, len(gauges)+len(counters))

	for name, value := range gauges {
		m := model.InternalMetric{
			Name:  name,
			Value: value,
			Type:  model.Gauge,
		}
		metrics = append(metrics, m)
	}

	for name, value := range counters {
		m := model.InternalMetric{
			Name:  name,
			Delta: value,
			Type:  model.Counter,
		}
		metrics = append(metrics, m)
	}

	return metrics, nil
}

// UpdateMetric updates vales in s using data provided in m.
func (s *Service) UpdateMetric(ctx context.Context, m model.InternalMetric) error {
	var err error
	switch m.Type {
	case model.Counter:
		err = s.storage.UpdateCounter(ctx, m.Name, m.Delta)
	case model.Gauge:
		err = s.storage.UpdateGauge(ctx, m.Name, m.Value)
	}
	return err
}

// UpdateMultipleMetrics updates values in all passed metrics.
func (s *Service) UpdateMultipleMetrics(ctx context.Context, metrics []model.InternalMetric) error {
	return s.storage.UpdateMultipleMetrics(ctx, metrics)
}

// GetMetric returns metric of specified type and name.
func (s *Service) GetMetric(ctx context.Context, name string, typ model.MetricType) (model.InternalMetric, error) {
	var err error
	m := model.InternalMetric{
		Name: name,
		Type: typ,
	}
	switch typ {
	case model.Counter:
		m.Delta, err = s.storage.GetCounter(ctx, m.Name)
	case model.Gauge:
		m.Value, err = s.storage.GetGauge(ctx, m.Name)
	}
	return m, err
}
