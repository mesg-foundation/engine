package commands

import (
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type walletSignCmd struct {
	baseCmd

	noPassphrase bool
	passphrase   string
	address      string

	e WalletExecutor
}

func newWalletSignCmd(e WalletExecutor) *walletSignCmd {
	c := &walletSignCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:   "list",
		Short: "list wallets",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runE,
	})

	c.cmd.Flags().BoolVarP(&c.passphrase, "no-passphrase", "-n", c.noPassphrase, "Leave passphrase empty")
	c.cmd.Flags().StringVarP(&c.passphrase, "passphrase", "p", c.passphrase, "Passphrase to encrypt the account")
	return c
}

func (c *walletSignCmd) preRunE(cmd *cobra.Command, args []string) error {
	if !c.noPassphrase && c.passphrase == "" {
		return survey.AskOne(&survey.Password{
			Message: "Enther passphrase",
		}, &c.passphrase, survey.MinLength(1))
	}
	c.address = args[0]
	return nil
}

func (c *walletSignCmd) runE(cmd *cobra.Command, args []string) error {
	return nil
}
