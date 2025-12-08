package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/ashershnyov/go-metrics-gatherer/internal/agent/gatherer"
	"github.com/ashershnyov/go-metrics-gatherer/internal/agent/model"
	"github.com/ashershnyov/go-metrics-gatherer/internal/retrier"
	"github.com/levigross/grequests"
)

const (
	// defaultAddress is a default address for an agent to send metrics to.
	defaultAddress = "http://localhost:8080"
	// defaultPollInterval is a default interval to gather metrics.
	defaultPollInterval = 2 * time.Second
	// deafultReportInterval is a default interval to send metrics to the server.
	deafultReportInterval = 10 * time.Second
	// defaultMaxRetries sets the default amount of retries upon send errors.
	defaultMaxRetries = 3
)

// config stores the server's configuration.
type config struct {
	address        string
	pollInterval   time.Duration
	reportInterval time.Duration
	maxRetries     int
}

// New constructs a config with default values, overrides with opts if passed.
func newConfig(opts ...option) *config {
	c := &config{
		address:        defaultAddress,
		pollInterval:   defaultPollInterval,
		reportInterval: deafultReportInterval,
		maxRetries:     defaultMaxRetries,
	}

	for _, opt := range opts {
		opt(c)
	}

	if !strings.HasPrefix(c.address, "http://") {
		c.address = "http://" + c.address
	}

	return c
}

type option func(*config)

// SetAddress sets custom address for an agent to send metrics to.
func SetAddress(addr *string) option {
	return func(c *config) {
		if addr != nil {
			c.address = *addr
		}
	}
}

// SetPollInterval sets custom polling interval.
func SetPollInterval(t *int) option {
	return func(c *config) {
		if t != nil {
			c.pollInterval = time.Duration(*t * int(time.Second))
		}
	}
}

// SetReportInterval sets custom report interval.
func SetReportInterval(t *int) option {
	return func(c *config) {
		if t != nil {
			c.reportInterval = time.Duration(*t * int(time.Second))
		}
	}
}

// Agent is a client that gathers and sends metrics to the server.
type Agent struct {
	cfg      *config
	gatherer *gatherer.Gatherer
}

// New creates an Agent with a provided cfg.
func New(opts ...option) *Agent {
	return &Agent{
		cfg:      newConfig(opts...),
		gatherer: gatherer.New(),
	}
}

// UpdateMetrics gathers and updates metrics.
func (a *Agent) UpdateMetrics() {
	a.gatherer.Gather()
}

// compressRequest returns gzip-compressed representation of b.
func compressRequest(b []byte) []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(b)
	gz.Close()
	return buf.Bytes()
}

// sendWithOpts sends passed metric to the url with provided opts.
func sendWithOpts(url string, opts grequests.RequestOptions, data any) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if opts.Headers["Content-Encoding"] == "gzip" {
		buf = compressRequest(buf)
	}

	reader := bytes.NewReader(buf)
	opts.RequestBody = io.NopCloser(reader)

	_, err = grequests.Post(url, grequests.FromRequestOptions(&opts))
	if err != nil {
		return err
	}

	return nil
}

// SendMetrics sends all metrics to the server.
func (a *Agent) SendMetrics() error {
	url := a.cfg.address + "/update/"
	opts := grequests.RequestOptions{
		Headers: map[string]string{
			"Content-Type":     "application/json",
			"Content-Encoding": "gzip",
		},
	}

	for name, value := range a.gatherer.GetGauges() {
		metric := model.Metric{
			ID:    name,
			Value: &value,
			Type:  "gauge",
		}

		if err := retrier.WithRetry(a.cfg.maxRetries, func() error { return sendWithOpts(url, opts, metric) }); err != nil {
			return fmt.Errorf("error occured when sending gauges: %w", err)
		}
	}

	for name, value := range a.gatherer.GetCounters() {
		metric := model.Metric{
			ID:    name,
			Delta: &value,
			Type:  "counter",
		}
		if err := retrier.WithRetry(a.cfg.maxRetries, func() error { return sendWithOpts(url, opts, metric) }); err != nil {
			return fmt.Errorf("error occured when sending counters: %w", err)
		}
	}
	return nil
}

// SendMetricsBatch sends all metrics to the server in single request.
func (a *Agent) SendMetricsBatch() error {
	url := a.cfg.address + "/updates/"
	opts := grequests.RequestOptions{
		Headers: map[string]string{
			"Content-Type":     "application/json",
			"Content-Encoding": "gzip",
		},
	}

	gauges := a.gatherer.GetGauges()
	counters := a.gatherer.GetCounters()
	metrics := make([]model.Metric, len(gauges)+len(counters))

	for name, value := range gauges {
		metric := model.Metric{
			ID:    name,
			Value: &value,
			Type:  "gauge",
		}
		metrics = append(metrics, metric)
	}

	for name, value := range counters {
		metric := model.Metric{
			ID:    name,
			Delta: &value,
			Type:  "counter",
		}
		metrics = append(metrics, metric)
	}

	return retrier.WithRetry(a.cfg.maxRetries, func() error { return sendWithOpts(url, opts, metrics) })
}

// Run starts the agent's loops.
func (a *Agent) Run() {
	pollTicker := time.NewTicker(a.cfg.pollInterval)
	reportTicker := time.NewTicker(a.cfg.reportInterval)
	for {
		select {
		case <-pollTicker.C:
			a.UpdateMetrics()
		case <-reportTicker.C:
			if err := a.SendMetricsBatch(); err != nil {
				log.Printf("error occured when sending metrics: %s", err.Error())
			}
		}
	}
}
