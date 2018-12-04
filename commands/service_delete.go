package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

var (
	errConfirmationNeeded = errors.New("can't continue without confirmation")
	errInputsNeededErr    = errors.New("at least one service id must be provided (or run with --all flag)")
)

type serviceDeleteCmd struct {
	baseCmd

	yes               bool
	deleteAllServices bool
	keepData          bool

	e      ServiceExecutor
	survey Survey
}

func newServiceDeleteCmd(e ServiceExecutor, survey Survey) *serviceDeleteCmd {
	c := &serviceDeleteCmd{
		e:      e,
		survey: survey,
	}
	c.cmd = newCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete one or many services",
		Example: `mesg-core service delete SERVICE [SERVICE...]
mesg-core service delete --all`,
		PreRunE: c.preRunE,
		RunE:    c.runE,
	})
	c.cmd.Flags().BoolVarP(&c.yes, "yes", "y", c.yes, `Automatically say 'Yes' to all prompts and run non-interactively`)
	c.cmd.Flags().BoolVar(&c.deleteAllServices, "all", c.deleteAllServices, "Delete all services")
	c.cmd.Flags().BoolVar(&c.keepData, "keep-data", c.keepData, "Delete services' persistent data along the way")
	return c
}

func (c *serviceDeleteCmd) preRunE(cmd *cobra.Command, args []string) error {
	if err := c.validateInputs(args); err != nil {
		return err
	}

	if c.yes {
		return nil
	}

	if err := c.confirmServiceDelete(); err != nil {
		return err
	}

	if !c.keepData {
		if err := c.confirmDataDelete(); err != nil {
			return err
		}
	}

	return nil
}

// validateInputs validates command arguments and flags.
func (c *serviceDeleteCmd) validateInputs(args []string) error {
	if len(args) == 0 && !c.deleteAllServices {
		return errInputsNeededErr
	}
	return nil
}

// confirmServiceDelete prompts a confirmation dialog for deleting services.
func (c *serviceDeleteCmd) confirmServiceDelete() error {
	var (
		confirmed bool
		confirm   = &survey.Confirm{Default: false}
	)

	if c.deleteAllServices {
		confirm.Message = "Are you sure to delete all services?"
	} else {
		confirm.Message = "Are you sure to delete service(s)?"
	}

	if err := c.survey.AskOne(confirm, &confirmed, nil); err != nil {
		return err
	}

	if !confirmed {
		return errConfirmationNeeded
	}
	return nil
}

// confirmServiceDelete prompts a confirmation dialog for deleting services' data.
func (c *serviceDeleteCmd) confirmDataDelete() error {
	if err := c.survey.AskOne(&survey.Confirm{
		Message: "Do you want to remove service(s)' persistent data as well?",
		Default: false,
	}, &c.keepData, nil); err != nil {
		return err
	}

	c.keepData = !c.keepData
	return nil
}

func (c *serviceDeleteCmd) runE(cmd *cobra.Command, args []string) error {
	var err error
	if c.deleteAllServices {
		pretty.Progress("Deleting all services...", func() {
			err = c.e.ServiceDeleteAll(!c.keepData)
		})
		if err != nil {
			return err
		}
		fmt.Printf("%s all services are deleted\n", pretty.SuccessSign)
		return nil
	}

	exitWithError := false
	for _, arg := range args {
		// build function to avoid using arg inside progress
		fn := func(id string) func() {
			return func() {
				err = c.e.ServiceDelete(!c.keepData, id)
			}
		}(arg)
		pretty.Progress(fmt.Sprintf("Deleting service %q...", arg), fn)
		if err != nil {
			exitWithError = true
			fmt.Fprintf(os.Stderr, "%s can't delete service %q: %s\n", pretty.FailSign, arg, err)
		} else {
			fmt.Printf("%s service %q deleted\n", pretty.SuccessSign, arg)
		}
	}

	if exitWithError {
		return errors.New("there was a problem with deleting some services")
	}
	return nil
}
