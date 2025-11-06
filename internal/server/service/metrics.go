package service

import (
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
)

type MetricStorage interface {
	GetGauges() map[string]float64
	GetCounters() map[string]int64
	UpdateGauge(name string, val float64)
	UpdateCounter(name string, val int64)
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
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
func (s *Service) ListMetrics() []model.InternalMetric {
	gauges := s.storage.GetGauges()
	counters := s.storage.GetCounters()
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

	return metrics
}

// UpdateMetric updates vales in s using data provided in m.
func (s *Service) UpdateMetric(m model.InternalMetric) {
	switch m.Type {
	case model.Counter:
		s.storage.UpdateCounter(m.Name, m.Delta)
	case model.Gauge:
		s.storage.UpdateGauge(m.Name, m.Value)
	}
}

// GetMetric returns metric of specified type and name.
func (s *Service) GetMetric(name string, typ model.MetricType) (model.InternalMetric, bool) {
	var ok bool
	m := model.InternalMetric{
		Name: name,
		Type: typ,
	}
	switch typ {
	case model.Counter:
		m.Delta, ok = s.storage.GetCounter(m.Name)
	case model.Gauge:
		m.Value, ok = s.storage.GetGauge(m.Name)
	}
	return m, ok
}
