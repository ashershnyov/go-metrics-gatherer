package storage

import (
	"testing"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
)

func TestUpdateCounter(t *testing.T) {
	storage := NewInMemory()

	testVal, testName := int64(123), "Test1"
	storage.UpdateCounter(t.Context(), testName, int64(testVal))
	actual, _ := storage.GetCounter(t.Context(), testName)
	if actual != testVal {
		t.Errorf("got value %v, want %v", actual, testVal)
	}

	storage.UpdateCounter(t.Context(), testName, int64(testVal))
	actual, _ = storage.GetCounter(t.Context(), testName)
	if actual != 2*testVal {
		t.Errorf("got value %v, want %v", actual, 2*testVal)
	}
}

func TestUpdateGauge(t *testing.T) {
	storage := NewInMemory()

	testVal, testName := float64(123.123), "Test1"
	storage.UpdateGauge(t.Context(), testName, testVal)
	actual, _ := storage.GetGauge(t.Context(), testName)
	if actual != testVal {
		t.Errorf("got value %v, want %v", actual, testVal)
	}

	testVal = 234
	storage.UpdateGauge(t.Context(), testName, testVal)
	actual, _ = storage.GetGauge(t.Context(), testName)
	if actual != testVal {
		t.Errorf("got value %v, want %v", actual, testVal)
	}
}

func TestGetGauge(t *testing.T) {
	storage := NewInMemory()

	testVal, testName := float64(123.123), "Test1"
	storage.UpdateGauge(t.Context(), testName, testVal)
	actual, err := storage.GetGauge(t.Context(), testName)
	if actual != testVal || err != nil {
		t.Errorf("got value %v, want %v", actual, testVal)
	}
	_, err = storage.GetCounter(t.Context(), testName)
	if err == nil {
		t.Error("written into wrong metric type")
	}
}

func TestGetCounter(t *testing.T) {
	storage := NewInMemory()

	testVal, testName := int64(123), "Test1"
	storage.UpdateCounter(t.Context(), testName, testVal)
	actual, err := storage.GetCounter(t.Context(), testName)
	if actual != testVal || err != nil {
		t.Errorf("got value %v, want %v", actual, testVal)
	}
	_, err = storage.GetGauge(t.Context(), testName)
	if err == nil {
		t.Error("written into wrong metric type")
	}
}

func TestUpdateMultipleMetrics(t *testing.T) {
	storage := NewInMemory()

	metrics := []model.InternalMetric{
		{
			Name:  "Test1",
			Value: 123.123,
			Delta: 0,
			Type:  "gauge",
		},
		{
			Name:  "Test2",
			Value: 0,
			Delta: 123,
			Type:  "counter",
		},
	}

	storage.UpdateMultipleMetrics(t.Context(), metrics)
	wantVal := metrics[0].Value
	actualGauge, _ := storage.GetGauge(t.Context(), metrics[0].Name)
	if actualGauge != metrics[0].Value {
		t.Errorf("got value %v, want %v", actualGauge, wantVal)
	}

	wantDelta := metrics[1].Delta
	actualCounter, _ := storage.GetCounter(t.Context(), metrics[1].Name)
	if actualCounter != wantDelta {
		t.Errorf("got value %v, want %v", actualGauge, wantDelta)
	}
}
