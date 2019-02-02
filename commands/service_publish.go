package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type servicePublishCmd struct {
	baseCmd

	path string

	yes bool

	e ServiceExecutor
}

func newServicePublishCmd(e ServiceExecutor) *servicePublishCmd {
	c := &servicePublishCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "publish",
		Short:   "Publish a service on the MESG Marketplace",
		Example: `mesg-core service publish PATH_TO_SERVICE`,
		PreRunE: c.preRunE,
		RunE:    c.runE,
		Args:    cobra.MaximumNArgs(1),
		Hidden:  true, // TODO: Remove when this feature is finished
	})
	c.cmd.Flags().BoolVarP(&c.yes, "yes", "y", c.yes, `Automatic "yes" to all prompts and run non-interactively`)
	return c
}

func (c *servicePublishCmd) preRunE(cmd *cobra.Command, args []string) error {
	c.path = getFirstOrDefault(args)
	return nil
}

func (c *servicePublishCmd) runE(cmd *cobra.Command, args []string) error {
	if !c.yes {
		var confirmed bool
		if err := survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("Are you sure to publish from path? (%s)", c.path),
		}, &confirmed, nil); err != nil {
			return err
		}
		if !confirmed {
			return errConfirmationNeeded
		}
	}

	definition, err := c.e.ServicePublishDefinitionFile(c.path)
	if err != nil {
		return err
	}

	fmt.Println("https://gateway.ipfs.io/ipfs/" + definition)

	return nil
}
