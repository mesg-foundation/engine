package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStartCmdRunE(t *testing.T) {
	m := newMockExecutor()
	c := newStartCmd(m)

	closeStd := captureStd(t)

	m.On("Start").Return(nil)
	c.cmd.Execute()

	stdout, stderr := closeStd()
	require.Contains(t, stdout, "Starting Core")
	require.Contains(t, stdout, "Core started")
	require.Empty(t, stderr)

	m.AssertExpectations(t)
}

func TestStartCmdFlags(t *testing.T) {
	c := newStartCmd(nil)
	require.Equal(t, "text", c.lfv.String())
	require.Equal(t, "info", c.llv.String())

	require.NoError(t, c.cmd.Flags().Set("log-format", "json"))
	require.Equal(t, "json", c.lfv.String())

	require.NoError(t, c.cmd.Flags().Set("log-level", "debug"))
	require.Equal(t, "debug", c.llv.String())
}
