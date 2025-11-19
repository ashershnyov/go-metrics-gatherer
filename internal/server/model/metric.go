package model

type MetricType = string

const (
	// Counter that has a new value added every update.
	Counter MetricType = "counter"
	// Gauge has its value repalced every update.
	Gauge MetricType = "gauge"
)

// InternalMetric describes a single metric internally.
type InternalMetric struct {
	Name  string
	Value float64
	Delta int64
	Type  MetricType
}

// Metric describes a single metric to communicate with an agent.
type Metric struct {
	ID    string     `json:"id"`
	Value *float64   `json:"value,omitempty"`
	Delta *int64     `json:"delta,omitempty"`
	Type  MetricType `json:"type"`
}
