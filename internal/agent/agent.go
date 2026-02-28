package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/ashershnyov/go-metrics-gatherer/internal/buildinfo"

	"github.com/ashershnyov/go-metrics-gatherer/internal/agent/config"
	"github.com/ashershnyov/go-metrics-gatherer/internal/agent/gatherer"
	"github.com/ashershnyov/go-metrics-gatherer/internal/agent/model"
	hg "github.com/ashershnyov/go-metrics-gatherer/pkg/hasher"
	"github.com/ashershnyov/go-metrics-gatherer/pkg/retrier"
	"github.com/levigross/grequests"
)

type hasher interface {
	Hash([]byte) string
}

// Agent is a client that gathers and sends metrics to the server.
type Agent struct {
	buildinfo buildinfo.BuildInfo
	cfg       *config.Config
	gatherer  *gatherer.Gatherer
	hasher    hasher
}

// New creates an Agent with a provided cfg.
func New(bi buildinfo.BuildInfo, opts ...config.Option) *Agent {
	cfg := config.New(opts...)
	return &Agent{
		buildinfo: bi,
		cfg:       cfg,
		gatherer:  gatherer.New(),
		hasher:    hg.NewHasher(cfg.Key),
	}
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
func (a *Agent) sendWithOpts(url string, opts grequests.RequestOptions, data any) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if a.cfg.Key != "" {
		opts.Headers["HashSHA256"] = a.hasher.Hash(buf)
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
	url := a.cfg.Address + "/update/"
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

		if err := retrier.WithRetry(a.cfg.MaxRetries, func() error { return a.sendWithOpts(url, opts, metric) }); err != nil {
			return fmt.Errorf("error occured when sending gauges: %w", err)
		}
	}

	for name, value := range a.gatherer.GetCounters() {
		metric := model.Metric{
			ID:    name,
			Delta: &value,
			Type:  "counter",
		}
		if err := retrier.WithRetry(a.cfg.MaxRetries, func() error { return a.sendWithOpts(url, opts, metric) }); err != nil {
			return fmt.Errorf("error occured when sending counters: %w", err)
		}
	}
	return nil
}

// SendMetricsBatch sends all metrics to the server in single request.
func (a *Agent) SendMetricsBatch() error {
	url := a.cfg.Address + "/updates/"
	opts := grequests.RequestOptions{
		Headers: map[string]string{
			"Content-Type":     "application/json",
			"Content-Encoding": "gzip",
		},
	}

	gauges := a.gatherer.GetGauges()
	counters := a.gatherer.GetCounters()
	metrics := make([]model.Metric, 0, len(gauges)+len(counters))

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

	return retrier.WithRetry(a.cfg.MaxRetries, func() error { return a.sendWithOpts(url, opts, metrics) })
}

func (a *Agent) metricSender(jobs <-chan struct{}, errs chan<- error) {
	for range jobs {
		if err := a.SendMetricsBatch(); err != nil {
			errs <- err
		}
	}
}

// Run starts the agent's loops.
func (a *Agent) Run() {
	fmt.Println(a.buildinfo.String())

	go a.gatherer.GatherAndUpdateLoop(a.cfg.PollInterval)

	jobsChan := make(chan struct{}, a.cfg.RateLimit)
	errChan := make(chan error, a.cfg.RateLimit)
	for i := 0; i < a.cfg.RateLimit; i++ {
		go a.metricSender(jobsChan, errChan)
	}

	reportTicker := time.NewTicker(a.cfg.ReportInterval)
	for {
		select {
		case <-reportTicker.C:
			jobsChan <- struct{}{}
		case err := <-errChan:
			if err != nil {
				log.Printf("error sending metrics: %s", err.Error())
			}
		}
	}
}
