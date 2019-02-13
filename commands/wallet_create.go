package commands

import (
	"fmt"

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type walletCreateCmd struct {
	baseCmd

	noPassphrase bool
	passphrase   string

	e WalletExecutor
}

func newWalletCreateCmd(e WalletExecutor) *walletCreateCmd {
	c := &walletCreateCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "create",
		Short:   "Create a new account",
		Long:    "Create a new account with a passphrase",
		Example: "mesg-core wallet create --no-passphrase",
		PreRunE: c.preRunE,
		RunE:    c.runE,
	})

	c.cmd.Flags().BoolVarP(&c.noPassphrase, "no-passphrase", "-n", c.noPassphrase, "Leave passphrase empty")
	c.cmd.Flags().StringVarP(&c.passphrase, "passphrase", "p", c.passphrase, "Passphrase to encrypt the account")
	return c
}

func (c *walletCreateCmd) preRunE(cmd *cobra.Command, args []string) error {
	if !c.noPassphrase && c.passphrase == "" {
		return survey.AskOne(&survey.Password{
			Message: "Enther passphrase",
		}, &c.passphrase, survey.MinLength(1))
	}
	return nil
}

func (c *walletCreateCmd) runE(cmd *cobra.Command, args []string) error {
	address, err := c.e.Create(c.passphrase)
	if err != nil {
		return err
	}

	fmt.Printf("NOTE: remember to save passphrase\n\n")
	fmt.Printf("%s Wallet created with address %s", pretty.SuccessSign, pretty.Success(address.String()))
	return nil
}
