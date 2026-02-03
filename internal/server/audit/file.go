package audit

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/config"
)

// FileDst defines a file destination for audit logs.
type FileDst struct {
	file *os.File
	mu   *sync.Mutex
}

// NewFileDst returns a new file destination.
func NewFileDst(cfg *config.Audit) (*FileDst, error) {
	file, err := os.OpenFile(cfg.DstFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return &FileDst{}, fmt.Errorf("error opening file: %w", err)
	}
	return &FileDst{
		file: file,
		mu:   &sync.Mutex{},
	}, nil
}

// Enabled returns whether to use the file destination.
// Would return false if path is not set.
func (fd *FileDst) Enabled() bool {
	return fd.file != nil
}

// Log logs the passed message to the file.
func (fd *FileDst) Log(ctx context.Context, record []byte) error {
	fd.mu.Lock()
	defer fd.mu.Unlock()

	if _, err := fd.file.Write(append(record, '\n')); err != nil {
		return fmt.Errorf("error writing audit record to file: %w", err)
	}

	return nil
}

// Close closes the file destination.
func (fd *FileDst) Close(_ context.Context) error {
	return fd.file.Close()
}
