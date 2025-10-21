package agent

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ashershnyov/go-metrics-gatherer/internal/agent/gatherer"
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
	client   *grequests.Session
}

// New creates an Agent with a provided cfg.
func New(opts ...option) *Agent {
	return &Agent{
		cfg:      newConfig(opts...),
		gatherer: gatherer.New(),
		client:   grequests.NewSession(nil),
	}
}

// UpdateMetrics gathers and updates metrics.
func (a *Agent) UpdateMetrics() {
	a.gatherer.Gather()
}

// GetGauges sends all metrics to the server.
func (a *Agent) SendMetrics() error {
	for name, value := range a.gatherer.GetGauges() {
		url := a.cfg.address + "/update/gauge/" + name + "/" + strconv.FormatFloat(value, 'f', 2, 64)
		_, err := a.client.Post(url, nil)
		if err != nil {
			return fmt.Errorf("error occured when sending gauges: %w", err)
		}
	}

	for name, value := range a.gatherer.GetCounters() {
		url := a.cfg.address + "/update/counter/" + name + "/" + strconv.FormatInt(value, 10)
		_, err := a.client.Post(url, nil)
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
				log.Fatalf("error occured when sending metrics: %w", err)
			}
		}
	}
}
