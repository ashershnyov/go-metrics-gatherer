package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

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
)

type cidr struct {
	*net.IPNet
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (c *cidr) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	_, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		return fmt.Errorf("invalid CIDR %q: %w", s, err)
	}

	c.IPNet = ipNet
	return nil
}

// Audit defines the configuration of the audit logger.
type Audit struct {
	DstFilePath string `json:"file_path"`
	DstURL      string `json:"url"`
}

// Config stores the server's configuration.
type Config struct {
	Audit          *Audit        `json:"audit"`
	Address        string        `json:"address"`
	FilePath       string        `json:"file_store_path"`
	DBAddress      string        `json:"database_dsn"`
	Key            string        `json:"key"`
	CryptoKeyPath  string        `json:"crypto_key_path"`
	StoreInterval  time.Duration `json:"store_interval"`
	MaxRetries     int           `json:"-"`
	RestoreMetrics bool          `json:"restore"`
	TrustedSubnet  cidr          `json:"trusted_subnet"`
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
		address       string
		storeInterval int
		filePath      string
		restore       bool
		dbAddress     string
		key           string
		auditURL      string
		auditFile     string
		cryptoKey     string
		cidrStr       string
	)

	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	fs.StringVar(&address, "a", "localhost:8080", "specifies the address for the server to start on")
	fs.IntVar(&storeInterval, "i", int(300*time.Second), "specifies the interval between writes to the specified file")
	fs.StringVar(&filePath, "f", "metrics.json", "specifies the filepath to store metric values in")
	fs.BoolVar(&restore, "r", true, "indicates whether the stored metrics should be loaded from the specified file on server startup")
	fs.StringVar(&dbAddress, "d", "", "specifies the address of the DB")
	fs.StringVar(&key, "k", "", "specifies the key to use to hash the response body")
	fs.StringVar(&auditURL, "audit-url", "", "specifies the URL to send audit logs to")
	fs.StringVar(&auditFile, "audit-file", "", "specifies the filepath to write audit logs to")
	fs.StringVar(&cryptoKey, "crypto-key", "", "specifies the filepath to server's private key")
	fs.StringVar(&cidrStr, "t", "", "specifies the subnet of ip's to accept requests from")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}

	var err error

	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "a":
			c.Address = address
		case "i":
			c.StoreInterval = time.Duration(storeInterval) * time.Second
		case "f":
			c.FilePath = filePath
		case "r":
			c.RestoreMetrics = restore
		case "d":
			c.DBAddress = dbAddress
		case "k":
			c.Key = key
		case "audit-url":
			c.Audit.DstURL = auditURL
		case "audit-file":
			c.Audit.DstFilePath = auditFile
		case "crypto-key":
			c.CryptoKeyPath = cryptoKey
		case "t":
			var cidrNet *net.IPNet
			_, cidrNet, err = net.ParseCIDR(cidrStr)
			c.TrustedSubnet.IPNet = cidrNet
		case "c":
		}
	})

	if err != nil {
		return fmt.Errorf("error parsing `-t` flag value: %w", err)
	}

	return nil
}

func (c *Config) loadFromEnvs() error {
	if v := os.Getenv("ADDRESS"); v != "" {
		c.Address = v
	}
	if v := os.Getenv("FILE_STORAGE_PATH"); v != "" {
		c.FilePath = v
	}
	if v := os.Getenv("RESTORE"); v != "" {
		c.RestoreMetrics = v == "true"
	}
	if v := os.Getenv("DATABASE_DSN"); v != "" {
		c.DBAddress = v
	}
	if v := os.Getenv("KEY"); v != "" {
		c.Key = v
	}
	if v := os.Getenv("AUDIT_URL"); v != "" {
		c.Audit.DstURL = v
	}
	if v := os.Getenv("AUDIT_FILE"); v != "" {
		c.Audit.DstFilePath = v
	}
	if v := os.Getenv("STORE_INTERVAL"); v != "" {
		interval, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("error parsing STORE_INTERVAL env: %w", err)
		}
		c.StoreInterval = time.Duration(interval) * time.Second
	}
	if v := os.Getenv("CRYPTO_KEY"); v != "" {
		c.CryptoKeyPath = v
	}
	if v := os.Getenv("TRUSTED_SUBNET"); v != "" {
		_, cidrNet, err := net.ParseCIDR(v)
		if err != nil {
			return fmt.Errorf("error parsing `TRUSTED_SUBNET` env: %w", err)
		}
		c.TrustedSubnet.IPNet = cidrNet
	}
	return nil
}

// NewConfig constructs a config with default values, overrides with opts if passed.
func NewConfig() (*Config, error) {
	c := &Config{
		Address:        defaultAddress,
		StoreInterval:  defaultStoreInterval,
		FilePath:       defaultFilePath,
		RestoreMetrics: defaultRestoreMetrics,
		MaxRetries:     defaultMaxRetries,
		Audit:          &Audit{},
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

	return c, nil
}
