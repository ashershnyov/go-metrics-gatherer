package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
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
	// defaultRateLimit is a default value used to determine max parallel requests count.
	defaultRateLimit = 1
)

// Config stores the server's configuration.
type Config struct {
	Address        string        `json:"address"`
	Key            string        `json:"key"`
	PollInterval   time.Duration `json:"poll_interval"`
	ReportInterval time.Duration `json:"report_interval"`
	MaxRetries     int           `json:"-"`
	RateLimit      int           `json:"rate_limit"`
	CryptoKeyPath  string        `json:"crypto_key"`
	cfgFilePath    string
}

func (c *Config) loadFromJSON() error {
	f, err := os.Open(c.cfgFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()

	return decoder.Decode(c)
}

func (c *Config) loadFromFlags() error {
	var (
		address        string
		reportInterval int
		pollInterval   int
		key            string
		rateLimit      int
		cryptoKey      string
	)

	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	fs.StringVar(&address, "a", "http://localhost:8080", "specifies the address for the agent to send metrics to")
	fs.IntVar(&reportInterval, "r", 10, "specifies the interval between metric sends")
	fs.IntVar(&pollInterval, "p", 2, "specifies the interval between metric gatherings")
	fs.StringVar(&key, "k", "", "specifies the key to use to hash the request body")
	fs.IntVar(&rateLimit, "l", 1, "specifies the maximum amount of parallel requests to the server")
	fs.StringVar(&cryptoKey, "crypto-key", "", "specifies the filepath to server's public key")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}

	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "a":
			c.Address = address
		case "r":
			c.ReportInterval = time.Duration(reportInterval) * time.Second
		case "p":
			c.PollInterval = time.Duration(pollInterval) * time.Second
		case "k":
			c.Key = key
		case "l":
			c.RateLimit = rateLimit
		case "crypto-key":
			c.CryptoKeyPath = cryptoKey
		case "c":
		}
	})

	return nil
}

func (c *Config) loadFromEnvs() error {
	if v := os.Getenv("ADDRESS"); v != "" {
		c.Address = v
	}
	if v := os.Getenv("KEY"); v != "" {
		c.Key = v
	}
	if v := os.Getenv("REPORT_INTERVAL"); v != "" {
		interval, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("error parsing REPORT_INTERVAL env: %w", err)
		}
		c.ReportInterval = interval * time.Second
	}
	if v := os.Getenv("POLL_INTERVAL"); v != "" {
		interval, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("error parsing POLL_INTERVAL env: %w", err)
		}
		c.PollInterval = interval * time.Second
	}
	if v := os.Getenv("RATE_LIMIT"); v != "" {
		limit, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("error parsing RATE_LIMIT env: %w", err)
		}
		c.RateLimit = limit
	}
	if v := os.Getenv("CRYPTO_KEY"); v != "" {
		c.CryptoKeyPath = v
	}
	return nil
}

// New constructs a config with default values, overrides with opts if passed.
func New() (*Config, error) {
	c := &Config{
		Address:        defaultAddress,
		PollInterval:   defaultPollInterval,
		ReportInterval: deafultReportInterval,
		MaxRetries:     defaultMaxRetries,
		RateLimit:      defaultRateLimit,
	}

	configPathFlag := ""
	if path := os.Getenv("CONFIG"); path != "" {
		configPathFlag = path
	}
	pathFs := flag.NewFlagSet("config path", flag.ContinueOnError)
	pathFs.StringVar(&configPathFlag, "c", configPathFlag, "specifies path to the config file")
	_ = pathFs.Parse(os.Args[1:])

	c.cfgFilePath = configPathFlag

	var err error

	if c.cfgFilePath != "" {
		err = c.loadFromJSON()
		if err != nil {
			return nil, fmt.Errorf("error loading config: %w", err)
		}
	}

	err = c.loadFromFlags()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	err = c.loadFromEnvs()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	if !strings.HasPrefix(c.Address, "http://") {
		c.Address = "http://" + c.Address
	}

	return c, nil
}
