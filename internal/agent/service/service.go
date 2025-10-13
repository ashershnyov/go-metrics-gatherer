package service

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ashershnyov/go-metrics-gatherer/internal/agent/config"
	"github.com/ashershnyov/go-metrics-gatherer/internal/agent/gatherer"
)

// Agent is a client that gathers and sends metrics to the server.
type Agent struct {
	cfg      *config.Config
	gatherer *gatherer.Gatherer
	client   *http.Client
}

// New creates an Agent with a provided cfg.
func New(cfg *config.Config) *Agent {
	return &Agent{
		cfg:      cfg,
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
		url := a.cfg.Address + "/update/counter/" + name + "/" + strconv.FormatInt(value, 10)
		req, err := prepareRequest(url)

		resp, err := a.client.Do(req)
		if err != nil {
			continue
		}

		resp.Body.Close()
	}
}

func (a *Agent) sendGauges() {
	for name, value := range a.gatherer.GetGauges() {
		url := a.cfg.Address + "/update/gauge/" + name + "/" + strconv.FormatFloat(value, 'f', 2, 64)
		req, err := prepareRequest(url)

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
	pollTicker := time.NewTicker(a.cfg.PollInterval)
	reportTicker := time.NewTicker(a.cfg.ReportInterval)
	for {
		select {
		case <-pollTicker.C:
			a.UpdateMetrics()
		case <-reportTicker.C:
			a.SendMetrics()
		}
	}
}
