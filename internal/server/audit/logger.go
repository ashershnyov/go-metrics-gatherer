package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/config"
	"golang.org/x/sync/errgroup"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func getBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func putBuffer(buf *bytes.Buffer) {
	// shouldn't put buffers with huge capacities to the pool
	if buf.Cap() > 4096 {
		return
	}
	bufferPool.Put(buf)
}

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

	buf := getBuffer()
	defer putBuffer(buf)
	encoder := json.NewEncoder(buf)

	if err := encoder.Encode(record); err != nil {
		return fmt.Errorf("error marshalling audit record: %w", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, dst := range l.dsts {
		if !dst.Enabled() {
			continue
		}
		eg.Go(func() error {
			return dst.Log(ctx, bytes.TrimSpace(buf.Bytes()))
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
