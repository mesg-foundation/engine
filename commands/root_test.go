package commands

import (
	"os"
	"strings"
	"testing"

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func findCommandChildByUsePrefix(root *cobra.Command, use string) bool {
	for _, cmd := range root.Commands() {
		if strings.HasPrefix(cmd.Use, use) {
			return true
		}
	}
	return false
}

func TestMain(m *testing.M) {
	pretty.DisableColor()
	pretty.DisableSpinner()
	os.Exit(m.Run())
}

func TestRootCmd(t *testing.T) {
	cmd := newRootCmd(nil).cmd
	for _, tt := range []struct {
		use string
	}{
		{"start"},
		{"status"},
		{"stop"},
		{"logs"},
		{"service"},
	} {
		require.Truef(t, findCommandChildByUsePrefix(cmd, tt.use), "command %q not found", tt.use)
	}
}

func TestRootCmdFlags(t *testing.T) {
	c := newRootCmd(nil)

	c.cmd.PersistentFlags().Set("no-color", "true")
	require.True(t, c.noColor)

	c.cmd.PersistentFlags().Set("no-spinner", "true")
	require.True(t, c.noSpinner)
}
