package config

// Audit defines the configuration of the audit logger.
type Audit struct {
	DstFilePath string
	DstURL      string
}

// SetAuditFilePath sets file path to the file to write audit log to.
func SetAuditFilePath(path *string) Option {
	return func(c *Config) {
		if path != nil {
			c.Audit.DstFilePath = *path
		}
	}
}

// SetAuditURL sets URL to write audit log to.
func SetAuditURL(url *string) Option {
	return func(c *Config) {
		if url != nil {
			c.Audit.DstURL = *url
		}
	}
}
