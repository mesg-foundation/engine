package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type marketplacePublishCmd struct {
	baseCmd

	path string

	e MarketplaceExecutor
}

func newMarketplacePublishCmd(e MarketplaceExecutor) *marketplacePublishCmd {
	c := &marketplacePublishCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "publish",
		Short:   "Publish a service on the MESG Marketplace",
		Example: `mesg-core marketplace publish PATH_TO_SERVICE`,
		PreRunE: c.preRunE,
		RunE:    c.runE,
		Args:    cobra.MaximumNArgs(1),
	})
	return c
}

func (c *marketplacePublishCmd) preRunE(cmd *cobra.Command, args []string) error {
	c.path = getFirstOrDefault(args)
	return nil
}

func (c *marketplacePublishCmd) runE(cmd *cobra.Command, args []string) error {
	definition, err := c.e.MarketplacePublishDefinitionFile(c.path)
	if err != nil {
		return err
	}

	fmt.Println("https://gateway.ipfs.io/ipfs/" + definition)

	return nil
}
