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
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type serviceDocsCmd struct {
	baseCmd

	force bool
	path  string

	e ServiceExecutor
}

func newServiceDocsCmd(e ServiceExecutor) *serviceDocsCmd {
	c := &serviceDocsCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:   "gen-doc",
		Short: "Generate the documentation for the service in a README.md file",
		Example: `mesg-core service gen-doc
mesg-core service gen-doc ./PATH_TO_SERVICE`,
		PreRunE: c.preRunE,
		RunE:    c.runE,
	})
	return c
}

func (c *serviceDocsCmd) preRunE(cmd *cobra.Command, args []string) error {
	c.path = getFirstOrDefault(args)
	readmePath := filepath.Join(c.path, "README.md")
	if _, err := os.Stat(readmePath); !c.force && err == nil {
		if err := survey.AskOne(&survey.Confirm{
			Message: "The file README.md already exists. Do you want to overwrite it?",
		}, &c.force, nil); err != nil {
			return err
		}
		if !c.force {
			return errors.New("can't continue without confirmation")
		}
	}
	return nil
}

func (c *serviceDocsCmd) runE(cmd *cobra.Command, args []string) error {
	if err := c.e.ServiceGenerateDocs(c.path); err != nil {
		return err
	}

	fmt.Printf("%s File README.md generated with success\n", pretty.SuccessSign)
	return nil
}
