package service

import (
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
)

type MetricStorage interface {
	UpdateGauge(name string, val float64)
	UpdateCounter(name string, val int64)
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
}

// UpdateMetric updates vales in s using data provided in m.
func UpdateMetric(s MetricStorage, m model.Metric) {
	switch m.Type {
	case model.Counter:
		s.UpdateCounter(m.Name, m.Delta)
	case model.Gauge:
		s.UpdateGauge(m.Name, m.Value)
	}
}

// GetMetric returns metric of specified type and name.
func GetMetric(s MetricStorage, name string, typ model.MetricType) (model.Metric, bool) {
	var ok bool
	m := model.Metric{
		Name: name,
		Type: typ,
	}
	switch typ {
	case model.Counter:
		m.Delta, ok = s.GetCounter(m.Name)
	case model.Gauge:
		m.Value, ok = s.GetGauge(m.Name)
	}
	return m, ok
}
