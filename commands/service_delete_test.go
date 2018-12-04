package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceDeleteCmdFlags(t *testing.T) {
	var (
		c     = newServiceDeleteCmd(nil)
		flags = c.cmd.Flags()
	)

	// check defaults.
	require.False(t, c.yes)
	require.False(t, c.all)
	require.False(t, c.keepData)

	require.Equal(t, flags.ShorthandLookup("y"), flags.Lookup("yes"))

	flags.Set("yes", "true")
	require.True(t, c.yes)

	flags.Set("all", "true")
	require.True(t, c.all)

	flags.Set("keep-data", "true")
	require.True(t, c.keepData)
}

func TestServiceDeletePreRunE(t *testing.T) {
	c := newServiceDeleteCmd(nil)
	c.discardOutput()
	require.Equal(t, errNoID, c.preRunE(c.cmd, nil))
}

func TestServiceDeleteRunE(t *testing.T) {
	var (
		m = newMockExecutor()
		c = newServiceDeleteCmd(m)
	)

	tests := []struct {
		all      bool
		keepData bool
		arg      string
	}{
		{all: false, keepData: false, arg: "0"},
		{all: false, keepData: true, arg: "0"},
		{all: true, keepData: false},
		{all: true, keepData: true},
	}

	for _, tt := range tests {
		c.all = tt.all
		c.keepData = tt.keepData

		if tt.all {
			m.On("ServiceDeleteAll", !tt.keepData).Once().Return(nil)
		} else {
			m.On("ServiceDelete", !tt.keepData, tt.arg).Once().Return(nil)
		}
		require.NoError(t, c.runE(c.cmd, []string{tt.arg}))
	}
	m.AssertExpectations(t)
}
