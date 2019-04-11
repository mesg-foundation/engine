package commands

import (
	"fmt"

	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type marketplaceCreateOfferCmd struct {
	baseMarketplaceCmd

	service  provider.MarketplaceService
	price    string // in MESG token
	duration string

	e Executor
}

func newMarketplaceCreateOfferCmd(e Executor) *marketplaceCreateOfferCmd {
	c := &marketplaceCreateOfferCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:   "create-offer",
		Short: "Create a new offer on a service on the MESG Marketplace",
		Example: `mesg-core marketplace create-offer SID
mesg-core marketplace create-offer SID --price 10 --duration 3600`,
		PreRunE: c.preRunE,
		RunE:    c.runE,
		Args:    cobra.MinimumNArgs(1),
	})
	c.setupFlags()
	c.cmd.Flags().StringVarP(&c.price, "price", "", c.price, "Price (in MESG token) of the offer to create")
	c.cmd.Flags().StringVarP(&c.duration, "duration", "", c.duration, "Duration (in second) of the offer to create")
	return c
}

func (c *marketplaceCreateOfferCmd) preRunE(cmd *cobra.Command, args []string) error {
	if err := c.askAccountAndPassphrase(); err != nil {
		return err
	}

	if c.price == "" {
		if err := askInput("Enter the price (in MESG Token) of the offer", &c.price); err != nil {
			return err
		}
	}
	if c.duration == "" {
		if err := askInput("Enter the duration (in second) of the offer", &c.duration); err != nil {
			return err
		}
	}

	var confirmed bool
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Are you sure to create offer for service %q with a price of %s MESG Tokens and a duration of %s seconds?", args[0], c.price, c.duration),
	}, &confirmed, nil); err != nil {
		return err
	}
	if !confirmed {
		return fmt.Errorf("cancelled")
	}

	return nil
}

func (c *marketplaceCreateOfferCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		tx                provider.Transaction
		signedTransaction string
		err               error
		sid, offerIndex   string
	)
	pretty.Progress("Creating offer on the marketplace...", func() {
		if tx, err = c.e.PrepareCreateServiceOffer(args[0], c.price, c.duration, c.account); err != nil {
			return
		}
		if signedTransaction, err = c.e.Sign(c.account, c.passphrase, tx); err != nil {
			return
		}
		sid, offerIndex, _, _, err = c.e.PublishCreateServiceOffer(signedTransaction)
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s Offer created with success with index %q\n", pretty.SuccessSign, offerIndex)
	fmt.Printf("%s See it on the marketplace: https://marketplace.mesg.com/services/%s#offers\n", pretty.SuccessSign, sid)

	return nil
}
