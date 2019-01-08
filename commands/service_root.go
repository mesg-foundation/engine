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
	"github.com/spf13/cobra"
)

type rootServiceCmd struct {
	baseCmd
}

func newRootServiceCmd(e ServiceExecutor) *rootServiceCmd {
	c := &rootServiceCmd{}
	c.cmd = newCommand(&cobra.Command{
		Use:   "service",
		Short: "Manage services",
	})

	c.cmd.AddCommand(
		newServiceDeployCmd(e).cmd,
		newServiceValidateCmd(e).cmd,
		newServiceStartCmd(e).cmd,
		newServiceStopCmd(e).cmd,
		newServiceDetailCmd(e).cmd,
		newServiceListCmd(e).cmd,
		newServiceInitCmd(e).cmd,
		newServiceDeleteCmd(e).cmd,
		newServiceLogsCmd(e).cmd,
		newServiceDocsCmd(e).cmd,
		newServiceDevCmd(e).cmd,
		newServiceExecuteCmd(e).cmd,
	)
	return c
}
