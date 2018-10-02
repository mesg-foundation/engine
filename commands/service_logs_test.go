package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceLogsCmdFlags(t *testing.T) {
	c := newServiceLogsCmd(nil)

	flags := c.cmd.Flags()
	require.Equal(t, flags.ShorthandLookup("d"), flags.Lookup("dependencies"))

	flags.Set("dependencies", "a")
	flags.Set("dependencies", "b")
	require.Equal(t, []string{"a", "b"}, c.dependencies)

}
