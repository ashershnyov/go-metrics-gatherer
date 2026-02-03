package config

import (
	"strings"
	"time"
)

const (
	// defaultAddress is a default address for an agent to send metrics to.
	defaultAddress = "http://localhost:8080"
	// defaultPollInterval is a default interval to gather metrics.
	defaultPollInterval = 2 * time.Second
	// deafultReportInterval is a default interval to send metrics to the server.
	deafultReportInterval = 10 * time.Second
	// defaultMaxRetries sets the default amount of retries upon send errors.
	defaultMaxRetries = 3
	// defaultKey is a default key value used to hash a reqeust body.
	defaultKey = ""
	// defaultRateLimit is a default value used to determine max parallel requests count.
	defaultRateLimit = 1
)

// Config stores the server's configuration.
type Config struct {
	Address        string
	PollInterval   time.Duration
	ReportInterval time.Duration
	MaxRetries     int
	Key            string
	RateLimit      int
}

// New constructs a config with default values, overrides with opts if passed.
func New(opts ...Option) *Config {
	c := &Config{
		Address:        defaultAddress,
		PollInterval:   defaultPollInterval,
		ReportInterval: deafultReportInterval,
		MaxRetries:     defaultMaxRetries,
		Key:            defaultKey,
		RateLimit:      defaultRateLimit,
	}

	for _, opt := range opts {
		opt(c)
	}

	if !strings.HasPrefix(c.Address, "http://") {
		c.Address = "http://" + c.Address
	}

	return c
}

type Option func(*Config)

// SetAddress sets custom address for an agent to send metrics to.
func SetAddress(addr *string) Option {
	return func(c *Config) {
		if addr != nil {
			c.Address = *addr
		}
	}
}

// SetPollInterval sets custom polling interval.
func SetPollInterval(t *int) Option {
	return func(c *Config) {
		if t != nil {
			c.PollInterval = time.Duration(*t * int(time.Second))
		}
	}
}

// SetReportInterval sets custom report interval.
func SetReportInterval(t *int) Option {
	return func(c *Config) {
		if t != nil {
			c.ReportInterval = time.Duration(*t * int(time.Second))
		}
	}
}

// SetKey sets the key used to hash the request.
func SetKey(key *string) Option {
	return func(c *Config) {
		if key != nil {
			c.Key = *key
		}
	}
}

// SetRateLimit sets the rate limit value.
func SetRateLimit(rl *int) Option {
	return func(c *Config) {
		if rl != nil {
			c.RateLimit = *rl
		}
	}
}
