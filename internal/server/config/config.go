package config

const (
	// defaultAddress specifies the address used unless overridden by starting params.
	defaultAddress = ":8080"
)

// Config stores the server's configuration.
type Config struct {
	Address string
}

// New constructs a config with default values, overrides with opts if passed.
func New() *Config {
	return &Config{
		Address: defaultAddress,
	}
}
