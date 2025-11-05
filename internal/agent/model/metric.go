package model

// Metric describes a single metric to communicate with a server.
type Metric struct {
	ID    string   `json:"id"`
	Value *float64 `json:"value,omitempty"`
	Delta *int64   `json:"delta,omitempty"`
	Type  string   `json:"type"`
}
