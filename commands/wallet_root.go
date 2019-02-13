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
		Use:   "service",
		Short: "Manage wallets",
	})

	c.cmd.AddCommand(
		newWalletListCmd(e).cmd,
		newWalletCreateCmd(e).cmd,
		newWalletDeleteCmd(e).cmd,
		newWalletExportCmd(e).cmd,
		newWalletImportCmd(e).cmd,
		newWalletSignCmd(e).cmd,
	)
	return c
}
