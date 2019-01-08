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

	"github.com/spf13/cobra"
)

type serviceValidateCmd struct {
	baseCmd

	path string

	e ServiceExecutor
}

func newServiceValidateCmd(e ServiceExecutor) *serviceValidateCmd {
	c := &serviceValidateCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:   "validate",
		Short: "Validate a service file",
		Long: `Validate a service file. Check the yml format and rules.

All the definitions of the service file can be found in the page [Service File from the documentation](https://docs.mesg.com/guide/service/service-file.html).`,
		Example: `mesg-core service validate
mesg-core service validate ./SERVICE_FOLDER`,
		Args:    cobra.MaximumNArgs(1),
		PreRunE: c.preRunE,
		RunE:    c.runE,
	})
	return c
}

func (c *serviceValidateCmd) preRunE(cmd *cobra.Command, args []string) error {
	c.path = getFirstOrDefault(args)
	return nil
}

func (c *serviceValidateCmd) runE(cmd *cobra.Command, args []string) error {
	msg, err := c.e.ServiceValidate(c.path)
	if err != nil {
		return err
	}
	fmt.Println(msg)
	return nil
}
