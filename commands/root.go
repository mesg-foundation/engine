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
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
)

type rootCmd struct {
	baseCmd

	noColor   bool
	noSpinner bool
}

func newRootCmd(e Executor) *rootCmd {
	c := &rootCmd{}
	c.cmd = newCommand(&cobra.Command{
		Use:              "mesg-core",
		Short:            "MESG Core",
		PersistentPreRun: c.persistentPreRun,
		SilenceUsage:     true,
		SilenceErrors:    true,
	})
	c.cmd.PersistentFlags().BoolVar(&c.noColor, "no-color", c.noColor, "disable colorized output")
	c.cmd.PersistentFlags().BoolVar(&c.noSpinner, "no-spinner", c.noSpinner, "disable spinners")

	c.cmd.AddCommand(
		newStartCmd(e).cmd,
		newStatusCmd(e).cmd,
		newStopCmd(e).cmd,
		newLogsCmd(e).cmd,
		newRootServiceCmd(e).cmd,
	)
	return c
}

func (c *rootCmd) persistentPreRun(cmd *cobra.Command, args []string) {
	if c.noColor {
		pretty.DisableColor()
	}
	if c.noSpinner {
		pretty.DisableSpinner()
	}
}
