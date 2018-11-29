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
