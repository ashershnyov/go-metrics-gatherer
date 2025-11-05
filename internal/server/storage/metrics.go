package storage

import "sync"

type Gauges map[string]float64
type Counters map[string]int64

// MetricStorage is an in-memory storage for Gauge and Counter metrics.
type MetricStorage struct {
	gauges   Gauges
	counters Counters
	mu       *sync.RWMutex
}

// NewMetricStorage creates an empty MetricStorage.
func NewMetricStorage() *MetricStorage {
	return &MetricStorage{
		gauges:   Gauges{},
		counters: Counters{},
		mu:       &sync.RWMutex{},
	}
}

// GetGauge returns the value of Gauge by the specified name and the indication whether the metric exists.
// If metric does not exist yet, will return (0.0, false).
func (m *MetricStorage) GetGauge(name string) (float64, bool) {
	val, ok := m.gauges[name]
	return val, ok
}

// UpdateGauge replaces the vaule of the Gauge by the specified name.
// Creates a new Gauge with the specified name if it does not exist yet.
func (m *MetricStorage) UpdateGauge(name string, val float64) {
	m.mu.Lock()
	m.gauges[name] = val
	m.mu.Unlock()
}

// GetCounter returns the value of Counter by the specified name and the indication whether the metric exists.
// If metric does not exist yet, will return (0.0, false).
func (m *MetricStorage) GetCounter(name string) (int64, bool) {
	val, ok := m.counters[name]
	return val, ok
}

// UpdateCounter adds val to the Counter with the specified name.
// Creates a new Counter with the specified name if it does not exist yet.
func (m *MetricStorage) UpdateCounter(name string, val int64) {
	_, ok := m.counters[name]
	m.mu.Lock()
	defer m.mu.Unlock()
	if !ok {
		m.counters[name] = val
		return
	}
	m.counters[name] += val
}
