package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
)

type metricService interface {
	ListMetrics() []model.InternalMetric
	UpdateMetric(model.InternalMetric)
	GetMetric(name string, typ model.MetricType) (model.InternalMetric, bool)
}

// MetricsHandler is a handler for updating and getting metrics.
type MetricsHandler struct {
	service metricService
	db      *sql.DB
}

// NewMetricsHandler returns an empty metrics update handler.
func NewMetricsHandler(s metricService, db *sql.DB) *MetricsHandler {
	return &MetricsHandler{
		service: s,
		db:      db,
	}
}

// ListMetrics returns plain test of metrics list in http response.
func (h *MetricsHandler) ListMetrics() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}

			metrics := h.service.ListMetrics()
			var sb strings.Builder

			for _, m := range metrics {
				metricString := "Name:" + m.Name + " Type:" + m.Type + " Value:"
				switch m.Type {
				case model.Gauge:
					metricString += strconv.FormatFloat(m.Value, 'f', -1, 64)
				case model.Counter:
					metricString += strconv.FormatInt(m.Delta, 10)
				}
				sb.WriteString(metricString + "\n")
			}

			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(sb.String()))
		},
	)
}

// UpdateMetricJSON updates value of passed metric.
func (h *MetricsHandler) UpdateMetricJSON() http.HandlerFunc {
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

			h.service.UpdateMetric(m)

			w.WriteHeader(http.StatusOK)
		},
	)
}

// UpdateMetrics updates metrics.
func (h *MetricsHandler) UpdateMetric() http.HandlerFunc {
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

			h.service.UpdateMetric(m)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		},
	)
}

// GetMetricJSON returns value of passed metric.
func (h *MetricsHandler) GetMetricJSON() http.HandlerFunc {
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

			m, _ := h.service.GetMetric(req.ID, req.Type)

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
func (h *MetricsHandler) GetMetric() http.HandlerFunc {
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

			m, ok := h.service.GetMetric(name, typ)
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

// PingDB checks connection to the DB.
func (h *MetricsHandler) PingDB() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			if h.db == nil {
				http.Error(w, "Could not connect to the database", http.StatusInternalServerError)
				return
			}
			err := h.db.PingContext(r.Context())
			if err != nil {
				http.Error(w, "Could not connect to the database", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		},
	)
}
