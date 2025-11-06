package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
)

type metricService interface {
	ListMetrics() []model.InternalMetric
	UpdateMetric(model.InternalMetric)
}

// MetricDumper dumps metrics to file.
type MetricDumper struct {
	storeInterval time.Duration
	filePath      string
	service       metricService
}

// NewMetricDumper returns a pointer to a newly created MetricDumper and restores metrics from the specified filePath if restoreMetrics is true.
func NewMetricDumper(service metricService, storeInterval time.Duration, filePath string, restoreMetrics bool) (*MetricDumper, error) {
	if restoreMetrics {
		buf, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("an error occurred when creating MetricDumper: %w", err)
		}

		if len(buf) != 0 {
			metrics := []model.Metric{}
			err = json.Unmarshal(buf, &metrics)
			if err != nil {
				return nil, fmt.Errorf("an error occurred when creating MetricDumper: %w", err)
			}

			for _, mExt := range metrics {
				m := model.InternalMetric{
					Name: mExt.ID,
					Type: mExt.Type,
				}
				switch mExt.Type {
				case model.Counter:
					m.Delta = *mExt.Delta
				case model.Gauge:
					m.Value = *mExt.Value
				}
				service.UpdateMetric(m)
			}
		}

	}

	return &MetricDumper{
		storeInterval: storeInterval,
		filePath:      filePath,
		service:       service,
	}, nil
}

// Dump writes metrics to file specified upon d's construction.
func (d *MetricDumper) Dump() error {
	if d.filePath == "" {
		return errors.New("filepath can't be empty")
	}
	metrics := d.service.ListMetrics()
	metricsOut := make([]model.Metric, len(metrics))
	for i, mInt := range metrics {
		mExt := model.Metric{
			ID:   mInt.Name,
			Type: mInt.Type,
		}
		switch mInt.Type {
		case model.Counter:
			mExt.Delta = &mInt.Delta
		case model.Gauge:
			mExt.Value = &mInt.Value
		}
		metricsOut[i] = mExt
	}

	buf, err := json.Marshal(&metricsOut)
	if err != nil {
		return fmt.Errorf("an error occurred while dumping: %w", err)
	}

	file, err := os.Create(d.filePath)
	if err != nil {
		return fmt.Errorf("an error occurred while dumping: %w", err)
	}
	file.Write(buf)
	file.Close()

	return nil
}

// DumperLoop executes dumper's loop.
func (d *MetricDumper) DumperLoop(s metricService) {
	if d.storeInterval <= 0 {
		return
	}
	go func() {
		t := time.NewTicker(d.storeInterval)
		for _ = range t.C {
			d.Dump()
		}
	}()
}

// Middleware is a middleware that dumps metrics to the file if storeInterval > 0.
func (d *MetricDumper) Middleware(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		if d.storeInterval > 0 {
			return
		}
		d.Dump()
	})
}
