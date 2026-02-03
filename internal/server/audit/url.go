package audit

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/config"
)

// URLDst defines an URL destination for audit logs.
type URLDst struct {
	path string
}

// NewUrlDst returns a new URL destination.
func NewUrlDst(cfg *config.Audit) *URLDst {
	return &URLDst{
		path: cfg.DstURL,
	}
}

// Enabled returns whether to use the URL destination.
func (ud *URLDst) Enabled() bool {
	return ud.path != ""
}

// Log sends audit log to the URL.
func (ud *URLDst) Log(ctx context.Context, record []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ud.path, bytes.NewBuffer(record))
	if err != nil {
		return fmt.Errorf("error sending audit record to url %s : %w", ud.path, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending audit record to url %s : %w", ud.path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK || resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("audit receiver on url %s returned status code %v", ud.path, resp.StatusCode)
	}

	return nil
}

// Close closes the URL destination.
func (ud *URLDst) Close(_ context.Context) error {
	return nil
}
