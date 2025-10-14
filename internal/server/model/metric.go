package model

type MetricType string

const (
	// Counter that has a new value added every update.
	Counter MetricType = "counter"
	// Gauge has its value repalced every update.
	Gauge MetricType = "gauge"
)

// Metric describes a single metric.
type Metric struct {
	Name  string
	Value float64
	Delta int64
	Type  MetricType
}
