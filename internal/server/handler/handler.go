package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/service"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/storage"
)

// MetricUpdateHandler a handler for updating metrics.
type MetricUpdateHandler struct {
	// TODO: тоже заменить интерфейсом? тогда где его объявить?
	metrics *storage.MetricStorage
}

// NewMetricUpdateHandler returns an empty metrics update handler.
func NewMetricUpdateHandler() *MetricUpdateHandler {
	return &MetricUpdateHandler{
		metrics: storage.NewMetricStorage(),
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

	m := model.Metric{
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

	service.UpdateMetrc(h.metrics, m)

	w.WriteHeader(http.StatusOK)
}
