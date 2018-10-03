package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceDeleteCmdFlags(t *testing.T) {
	c := newServiceDeleteCmd(nil)

	flags := c.cmd.Flags()
	require.Equal(t, flags.ShorthandLookup("f"), flags.Lookup("force"))

	flags.Set("force", "true")
	require.True(t, c.force)

	flags.Set("all", "true")
	require.True(t, c.all)
}
