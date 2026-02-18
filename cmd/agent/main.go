package main

import (
	"flag"
	"log"

	"github.com/ashershnyov/go-metrics-gatherer/internal/agent"
	"github.com/ashershnyov/go-metrics-gatherer/internal/agent/buildinfo"
	"github.com/ashershnyov/go-metrics-gatherer/internal/agent/config"
	"github.com/caarlos0/env"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

type envVars struct {
	Address        string `env:"ADDRESS" envDefault:""`
	Key            string `env:"KEY" envDefault:""`
	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"-1"`
	PollInterval   int    `env:"POLL_INTERVAL" envDefault:"-1"`
	RateLimit      int    `env:"RATE_LIMIT" envDefault:"1"`
}

func main() {
	var err error

	address := flag.String("a", "http://localhost:8080", "specifies the address for the agent to send metrics to")
	reportInterval := flag.Int("r", 10, "specifies the interval between metric sends")
	pollInterval := flag.Int("p", 2, "specifies the interval between metric gatherings")
	key := flag.String("k", "", "specifies the key to use to hash the request body")
	rateLimit := flag.Int("l", 1, "specifies the maximum amount of parallel requests to the server")
	flag.Parse()

	var envs envVars
	err = env.Parse(&envs)
	if err != nil {
		log.Fatal(err)
	}

	if envs.Address != "" {
		address = &envs.Address
	}

	if envs.PollInterval > 0 {
		pollInterval = &envs.PollInterval
	}

	if envs.ReportInterval > 0 {
		reportInterval = &envs.ReportInterval
	}

	if envs.Key != "" {
		key = &envs.Key
	}

	if envs.RateLimit > 0 {
		rateLimit = &envs.RateLimit
	}

	bi := buildinfo.New(buildVersion, buildDate, buildCommit)

	agent := agent.New(
		bi,
		config.SetAddress(address),
		config.SetPollInterval(pollInterval),
		config.SetReportInterval(reportInterval),
		config.SetKey(key),
		config.SetRateLimit(rateLimit),
	)
	agent.Run()
}
