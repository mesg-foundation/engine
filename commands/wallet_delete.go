package commands

import (
	"fmt"

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type walletDeleteCmd struct {
	baseCmd

	noPassphrase bool
	passphrase   string
	address      string

	e WalletExecutor
}

func newWalletDeleteCmd(e WalletExecutor) *walletDeleteCmd {
	c := &walletDeleteCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "delete",
		Short:   "Delete an account",
		Long:    "Delete an account from the wallet",
		Example: "mesg-core wallet delete 0x0",
		Args:    cobra.ExactArgs(1),
		RunE:    c.runE,
	})

	c.cmd.Flags().BoolVarP(&c.passphrase, "no-passphrase", "-n", c.noPassphrase, "Leave passphrase empty")
	c.cmd.Flags().StringVarP(&c.passphrase, "passphrase", "p", c.passphrase, "Passphrase to encrypt the account")
	return c
}

func (c *walletDeleteCmd) preRunE(cmd *cobra.Command, args []string) error {
	if !c.noPassphrase && c.passphrase == "" {
		return survey.AskOne(&survey.Password{
			Message: "Enther passphrase",
		}, &c.passphrase, survey.MinLength(1))
	}
	c.address = args[0]
	return nil
}

func (c *walletDeleteCmd) runE(cmd *cobra.Command, args []string) error {
	if err := c.e.Delete(c.address, c.passphrase); err != nil {
		return err
	}
	fmt.Printf("%s Wallet deleted", pretty.SuccessSign)
	return nil
}
