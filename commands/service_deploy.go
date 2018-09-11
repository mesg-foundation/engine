package commands

import (
	"errors"
	"fmt"

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
)

type serviceDeployCmd struct {
	baseCmd

	path string

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
	return c
}

func (c *serviceDeployCmd) preRunE(cmd *cobra.Command, args []string) error {
	c.path = getFirstOrDefault(args, "./")
	return nil
}

func (c *serviceDeployCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		id    string
		valid bool
		err   error
	)
	pretty.Progress("Deploying the service...", func() {
		id, valid, err = c.e.ServiceDeploy(c.path)
	})
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("Service is invalid. To get more information, run: mesg-core service validate")
	}
	fmt.Println("Service deployed with ID:", pretty.Success(id))
	fmt.Printf("To start it, run the command: mesg-core service start %s\n", id)
	return nil
}
