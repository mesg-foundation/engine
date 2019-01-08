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
	"encoding/json"
	"fmt"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
)

type serviceDetailCmd struct {
	baseCmd
	e ServiceExecutor
}

func newServiceDetailCmd(e ServiceExecutor) *serviceDetailCmd {
	c := &serviceDetailCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "detail SERVICE",
		Short:   "Show details of a published service",
		Args:    cobra.ExactArgs(1),
		Example: "mesg-core service detail SERVICE",
		RunE:    c.runE,
	})
	return c
}

func (c *serviceDetailCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		err     error
		service *coreapi.Service
	)
	pretty.Progress("Loading the service...", func() {
		service, err = c.e.ServiceByID(args[0])
	})
	if err != nil {
		return err
	}
	// dump service definition.
	bytes, err := json.Marshal(service)
	if err != nil {
		return err
	}
	fmt.Println(string(pretty.ColorizeJSON(pretty.FgCyan, nil, true, bytes)))
	return nil
}
