package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type servicePublishCmd struct {
	baseCmd

	path string

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
	return c
}

func (c *servicePublishCmd) preRunE(cmd *cobra.Command, args []string) error {
	c.path = getFirstOrCurrentPath(args)
	return nil
}

func (c *servicePublishCmd) runE(cmd *cobra.Command, args []string) error {
	definition, err := c.e.ServicePublishDefinitionFile(c.path)
	if err != nil {
		return err
	}

	fmt.Println("https://gateway.ipfs.io/ipfs/" + definition)

	return nil
}
