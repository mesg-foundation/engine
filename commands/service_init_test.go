package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceInitCmdFlags(t *testing.T) {
	c := newServiceInitCmd(nil)

	flags := c.cmd.Flags()
	require.Equal(t, flags.ShorthandLookup("t"), flags.Lookup("template"))

	flags.Set("dir", "/")
	require.Equal(t, "/", c.dir)

	flags.Set("template", "github.com/mesg-foundation/awesome")
	require.Equal(t, "github.com/mesg-foundation/awesome", c.templateURL)
}
