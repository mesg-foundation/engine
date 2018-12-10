package commands

import (
	"errors"
	"fmt"
	"sync"

	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/mesg-foundation/core/x/xerrors"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type serviceDeployCmd struct {
	baseCmd

	path string

	force bool

	e ServiceExecutor
}

func newServiceDeployCmd(e ServiceExecutor) *serviceDeployCmd {
	c := &serviceDeployCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:   "deploy",
		Short: "Deploy a service",
		Long: `Deploy a service.

To get more information, see the [deploy page from the documentation](https://docs.mesg.com/guide/service/deploy-a-service.html)`,
		Example: `mesg-core service deploy PATH_TO_SERVICE`,
		PreRunE: c.preRunE,
		RunE:    c.runE,
		Args:    cobra.MaximumNArgs(1),
	})
	c.cmd.Flags().BoolVarP(&c.force, "force", "f", c.force, "Force deploying overwrites existing service with the same sid")
	return c
}

func (c *serviceDeployCmd) preRunE(cmd *cobra.Command, args []string) error {
	c.path = getFirstOrDefault(args)
	return nil
}

func (c *serviceDeployCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		statuses = make(chan provider.DeployStatus)
		wg       sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		printDeployStatuses(statuses)
	}()

	var confirmation *bool
	if c.force {
		confirmation = &c.force
	}

	id, validationError, err := c.e.ServiceDeploy(c.path, confirmation, func(sid string) (bool, error) {
		pretty.DestroySpinner()
		var confirm bool
		if err := survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("A service with the same sid %q already exists. Do you want to replace it?", sid),
		}, &confirm, nil); err != nil {
			return false, errors.New("confirmation rejected")
		}
		pretty.UseSpinner("Sending confirmation status")
		return confirm, nil
	}, statuses)
	wg.Wait()

	pretty.DestroySpinner()
	if err != nil {
		return err
	}
	if validationError != nil {
		return xerrors.Errors{
			validationError,
			errors.New("to get more information, run: mesg-core service validate"),
		}
	}
	fmt.Printf("%s Service deployed with hash: %v\n", pretty.SuccessSign, pretty.Success(id))
	fmt.Printf("To start it, run the command:\n\tmesg-core service start %s\n", id)
	return nil
}

func printDeployStatuses(statuses chan provider.DeployStatus) {
	for status := range statuses {
		switch status.Type {
		case provider.Running:
			pretty.UseSpinner(status.Message)
		default:
			var sign string
			switch status.Type {
			case provider.DonePositive:
				sign = pretty.SuccessSign
			case provider.DoneNegative:
				sign = pretty.FailSign
			}
			pretty.DestroySpinner()
			fmt.Printf("%s %s\n", sign, status.Message)
		}
	}
}
