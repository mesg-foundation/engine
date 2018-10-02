package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceDevCmdFlags(t *testing.T) {
	c := newServiceDevCmd(nil)

	flags := c.cmd.Flags()
	require.Equal(t, flags.ShorthandLookup("e"), flags.Lookup("event-filter"))
	require.Equal(t, flags.ShorthandLookup("t"), flags.Lookup("task-filter"))
	require.Equal(t, flags.ShorthandLookup("o"), flags.Lookup("output-filter"))

	flags.Set("event-filter", "ef")
	require.Equal(t, "ef", c.eventFilter)

	flags.Set("task-filter", "tf")
	require.Equal(t, "tf", c.taskFilter)

	flags.Set("output-filter", "of")
	require.Equal(t, "of", c.outputFilter)
}
