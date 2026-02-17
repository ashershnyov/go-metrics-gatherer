package audit

// Record defines the audit log record.
type Record struct {
	IPAddress string   `json:"ip_address"`
	Metrics   []string `json:"metrics"`
	TS        int64    `json:"ts"`
}
