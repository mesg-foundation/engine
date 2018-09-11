package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type serviceDeleteCmd struct {
	baseCmd

	all   bool
	force bool

	e ServiceExecutor
}

func newServiceDeleteCmd(e ServiceExecutor) *serviceDeleteCmd {
	c := &serviceDeleteCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete one or many services",
		Example: `mesg-core service delete SERVICE_ID [SERVICE_ID...]
mesg-core service delete --all`,
		PreRunE: c.preRunE,
		RunE:    c.runE,
	})

	c.cmd.Flags().BoolVar(&c.all, "all", c.all, "Delete all services")
	c.cmd.Flags().BoolVarP(&c.force, "force", "f", c.force, "Force delete all services")
	return c
}

func (c *serviceDeleteCmd) preRunE(cmd *cobra.Command, args []string) error {
	if len(args) == 0 && !c.all {
		return errors.New("at least one service id must be provided (or run with --all flag)")
	}

	if !c.all || (c.all && c.force) {
		return nil
	}

	if err := survey.AskOne(&survey.Confirm{Message: "Are you sure to delete all services?"}, &c.force, nil); err != nil {
		return err
	}

	// is still no confirm .
	if !c.force {
		return errors.New("can't continue without confirmation")
	}
	return nil
}

func (c *serviceDeleteCmd) runE(cmd *cobra.Command, args []string) error {
	if c.all {
		if err := c.e.ServiceDeleteAll(); err != nil {
			return err
		}
		fmt.Printf("%s all services are deleted", pretty.SuccessSign)
		return nil
	}

	exitWithError := false
	for _, arg := range args {
		if err := c.e.ServiceDelete(arg); err != nil {
			exitWithError = true
			fmt.Fprintf(os.Stderr, "%s can't delete %s service: %s\n", pretty.FailSign, arg, err)
		} else {
			fmt.Printf("%s service %s deleted\n", pretty.SuccessSign, arg)
		}
	}

	if exitWithError {
		return errors.New("There was a problem with deleting some services")
	}
	return nil
}
