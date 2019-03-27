package commands

import (
	"fmt"
	"strconv"

	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type marketplacePurchaseCmd struct {
	baseMarketplaceCmd

	offerIndex string

	e Executor
}

func newMarketplacePurchaseCmd(e Executor) *marketplacePurchaseCmd {
	c := &marketplacePurchaseCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:   "purchase",
		Short: "Purchase a service on the MESG Marketplace",
		Example: `mesg-core marketplace purchase SID
mesg-core marketplace purchase SID --offer-index 1`,
		PreRunE: c.preRunE,
		RunE:    c.runE,
		Args:    cobra.ExactArgs(1),
	})
	c.setupFlags()
	c.cmd.Flags().StringVarP(&c.offerIndex, "offer-index", "o", c.offerIndex, "Offer index to purchase")
	return c
}

func (c *marketplacePurchaseCmd) preRunE(cmd *cobra.Command, args []string) error {
	var (
		service provider.MarketplaceService
		err     error
	)

	if err := c.askAccountAndPassphrase(); err != nil {
		return err
	}

	if c.offerIndex == "" {
		// TODO: should display the list of offers and ask to select one
		if err := askInput("Enter the offer index to purchase", &c.offerIndex); err != nil {
			return err
		}
	}

	pretty.Progress("Getting service data...", func() {
		service, err = c.e.GetService(args[0])
	})
	if err != nil {
		return err
	}

	offerIndex, err := strconv.Atoi(c.offerIndex)
	if err != nil {
		return err
	}
	if offerIndex >= len(service.Offers) {
		return fmt.Errorf("offer index %d doesn't exist", offerIndex)
	}
	offer := service.Offers[offerIndex]
	if !offer.Active {
		return fmt.Errorf("cannot purchase this offer because it is not active")
	}

	var confirmed bool
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Are you sure to purchase service %q for price %q and duration %q?", args[0], offer.Price, offer.Duration),
	}, &confirmed, nil); err != nil {
		return err
	}
	if !confirmed {
		return fmt.Errorf("cancelled")
	}

	return nil
}

func (c *marketplacePurchaseCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		transactions []provider.Transaction
		err          error
	)
	pretty.Progress("Purchasing offer on the marketplace...", func() {
		transactions, err = c.e.Purchase(args[0], c.offerIndex, c.account)
		if err != nil {
			return
		}
		for _, tx := range transactions {
			_, err = c.signAndSendTransaction(c.e, tx)
			if err != nil {
				return
			}
		}
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s Offer purchased with success\n", pretty.SuccessSign)

	return nil
}
