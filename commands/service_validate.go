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
	c.path = getFirstOrDefault(args, "./")
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
