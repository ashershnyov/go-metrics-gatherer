package config

import "time"

const (
	// defaultAddress specifies the address used unless overridden by starting params.
	defaultAddress = ":8080"
	// defaultStoreInterval specifies the default interval to dump metrics to file.
	defaultStoreInterval = 300 * time.Second
	// defaultFilePath specifies the default path to file to dump metrics to.
	defaultFilePath = "metrics.json"
	// defaultRestoreMetrics specifies the default value of metrics restoration flag.
	defaultRestoreMetrics = true
	// defaultMaxRetries sets the default amount of retries upon send errors.
	defaultMaxRetries = 3
	// defaultKey is a default key value used to hash a response body.
	defaultKey = ""
)

// Config stores the server's configuration.
type Config struct {
	Audit          *Audit
	Address        string
	FilePath       string
	DBAddress      string
	Key            string
	StoreInterval  time.Duration
	MaxRetries     int
	RestoreMetrics bool
}

// NewConfig constructs a config with default values, overrides with opts if passed.
func NewConfig(opts ...Option) *Config {
	c := &Config{
		Address:        defaultAddress,
		StoreInterval:  defaultStoreInterval,
		FilePath:       defaultFilePath,
		RestoreMetrics: defaultRestoreMetrics,
		MaxRetries:     defaultMaxRetries,
		Key:            defaultKey,
		Audit:          &Audit{},
	}

	for _, opt := range opts {
		opt(c)
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

// SetStoreInterval sets custom interval between metric dumps.
func SetStoreInterval(interval *int) Option {
	return func(c *Config) {
		if interval != nil {
			c.StoreInterval = time.Duration(*interval) * time.Second
		}
	}
}

// SetFilePath sets custom path to the file to dump metrics to.
func SetFilePath(path *string) Option {
	return func(c *Config) {
		if path != nil {
			c.FilePath = *path
		}
	}
}

// SetRestoreMetrics sets metric restoration from file flag.
func SetRestoreMetrics(flag *bool) Option {
	return func(c *Config) {
		if flag != nil {
			c.RestoreMetrics = *flag
		}
	}
}

func SetDBAddress(address *string) Option {
	return func(c *Config) {
		if address != nil {
			c.DBAddress = *address
		}
	}
}

func SetKey(key *string) Option {
	return func(c *Config) {
		if key != nil {
			c.Key = *key
		}
	}
}
