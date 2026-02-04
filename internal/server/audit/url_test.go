package audit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/config"
	"github.com/stretchr/testify/require"
)

func Test_URL_New(t *testing.T) {
	t.Run("Empty URL", func(t *testing.T) {
		cfg := &config.Audit{DstURL: ""}
		f := NewURLDst(cfg)
		require.False(t, f.Enabled())
	})

	t.Run("Valid URL", func(t *testing.T) {
		cfg := &config.Audit{DstURL: "0.0.0.0:8080/audit"}
		f := NewURLDst(cfg)
		require.True(t, f.Enabled())
	})
}

func Test_URL_Log(t *testing.T) {
	t.Run("Empty URL", func(t *testing.T) {
		cfg := &config.Audit{DstURL: ""}
		f := NewURLDst(cfg)
		err := f.Log(context.Background(), []byte(`{"ts":123324234,"metrics":["TestMetric"],"ip_address":"10.10.10.10"}`))
		require.Error(t, err)
	})

	t.Run("Valid URL", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		cfg := &config.Audit{DstURL: srv.URL}
		f := NewURLDst(cfg)
		err := f.Log(context.Background(), []byte(`{"ts":123324234,"metrics":["TestMetric"],"ip_address":"10.10.10.10"}`))
		require.NoError(t, err)
	})
}
