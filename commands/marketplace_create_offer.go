package commands

import (
	"fmt"
	"strings"

	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type marketplaceCreateOfferCmd struct {
	baseMarketplaceCmd

	service  *provider.MarketplaceService
	price    string // in MESG token
	duration string

	e Executor
}

func newMarketplaceCreateOfferCmd(e Executor) *marketplaceCreateOfferCmd {
	c := &marketplaceCreateOfferCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "create-offer",
		Short:   "Create a new offer on a service on the MESG Marketplace",
		Example: `mesg-core marketplace create-offer SID`,
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
	var (
		err error
	)

	if err := c.askAccountAndPassphrase(); err != nil {
		return err
	}

	if c.price == "" {
		if err := askInput("Enter the price (in MESG Token) of the offer", &c.price); err != nil {
			return err
		}
	}
	// if c.price < 0 {
	// 	return fmt.Errorf("Price cannot be negative")
	// }
	if c.duration == "" {
		if err := askInput("Enter the duration (in second) of the offer", &c.duration); err != nil {
			return err
		}
	}

	pretty.Progress("Getting service data...", func() {
		c.service, err = c.e.GetService(args[0])
	})
	if err != nil {
		return err
	}
	if !strings.EqualFold(c.service.Owner, c.account) {
		return fmt.Errorf("the service's owner %q is different than the specified account", c.service.Owner)
	}

	var confirmed bool
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Are you sure to create offer on service %q with price %q and duration %q?", args[0], c.price, c.duration),
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
		tx  *provider.Transaction
		err error
	)
	pretty.Progress("Creating offer on the marketplace...", func() {
		tx, err = c.e.CreateServiceOffer(args[0], c.price, c.duration, c.account)
		if err != nil {
			return
		}
		_, err = c.signAndSendTransaction(c.e, tx)
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s Offer created with success\n", pretty.SuccessSign)
	fmt.Printf("%s See it on the marketplace: https://marketplace.mesg.com/services/%s\n", pretty.SuccessSign, c.sha3(args[0]))

	return nil
}
