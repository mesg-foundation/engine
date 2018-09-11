package commands

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestLogFormatValue(t *testing.T) {
	lfv := logFormatValue("")
	require.Equal(t, lfv.Type(), "string")

	pflag.Var(&lfv, "test-log-format", "")
	require.NoError(t, pflag.Set("test-log-format", "text"))
	require.NoError(t, pflag.Set("test-log-format", "json"))
	require.Error(t, pflag.Set("test-log-format", "unknown"))
}

func TestLogLevelValue(t *testing.T) {
	llv := logLevelValue("")
	require.Equal(t, llv.Type(), "string")

	pflag.Var(&llv, "test-log-level", "")
	require.NoError(t, pflag.Set("test-log-level", "debug"))
	require.NoError(t, pflag.Set("test-log-level", "info"))
	require.NoError(t, pflag.Set("test-log-level", "warn"))
	require.NoError(t, pflag.Set("test-log-level", "error"))
	require.NoError(t, pflag.Set("test-log-level", "fatal"))
	require.NoError(t, pflag.Set("test-log-level", "panic"))
	require.Error(t, pflag.Set("test-log-level", "unknown"))
}
