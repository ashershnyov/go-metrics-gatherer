package storage

import (
	"context"
	"fmt"
	"sync"
)

// InMemory is an in-memory storage for Gauge and Counter metrics.
type InMemory struct {
	gauges   Gauges
	counters Counters
	mu       *sync.RWMutex
}

// NewInMemory creates an empty MetricStorage.
func NewInMemory() *InMemory {
	return &InMemory{
		gauges:   Gauges{},
		counters: Counters{},
		mu:       &sync.RWMutex{},
	}
}

// GetGauges returns all gauges stored upon calling.
func (m *InMemory) GetGauges(_ context.Context) (Gauges, error) {
	return m.gauges, nil
}

// GetCounters returns all counters stored upon calling.
func (m *InMemory) GetCounters(_ context.Context) (Counters, error) {
	return m.counters, nil
}

// GetGauge returns the value of Gauge by the specified name and the indication whether the metric exists.
// If metric does not exist yet, will return (0.0, false).
func (m *InMemory) GetGauge(_ context.Context, name string) (float64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.gauges[name]
	var err error
	if !ok {
		err = fmt.Errorf("no such gauge: %s", name)
	}
	return val, err
}

// UpdateGauge replaces the vaule of the Gauge by the specified name.
// Creates a new Gauge with the specified name if it does not exist yet.
func (m *InMemory) UpdateGauge(_ context.Context, name string, val float64) error {
	m.mu.Lock()
	m.gauges[name] = val
	m.mu.Unlock()
	return nil
}

// GetCounter returns the value of Counter by the specified name and the indication whether the metric exists.
// If metric does not exist yet, will return (0.0, false).
func (m *InMemory) GetCounter(_ context.Context, name string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.counters[name]
	var err error
	if !ok {
		err = fmt.Errorf("no such counter: %s", name)
	}
	return val, err
}

// UpdateCounter adds val to the Counter with the specified name.
// Creates a new Counter with the specified name if it does not exist yet.
func (m *InMemory) UpdateCounter(_ context.Context, name string, val int64) error {
	_, ok := m.counters[name]
	m.mu.Lock()
	defer m.mu.Unlock()
	if !ok {
		m.counters[name] = val
		return nil
	}
	m.counters[name] += val
	return nil
}
