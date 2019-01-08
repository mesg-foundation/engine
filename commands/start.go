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
	"fmt"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
)

type startCmd struct {
	baseCmd

	lfv logFormatValue
	llv logLevelValue

	e RootExecutor
}

func newStartCmd(e RootExecutor) *startCmd {
	c := &startCmd{
		lfv: logFormatValue("text"),
		llv: logLevelValue("info"),
		e:   e,
	}
	c.cmd = newCommand(&cobra.Command{
		Use:     "start",
		Short:   "Start the Core",
		PreRunE: c.preRunE,
		RunE:    c.runE,
	})

	c.cmd.Flags().Var(&c.lfv, "log-format", "log format [text|json]")
	c.cmd.Flags().Var(&c.llv, "log-level", "log level [debug|info|warn|error|fatal|panic]")
	return c
}

func (c *startCmd) preRunE(cmd *cobra.Command, args []string) error {
	cfg, err := config.Global()
	if err != nil {
		return err
	}

	cfg.Log.Format = string(c.lfv)
	cfg.Log.Level = string(c.llv)
	return nil
}

func (c *startCmd) runE(cmd *cobra.Command, args []string) error {
	var err error
	pretty.Progress("Starting Core...", func() { err = c.e.Start() })
	if err != nil {
		return err
	}

	fmt.Printf("%s Core started\n", pretty.SuccessSign)
	return nil
}
