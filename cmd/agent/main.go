package main

import (
	"flag"

	"github.com/ashershnyov/go-metrics-gatherer/internal/agent"
)

func main() {
	address := flag.String("a", "http://localhost:8080", "specifies the address for the agent to send metrics to")
	reportInterval := flag.Int("r", 10, "specifies the interval beteween metric sends")
	pollInterval := flag.Int("p", 2, "specifies the interval between metirc gatherings")
	flag.Parse()
	agent := agent.New(
		agent.SetAddress(address),
		agent.SetPollInterval(pollInterval),
		agent.SetReportInterval(reportInterval),
	)
	agent.Run()
}
