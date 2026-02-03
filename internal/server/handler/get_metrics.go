package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
)

// ListMetrics returns plain test of metrics list in http response.
func (h *MetricsHandler) ListMetrics() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			metrics, err := h.service.ListMetrics(r.Context())
			if err != nil {
				http.Error(w, "Could not list metrics", http.StatusInternalServerError)
			}
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

			m, _ := h.service.GetMetric(r.Context(), req.ID, req.Type)

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

			m, err := h.service.GetMetric(r.Context(), name, typ)
			if err != nil {
				http.Error(w, fmt.Sprintf("No such metric: %s", name), http.StatusNotFound)
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
