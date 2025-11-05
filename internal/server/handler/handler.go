package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/service"
)

// MetricsHandler a handler for updating and getting metrics.
type MetricsHandler struct {
	metrics service.MetricStorage
}

// NewMetricsHandler returns an empty metrics update handler.
func NewMetricsHandler(s service.MetricStorage) *MetricsHandler {
	return &MetricsHandler{
		metrics: s,
	}
}

// UpdateMetricJSON updates value of passed metric.
func (h *MetricsHandler) UpdateMetricJSON() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var buf bytes.Buffer
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			req := &model.Metric{}
			if err := json.Unmarshal(buf.Bytes(), req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			m := model.InternalMetric{
				Name: req.ID,
				Type: req.Type,
			}

			switch m.Type {
			case model.Counter:
				if req.Delta == nil {
					http.Error(w, "counter must have a set delta", http.StatusBadRequest)
					return
				}
				m.Delta = *req.Delta
			case model.Gauge:
				if req.Value == nil {
					http.Error(w, "gauge must have a set value", http.StatusBadRequest)
					return
				}
				m.Value = *req.Value
			}

			service.UpdateMetric(h.metrics, m)

			w.WriteHeader(http.StatusOK)
		},
	)
}

// UpdateMetrics updates metrics.
func (h *MetricsHandler) UpdateMetric() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			path := strings.Split(strings.TrimPrefix(r.URL.Path, "/update/"), "/")
			if len(path) != 3 {
				http.Error(w, "Invalid URI", http.StatusNotFound)
				return
			}

			m := model.InternalMetric{
				Name: path[1],
				Type: model.MetricType(path[0]),
			}

			if m.Type != model.Counter && m.Type != model.Gauge {
				http.Error(w, fmt.Sprintf("No such metric type: %v", m.Type), http.StatusBadRequest)
				return
			}

			var err error

			rawMetricValue := path[2]

			switch m.Type {
			case model.Counter:
				m.Delta, err = strconv.ParseInt(rawMetricValue, 10, 64)
				if err != nil {
					http.Error(w, fmt.Sprintf("Invalid metric value: %v for metric: %v", rawMetricValue, m.Name), http.StatusBadRequest)
					return
				}
			case model.Gauge:
				m.Value, err = strconv.ParseFloat(rawMetricValue, 64)
				if err != nil {
					http.Error(w, fmt.Sprintf("Invalid metric value: %v for metric: %v", rawMetricValue, m.Name), http.StatusBadRequest)
					return
				}
			}

			service.UpdateMetric(h.metrics, m)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		},
	)
}

// GetMetricJSON returns value of passed metric.
func (h *MetricsHandler) GetMetricJSON() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var buf bytes.Buffer
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			req := &model.Metric{}
			if err := json.Unmarshal(buf.Bytes(), req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			m, _ := service.GetMetric(h.metrics, req.ID, req.Type)

			switch m.Type {
			case model.Counter:
				req.Delta = &m.Delta
			case model.Gauge:
				req.Value = &m.Value
			}

			w.Header().Set("Content-Type", "application/json")
			resp, err := json.Marshal(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		},
	)
}

// GetMetric returns metric's value.
func (h *MetricsHandler) GetMetric() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			path := strings.Split(strings.TrimPrefix(r.URL.Path, "/value/"), "/")
			if len(path) != 2 {
				http.Error(w, "Invalid URI", http.StatusNotFound)
				return
			}

			name, typ := path[1], model.MetricType(path[0])

			if typ != model.Counter && typ != model.Gauge {
				http.Error(w, fmt.Sprintf("No such metric type: %v", typ), http.StatusBadRequest)
				return
			}

			m, ok := service.GetMetric(h.metrics, name, typ)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusOK)
			switch m.Type {
			case model.Counter:
				w.Write([]byte(strconv.FormatInt(m.Delta, 10)))
			case model.Gauge:
				w.Write([]byte(strconv.FormatFloat(m.Value, 'f', -1, 64)))
			}
		},
	)
}
