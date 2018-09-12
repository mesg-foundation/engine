package commands

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func findCommandChildByUse(root *cobra.Command, use string) bool {
	for _, cmd := range root.Commands() {
		if cmd.Use == use {
			return true
		}
	}
	return false
}

func TestRootCmd(t *testing.T) {
	cmd := Build(nil)
	for _, tt := range []struct {
		use string
	}{
		{"start"},
		{"status"},
		{"stop"},
		{"logs"},
		{"service"},
	} {
		require.Truef(t, findCommandChildByUse(cmd, tt.use), "command %q not found", tt.use)
	}
}
