package service

import (
	"testing"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/storage"
)

func TestUpdateMetric(t *testing.T) {
	cases := []struct {
		actual model.Metric
		want   model.Metric
	}{
		{
			model.Metric{Name: "Test1", Value: 12.5, Delta: 0, Type: model.Gauge},
			model.Metric{Name: "Test1", Value: 12.5, Delta: 0, Type: model.Gauge},
		},
		{
			model.Metric{Name: "Test1", Value: 22.0, Delta: 0, Type: model.Gauge},
			model.Metric{Name: "Test1", Value: 22.0, Delta: 0, Type: model.Gauge},
		},
		{
			model.Metric{Name: "Test2", Value: 0, Delta: 12, Type: model.Counter},
			model.Metric{Name: "Test2", Value: 0, Delta: 12, Type: model.Counter},
		},
		{
			model.Metric{Name: "Test2", Value: 0, Delta: 12, Type: model.Counter},
			model.Metric{Name: "Test2", Value: 0, Delta: 24, Type: model.Counter},
		},
	}

	storage := storage.NewMetricStorage()

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			UpdateMetrc(storage, tt.actual)
			switch tt.actual.Type {
			case model.Gauge:
				v, _ := storage.GetGauge(tt.want.Name)
				if v != tt.want.Value {
					t.Errorf("got value %v, want %v", v, tt.actual.Value)
				}
			case model.Counter:
				d, _ := storage.GetCounter(tt.want.Name)
				if d != tt.want.Delta {
					t.Errorf("got value %v, want %v", d, tt.actual.Delta)
				}
			}
		})
	}
}
