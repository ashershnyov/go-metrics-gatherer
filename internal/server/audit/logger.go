package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/config"
	"golang.org/x/sync/errgroup"
)

// Destination defines the destination to write audit log to.
type Destination interface {
	Log(context.Context, []byte) error
	Enabled() bool
	Close(context.Context) error
}

// Logger defines audit log writer.
type Logger struct {
	dsts []Destination
}

// NewLogger returns a new logger.
func NewLogger(_ *config.Audit, dsts ...Destination) *Logger {
	return &Logger{
		dsts: dsts,
	}
}

// Log logs the audit record to all enabled destinations.
func (l *Logger) Log(ctx context.Context, metrics []string, ipAddress string) error {
	record := Record{
		TS:        time.Now().Unix(),
		Metrics:   metrics,
		IPAddress: ipAddress,
	}

	bytes, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("error marshalling audit record: %w", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, dst := range l.dsts {
		if !dst.Enabled() {
			continue
		}
		eg.Go(func() error {
			return dst.Log(ctx, bytes)
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("error writing audit log: %w", err)
	}

	return nil
}

// CloseDestinations closes all destinations used by the logger.
func (l *Logger) CloseDestinations() error {
	eg, ctx := errgroup.WithContext(context.Background())
	for _, dst := range l.dsts {
		if !dst.Enabled() {
			continue
		}
		eg.Go(func() error {
			return dst.Close(ctx)
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("error writing audit log: %w", err)
	}

	return nil
}
