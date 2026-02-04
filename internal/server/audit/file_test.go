package audit

import (
	"context"
	"testing"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/config"
	"github.com/stretchr/testify/require"
)

func Test_File_New(t *testing.T) {
	t.Run("No File", func(t *testing.T) {
		cfg := &config.Audit{DstFilePath: "./no_dir/audit.log"}
		f, err := NewFileDst(cfg)
		require.Error(t, err)
		require.False(t, f.Enabled())
	})

	t.Run("Valid file", func(t *testing.T) {
		cfg := &config.Audit{DstFilePath: "./audit.log"}
		f, err := NewFileDst(cfg)
		require.NoError(t, err)
		require.True(t, f.Enabled())
	})
}

func Test_File_Log(t *testing.T) {
	t.Run("No File", func(t *testing.T) {
		cfg := &config.Audit{DstFilePath: "./no_dir/audit.log"}
		f, _ := NewFileDst(cfg)
		err := f.Log(context.Background(), []byte(`{"ts":123324234,"metrics":["TestMetric"],"ip_address":"10.10.10.10"}`))
		require.Error(t, err)
	})

	t.Run("Valid File", func(t *testing.T) {
		cfg := &config.Audit{DstFilePath: "./audit.log"}
		f, _ := NewFileDst(cfg)
		err := f.Log(context.Background(), []byte(`{"ts":123324234,"metrics":["TestMetric"],"ip_address":"10.10.10.10"}`))
		require.NoError(t, err)
	})
}
