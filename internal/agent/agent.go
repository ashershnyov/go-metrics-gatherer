package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ashershnyov/go-metrics-gatherer/internal/agent/gatherer"
	"github.com/ashershnyov/go-metrics-gatherer/internal/agent/model"
	"github.com/levigross/grequests"
)

const (
	// defaultAddress is a default address for an agent to send metrics to.
	defaultAddress = "http://localhost:8080"
	// defaultPollInterval is a default interval to gather metrics.
	defaultPollInterval = 2 * time.Second
	// deafultReportInterval is a default interval to send metrics to the server.
	deafultReportInterval = 10 * time.Second
)

// config stores the server's configuration.
type config struct {
	address        string
	pollInterval   time.Duration
	reportInterval time.Duration
}

// New constructs a config with default values, overrides with opts if passed.
func newConfig(opts ...option) *config {
	c := &config{
		address:        defaultAddress,
		pollInterval:   defaultPollInterval,
		reportInterval: deafultReportInterval,
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

// SendMetrics sends all metrics to the server.
func (a *Agent) SendMetrics() error {
	opts := &grequests.RequestOptions{
		Headers: map[string]string{"Content-Type": "application/json"},
	}
	url := a.cfg.address + "/update/"

	for name, value := range a.gatherer.GetGauges() {
		req := model.Metric{
			ID:    name,
			Value: &value,
			Type:  "gauge",
		}

		buf, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("error occured when sending gauges: %w", err)
		}

		reader := bytes.NewReader(buf)
		opts.RequestBody = reader

		_, err = grequests.Post(url, grequests.FromRequestOptions(opts))
		if err != nil {
			return fmt.Errorf("error occured when sending gauges: %w", err)
		}
	}

	for name, value := range a.gatherer.GetCounters() {
		req := model.Metric{
			ID:    name,
			Delta: &value,
			Type:  "counter",
		}

		buf, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("error occured when sending gauges: %w", err)
		}

		reader := bytes.NewReader(buf)
		opts.RequestBody = reader

		_, err = grequests.Post(url, grequests.FromRequestOptions(opts))
		if err != nil {
			return fmt.Errorf("error when sending counters: %w", err)
		}
	}
	return nil
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
			if err := a.SendMetrics(); err != nil {
				log.Printf("error occured when sending metrics: %w", err)
			}
		}
	}
}
