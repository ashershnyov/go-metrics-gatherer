package service

import (
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
)

type metricStorage interface {
	UpdateGauge(name string, val float64)
	UpdateCounter(name string, val int64)
}

// UpdateMetrc updates vales in s using data provided in m.
func UpdateMetrc(s metricStorage, m model.Metric) {
	switch m.Type {
	case model.Counter:
		s.UpdateCounter(m.Name, m.Delta)
	case model.Gauge:
		s.UpdateGauge(m.Name, m.Value)
	}
}
