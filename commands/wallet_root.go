package commands

import (
	"errors"

	"github.com/spf13/cobra"
)

var errInvalidAddress = errors.New("invalid wallet address")

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
