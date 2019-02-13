package commands

import (
	"github.com/spf13/cobra"
)

type rootMarketplaceCmd struct {
	baseCmd
}

func newRootMarketplaceCmd(e MarketplaceExecutor) *rootMarketplaceCmd {
	c := &rootMarketplaceCmd{}
	c.cmd = newCommand(&cobra.Command{
		Use:   "marketplace",
		Short: "Interact with the Marketplace",
	})

	c.cmd.AddCommand(
		newMarketplacePublishCmd(e).cmd,
	)
	return c
}
