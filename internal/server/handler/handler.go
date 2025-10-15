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

// TODO: насколько вообще стоит делать такие объекты для хендлеров?
// TODO: стоит ли под каждый хендлер/очень близкую по смыслу группу хендлеров делать свой объект?

// MetricsHandler a handler for updating and getting metrics.
type MetricsHandler struct {
	// TODO: тоже заменить интерфейсом? тогда где его объявить?
	metrics *storage.MetricStorage
}

// NewMetricsHandler returns an empty metrics update handler.
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{
		metrics: storage.NewMetricStorage(),
	}
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
