package storage

import "testing"

func TestUpdateCounter(t *testing.T) {
	storage := NewMetricStorage()

	testVal, testName := int64(123), "Test1"
	storage.UpdateCounter(testName, int64(testVal))
	actual, _ := storage.GetCounter(testName)
	if actual != testVal {
		t.Errorf("got value %v, want %v", actual, testVal)
	}

	storage.UpdateCounter(testName, int64(testVal))
	actual, _ = storage.GetCounter(testName)
	if actual != 2*testVal {
		t.Errorf("got value %v, want %v", actual, 2*testVal)
	}
}

func TestUpdateGauge(t *testing.T) {
	storage := NewMetricStorage()

	testVal, testName := float64(123.123), "Test1"
	storage.UpdateGauge(testName, testVal)
	actual, _ := storage.GetGauge(testName)
	if actual != testVal {
		t.Errorf("got value %v, want %v", actual, testVal)
	}

	testVal = 234
	storage.UpdateGauge(testName, testVal)
	actual, _ = storage.GetGauge(testName)
	if actual != testVal {
		t.Errorf("got value %v, want %v", actual, testVal)
	}
}

func TestGetGauge(t *testing.T) {
	storage := NewMetricStorage()

	testVal, testName := float64(123.123), "Test1"
	storage.UpdateGauge(testName, testVal)
	actual, ok := storage.GetGauge(testName)
	if actual != testVal || !ok {
		t.Errorf("got value %v, want %v", actual, testVal)
	}
	_, ok = storage.GetCounter(testName)
	if ok {
		t.Error("written into wrong metric type")
	}
}

func TestGetCounter(t *testing.T) {
	storage := NewMetricStorage()

	testVal, testName := int64(123), "Test1"
	storage.UpdateCounter(testName, testVal)
	actual, ok := storage.GetCounter(testName)
	if actual != testVal || !ok {
		t.Errorf("got value %v, want %v", actual, testVal)
	}
	_, ok = storage.GetGauge(testName)
	if ok {
		t.Error("written into wrong metric type")
	}
}
