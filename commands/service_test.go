package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootServiceCmd(t *testing.T) {
	cmd := newRootServiceCmd(nil).cmd
	for _, tt := range []struct {
		use string
	}{
		{"deploy"},
		{"validate"},
		{"start"},
		{"stop"},
		{"detail"},
		{"list"},
		{"init"},
		{"delete"},
		{"logs"},
		{"gen-doc"},
		{"dev"},
		{"execute"},
	} {
		require.Truef(t, findCommandChildByUsePrefix(cmd, tt.use), "command %q not found", tt.use)
	}
}
