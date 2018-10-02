package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceExecuteCmdFlags(t *testing.T) {
	c := newServiceExecuteCmd(nil)

	flags := c.cmd.Flags()
	require.Equal(t, flags.ShorthandLookup("t"), flags.Lookup("task"))
	require.Equal(t, flags.ShorthandLookup("d"), flags.Lookup("data"))
	require.Equal(t, flags.ShorthandLookup("j"), flags.Lookup("json"))

	flags.Set("task", "t")
	require.Equal(t, "t", c.taskKey)

	flags.Set("data", "k=v")
	require.Equal(t, map[string]string{"k": "v"}, c.executeData)

	flags.Set("json", "data.json")
	require.Equal(t, "data.json", c.jsonFile)
}
