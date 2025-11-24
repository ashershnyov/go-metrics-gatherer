package service

import (
	"errors"
	"testing"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/storage"
)

func TestUpdateMetric(t *testing.T) {
	cases := []struct {
		actual model.InternalMetric
		want   model.InternalMetric
	}{
		{
			model.InternalMetric{Name: "Test1", Value: 12.5, Delta: 0, Type: model.Gauge},
			model.InternalMetric{Name: "Test1", Value: 12.5, Delta: 0, Type: model.Gauge},
		},
		{
			model.InternalMetric{Name: "Test1", Value: 22.0, Delta: 0, Type: model.Gauge},
			model.InternalMetric{Name: "Test1", Value: 22.0, Delta: 0, Type: model.Gauge},
		},
		{
			model.InternalMetric{Name: "Test2", Value: 0, Delta: 12, Type: model.Counter},
			model.InternalMetric{Name: "Test2", Value: 0, Delta: 12, Type: model.Counter},
		},
		{
			model.InternalMetric{Name: "Test2", Value: 0, Delta: 12, Type: model.Counter},
			model.InternalMetric{Name: "Test2", Value: 0, Delta: 24, Type: model.Counter},
		},
	}

	storage := storage.NewInMemory()
	service := NewService(storage)

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			service.UpdateMetric(t.Context(), tt.actual)
			switch tt.actual.Type {
			case model.Gauge:
				v, _ := storage.GetGauge(t.Context(), tt.want.Name)
				if v != tt.want.Value {
					t.Errorf("got value %v, want %v", v, tt.actual.Value)
				}
			case model.Counter:
				d, _ := storage.GetCounter(t.Context(), tt.want.Name)
				if d != tt.want.Delta {
					t.Errorf("got value %v, want %v", d, tt.actual.Delta)
				}
			}
		})
	}
}

func TestGetMetric(t *testing.T) {
	storage := storage.NewInMemory()
	storage.UpdateCounter(t.Context(), "TestCounter", 123)
	storage.UpdateGauge(t.Context(), "TestGauge", 234.123)
	service := NewService(storage)

	cases := []struct {
		name    string
		typ     model.MetricType
		want    model.InternalMetric
		wantErr error
	}{
		{
			name: "TestCounter",
			typ:  model.Counter,
			want: model.InternalMetric{
				Name:  "TestCounter",
				Delta: 123,
				Type:  model.Counter,
			},
			wantErr: nil,
		},
		{
			name: "TestGauge",
			typ:  model.Gauge,
			want: model.InternalMetric{
				Name:  "TestGauge",
				Value: 234.123,
				Type:  model.Gauge,
			},
			wantErr: nil,
		},
		{
			name: "TestFail",
			typ:  model.Gauge,
			want: model.InternalMetric{
				Name: "TestFail",
				Type: model.Gauge,
			},
			wantErr: errors.New("no such gauge: TestFail"),
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			m, err := service.GetMetric(t.Context(), tt.name, tt.typ)
			if err != tt.wantErr && err.Error() != tt.wantErr.Error() {
				t.Error("metric present where it shouldn't be or vice-versa")
			}
			if m.Name != tt.want.Name || m.Type != tt.want.Type || m.Value != tt.want.Value || m.Delta != tt.want.Delta {
				t.Errorf("want %v, got %v", tt.want, m)
			}
		})
	}
}
