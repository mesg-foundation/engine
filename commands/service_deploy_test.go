package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceDeployCmdFlags(t *testing.T) {
	c := newServiceDeployCmd(nil)

	flags := c.cmd.Flags()
	flags.Set("env", "a=1,b=2")
	require.Equal(t, map[string]string{"a": "1", "b": "2"}, c.env)
}
