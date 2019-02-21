package commands

import (
	"fmt"

	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type marketplaceTransferOwnershipCmd struct {
	baseMarketplaceCmd

	service *provider.MarketplaceService

	e Executor
}

func newMarketplaceTransferOwnershipCmd(e Executor) *marketplaceTransferOwnershipCmd {
	c := &marketplaceTransferOwnershipCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "transfer-ownership",
		Short:   "Transfer the ownership of the service to another account",
		Example: `mesg-core marketplace transfer-ownership SID NEW_OWNER`,
		PreRunE: c.preRunE,
		RunE:    c.runE,
		Args:    cobra.MinimumNArgs(2),
	})
	c.setupFlags()
	return c
}

func (c *marketplaceTransferOwnershipCmd) preRunE(cmd *cobra.Command, args []string) error {
	var (
		err error
	)

	if err := c.askAccountAndPassphrase(); err != nil {
		return nil
	}

	pretty.Progress("Getting service data...", func() {
		c.service, err = c.e.GetService(args[0])
	})
	if err != nil {
		return err
	}
	if c.service.Owner != c.account {
		return fmt.Errorf("the service's owner %q is different than the specified account", c.service.Owner)
	}

	var confirmed bool
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Are you sure to transfer ownership of service %q to %q?", args[0], args[1]),
	}, &confirmed, nil); err != nil {
		return err
	}
	if !confirmed {
		return fmt.Errorf("cancelled")
	}

	return nil
}

func (c *marketplaceTransferOwnershipCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		tx  *provider.Transaction
		err error
	)
	pretty.Progress("Transfering ownership...", func() {
		tx, err = c.e.TransferServiceOwnership(args[0], args[1], c.account)
		if err != nil {
			return
		}
		_, err = c.signAndSendTransaction(c.e, tx)
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s Ownership transferred with success\n", pretty.SuccessSign)
	fmt.Printf("%s See it on the marketplace: https://marketplace.mesg.com/services/%s\n", pretty.SuccessSign, c.sha3(args[0]))

	return nil
}
