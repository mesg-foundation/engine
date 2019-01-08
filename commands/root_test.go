// Copyright 2018 MESG Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
