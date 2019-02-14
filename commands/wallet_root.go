package commands

import (
	"github.com/spf13/cobra"
)

type rootWalletCmd struct {
	baseCmd
}

func newRootWalletCmd(e WalletExecutor) *rootWalletCmd {
	c := &rootWalletCmd{}
	c.cmd = newCommand(&cobra.Command{
		Use:   "wallet",
		Short: "Manage wallets",
	})

	c.cmd.AddCommand(
		newWalletListCmd(e).cmd,
		newWalletCreateCmd(e).cmd,
		newWalletDeleteCmd(e).cmd,
		newWalletExportCmd(e).cmd,
		newWalletImportCmd(e).cmd,
	)
	return c
}

// baseWalletCmd is basic command for wallet subcommands
// that require passphrase.
type baseWalletCmd struct {
	baseCmd
	noPassphrase bool
	passphrase   string
}

func (c *baseWalletCmd) setupFlags() {
	c.cmd.Flags().BoolVarP(&c.noPassphrase, "no-passphrase", "n", c.noPassphrase, "Leave passphrase empty")
	c.cmd.Flags().StringVarP(&c.passphrase, "passphrase", "p", c.passphrase, "Passphrase to encrypt the account")
}
