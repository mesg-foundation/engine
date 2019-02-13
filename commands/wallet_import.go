package commands

import (
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type walletImportCmd struct {
	baseCmd

	noPassphrase bool
	passphrase   string
	address      string

	e WalletExecutor
}

func newWalletImportCmd(e WalletExecutor) *walletImportCmd {
	c := &walletImportCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "import ",
		Short:   "import an account",
		Long:    "import an account to wallet",
		Example: "mesg-core wallet import",
		Args:    cobra.ExactArgs(1),
		RunE:    c.runE,
	})

	c.cmd.Flags().StringBoolP(&c.passphrase, "no-passphrase", "-n", c.noPassphrase, "Leave passphrase empty")
	c.cmd.Flags().StringVarP(&c.passphrase, "passphrase", "p", c.passphrase, "Passphrase to encrypt the account")
	return c
}

func (c *walletImportCmd) preRunE(cmd *cobra.Command, args []string) error {
	if !c.noPassphrase && c.passphrase == "" {
		return survey.AskOne(&survey.Password{
			Message: "Enther passphrase",
		}, &c.passphrase, survey.MinLength(1))
	}
	c.address = args[0]
	return nil
}

func (c *walletImportCmd) runE(cmd *cobra.Command, args []string) error {
	return nil
}
