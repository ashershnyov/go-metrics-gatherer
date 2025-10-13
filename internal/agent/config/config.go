package config

import (
	"time"
)

const (
	// defaultAddress is a default address for an agent to send metrics to.
	defaultAddress = "http://localhost:8080"
	// defaultPollInterval is a default interval to gather metrics.
	defaultPollInterval = 2 * time.Second
	// deafultRerportInterval is a default interval to send metrics to the server.
	deafultRerportInterval = 10 * time.Second
)

// Config stores the server's configuration.
type Config struct {
	Address        string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

// New constructs a config with default values, overrides with opts if passed.
func New() *Config {
	return &Config{
		Address:        defaultAddress,
		PollInterval:   defaultPollInterval,
		ReportInterval: deafultRerportInterval,
	}
}
