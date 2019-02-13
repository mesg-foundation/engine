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
	// definition, err := c.e.PublishDefinitionFile(c.path)
	tx, err := c.e.CreateService("test1", "0xf3C21FD07B1D4c40d3cE6EfaC81a3E49f6c04592")
	if err != nil {
		return err
	}
	fmt.Println("tx", tx)

	// fmt.Println("https://gateway.ipfs.io/ipfs/" + definition)
	fmt.Println("published")

	return nil
}
