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

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
)

type stopCmd struct {
	baseCmd
	e RootExecutor
}

func newStopCmd(e RootExecutor) *stopCmd {
	c := &stopCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:   "stop",
		Short: "Stop the Core",
		RunE:  c.runE,
	})
	return c
}

func (c *stopCmd) runE(cmd *cobra.Command, args []string) error {
	var err error
	pretty.Progress("Stopping Core...", func() { err = c.e.Stop() })
	if err != nil {
		return err
	}
	fmt.Printf("%s Core stopped\n", pretty.SuccessSign)
	return nil
}
