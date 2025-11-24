package storage

import (
	"testing"
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
