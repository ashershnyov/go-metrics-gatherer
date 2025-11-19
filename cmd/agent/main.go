package main

import (
	"flag"
	"log"

	"github.com/ashershnyov/go-metrics-gatherer/internal/agent"
	"github.com/caarlos0/env"
)

type envs struct {
	Address        string `env:"ADDRESS" envDefault:""`
	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"-1"`
	PollInterval   int    `env:"POLL_INTERVAL" envDefault:"-1"`
}

func main() {
	var err error

	address := flag.String("a", "http://localhost:8080", "specifies the address for the agent to send metrics to")
	reportInterval := flag.Int("r", 10, "specifies the interval between metric sends")
	pollInterval := flag.Int("p", 2, "specifies the interval between metric gatherings")
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

	agent := agent.New(
		agent.SetAddress(address),
		agent.SetPollInterval(pollInterval),
		agent.SetReportInterval(reportInterval),
	)
	agent.Run()
}
