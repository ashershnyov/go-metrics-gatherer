package main

import (
	"flag"
	"log"

	"github.com/ashershnyov/go-metrics-gatherer/internal/agent"
	"github.com/ashershnyov/go-metrics-gatherer/internal/agent/config"
	"github.com/caarlos0/env"
)

type envs struct {
	Address        string `env:"ADDRESS" envDefault:""`
	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"-1"`
	PollInterval   int    `env:"POLL_INTERVAL" envDefault:"-1"`
	Key            string `env:"KEY" envDefault:""`
}

func main() {
	var err error

	address := flag.String("a", "http://localhost:8080", "specifies the address for the agent to send metrics to")
	reportInterval := flag.Int("r", 10, "specifies the interval between metric sends")
	pollInterval := flag.Int("p", 2, "specifies the interval between metric gatherings")
	key := flag.String("k", "", "specifies the key to use to hash the request body")
	flag.Parse()

	var envs envs
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

	agent := agent.New(
		config.SetAddress(address),
		config.SetPollInterval(pollInterval),
		config.SetReportInterval(reportInterval),
		config.SetKey(key),
	)
	agent.Run()
}
