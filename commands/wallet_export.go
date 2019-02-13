package commands

import (
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type walletExportCmd struct {
	baseCmd

	noPassphrase bool
	passphrase   string
	address      string

	e WalletExecutor
}

func newWalletExportCmd(e WalletExecutor) *walletExportCmd {
	c := &walletExportCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "export",
		Short:   "export an account",
		Long:    "export an existing account in order to backup it and import it in an other wallet",
		Example: "mesg-core wallet export",
		Args:    cobra.ExactArgs(1),
		RunE:    c.runE,
	})

	c.cmd.Flags().StringBoolP(&c.passphrase, "no-passphrase", "-n", c.noPassphrase, "Leave passphrase empty")
	c.cmd.Flags().StringVarP(&c.passphrase, "passphrase", "p", c.passphrase, "Passphrase to encrypt the account")
	return c
}

func (c *walletExportCmd) preRunE(cmd *cobra.Command, args []string) error {
	if !c.noPassphrase && c.passphrase == "" {
		return survey.AskOne(&survey.Password{
			Message: "Enther passphrase",
		}, &c.passphrase, survey.MinLength(1))
	}
	c.address = args[0]
	return nil
}

func (c *walletExportCmd) runE(cmd *cobra.Command, args []string) error {
	return nil
}
