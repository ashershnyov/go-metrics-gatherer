package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
)

const (
	counterMetricType = "counter"
	gaugeMetricType   = "gauge"
)

// MetricUpdateHandler a handler for updating metrics.
type MetricUpdateHandler struct {
	metrics *model.MetricStorage
}

// NewMetricUpdateHandler returns an empty metrics update handler.
func NewMetricUpdateHandler() *MetricUpdateHandler {
	return &MetricUpdateHandler{
		metrics: model.NewMetricStorage(),
	}
}

// ServeHTTP updates metrics.
func (h *MetricUpdateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.Split(strings.TrimPrefix(req.URL.Path, "/update/"), "/")
	if len(path) != 3 {
		http.Error(w, "Invalid URI", http.StatusNotFound)
		return
	}

	metricType, metricName, metricValueString := path[0], path[1], path[2]
	if metricType != counterMetricType && metricType != gaugeMetricType {
		http.Error(w, fmt.Sprintf("No such metric type: %v", metricType), http.StatusBadRequest)
		return
	}

	switch metricType {
	case counterMetricType:
		v, err := strconv.ParseInt(metricValueString, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid metric value: %v for metric: %v", metricValueString, metricName), http.StatusBadRequest)
			return
		}
		h.metrics.UpdateCounter(metricName, v)
	case gaugeMetricType:
		v, err := strconv.ParseFloat(metricValueString, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid metric value: %v for metric: %v", metricValueString, metricName), http.StatusBadRequest)
			return
		}
		h.metrics.UpdateGauge(metricName, v)
	}

	w.WriteHeader(http.StatusOK)
}
