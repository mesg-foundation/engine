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
	"os"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
)

type logsCmd struct {
	baseCmd

	e RootExecutor
}

func newLogsCmd(e RootExecutor) *logsCmd {
	c := &logsCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:   "logs",
		Short: "Show the Core's logs",
		RunE:  c.runE,
	})
	return c
}

func (c *logsCmd) runE(cmd *cobra.Command, args []string) error {
	status, err := c.e.Status()
	if err != nil {
		return err
	}

	if status == container.STOPPED {
		fmt.Printf("%s Core is stopped", pretty.WarnSign)
		return nil
	}

	reader, err := c.e.Logs()
	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, reader)
	return err
}
