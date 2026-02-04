package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
)

// UpdateMetricJSON updates value of passed metric.
func (h *MetricsHandler) UpdateMetricJSON() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			req := model.Metric{}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

			err := h.service.UpdateMetric(r.Context(), m)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			err = h.audit.Log(r.Context(), []string{m.Name}, r.RemoteAddr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w.WriteHeader(http.StatusOK)
		},
	)
}

// UpdateMultipleJSON updates multiple metrics from json body.
func (h *MetricsHandler) UpdateMultipleJSON() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var data []model.Metric
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			metrics := make([]model.InternalMetric, len(data))
			metricsIDs := make([]string, len(data))
			for i, extMetric := range data {
				m := model.InternalMetric{
					Name: extMetric.ID,
					Type: extMetric.Type,
				}

				switch m.Type {
				case model.Counter:
					if extMetric.Delta == nil {
						http.Error(w, "counter must have a set delta", http.StatusBadRequest)
						return
					}
					m.Delta = *extMetric.Delta
				case model.Gauge:
					if extMetric.Value == nil {
						http.Error(w, "gauge must have a set value", http.StatusBadRequest)
						return
					}
					m.Value = *extMetric.Value
				}

				metrics[i] = m
				metricsIDs[i] = extMetric.ID
			}

			err := h.service.UpdateMultipleMetrics(r.Context(), metrics)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = h.audit.Log(r.Context(), metricsIDs, r.RemoteAddr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

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

			err = h.service.UpdateMetric(r.Context(), m)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			err = h.audit.Log(r.Context(), []string{m.Name}, r.RemoteAddr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		},
	)
}
