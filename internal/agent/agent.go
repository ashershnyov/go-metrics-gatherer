package agent

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ashershnyov/go-metrics-gatherer/internal/agent/gatherer"
)

const (
	// defaultAddress is a default address for an agent to send metrics to.
	defaultAddress = "http://localhost:8080"
	// defaultPollInterval is a default interval to gather metrics.
	defaultPollInterval = 2 * time.Second
	// deafultRerportInterval is a default interval to send metrics to the server.
	deafultRerportInterval = 10 * time.Second
)

// config stores the server's configuration.
type config struct {
	address        string
	pollInterval   time.Duration
	reportInterval time.Duration
}

// New constructs a config with default values, overrides with opts if passed.
func newConfig() *config {
	return &config{
		address:        defaultAddress,
		pollInterval:   defaultPollInterval,
		reportInterval: deafultRerportInterval,
	}
}

// Agent is a client that gathers and sends metrics to the server.
type Agent struct {
	cfg      *config
	gatherer *gatherer.Gatherer
	client   *http.Client
}

// New creates an Agent with a provided cfg.
func New() *Agent {
	return &Agent{
		cfg:      newConfig(),
		gatherer: gatherer.New(),
		client:   &http.Client{},
	}
}

// UpdateMetrics gathers and updates metrics.
func (a *Agent) UpdateMetrics() {
	a.gatherer.Gather()
}

func prepareRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Content-Type", "text/plain")
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (a *Agent) sendCounters() {
	for name, value := range a.gatherer.GetCounters() {
		url := a.cfg.address + "/update/counter/" + name + "/" + strconv.FormatInt(value, 10)

		req, err := prepareRequest(url)
		if err != nil {
			continue
		}

		resp, err := a.client.Do(req)
		if err != nil {
			continue
		}

		resp.Body.Close()
	}
}

func (a *Agent) sendGauges() {
	for name, value := range a.gatherer.GetGauges() {
		url := a.cfg.address + "/update/gauge/" + name + "/" + strconv.FormatFloat(value, 'f', 2, 64)
		req, err := prepareRequest(url)
		if err != nil {
			continue
		}

		resp, err := a.client.Do(req)
		if err != nil {
			continue
		}

		resp.Body.Close()
	}
}

// GetGauges sends all metrics to the server.
func (a *Agent) SendMetrics() {
	a.sendCounters()
	a.sendGauges()
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
			a.SendMetrics()
		}
	}
}
