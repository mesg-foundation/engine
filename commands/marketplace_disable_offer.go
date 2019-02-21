package commands

import (
	"fmt"
	"strconv"

	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type marketplaceDisableOfferCmd struct {
	baseMarketplaceCmd

	service *provider.MarketplaceService

	e Executor
}

func newMarketplaceDisableOfferCmd(e Executor) *marketplaceDisableOfferCmd {
	c := &marketplaceDisableOfferCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "disable-offer",
		Short:   "Disable an offer on a service on the MESG Marketplace",
		Example: `mesg-core marketplace disable-offer SID OFFER_INDEX`,
		PreRunE: c.preRunE,
		RunE:    c.runE,
		Args:    cobra.MinimumNArgs(2),
	})
	c.setupFlags()
	return c
}

func (c *marketplaceDisableOfferCmd) preRunE(cmd *cobra.Command, args []string) error {
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
	offerIndex, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}
	if offerIndex >= len(c.service.Offers) {
		return fmt.Errorf("offer index is out of range")
	}

	var confirmed bool
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Are you sure to disable offer with index %q on service %q?", args[1], args[0]),
	}, &confirmed, nil); err != nil {
		return err
	}
	if !confirmed {
		return fmt.Errorf("cancelled")
	}

	return nil
}

func (c *marketplaceDisableOfferCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		tx  *provider.Transaction
		err error
	)
	pretty.Progress("Disabling offer on the marketplace...", func() {
		tx, err = c.e.DisableServiceOffer(args[0], args[1], c.account)
		if err != nil {
			return
		}
		_, err = c.signAndSendTransaction(c.e, tx)
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s Offer disabled with success\n", pretty.SuccessSign)
	fmt.Printf("%s See it on the marketplace: https://marketplace.mesg.com/services/%s\n", pretty.SuccessSign, c.sha3(args[0]))

	return nil
}
