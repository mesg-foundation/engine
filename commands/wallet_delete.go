package commands

import (
	"fmt"

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
)

type walletDeleteCmd struct {
	baseWalletCmd

	address string

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
		PreRunE: c.preRunE,
		RunE:    c.runE,
	})
	c.setupFlags()
	return c
}

func (c *walletDeleteCmd) preRunE(cmd *cobra.Command, args []string) error {
	if !c.noPassphrase && c.passphrase == "" {
		if err := askPass("Enter passphrase", &c.passphrase); err != nil {
			return err
		}
	}
	c.address = args[0]
	return nil
}

func (c *walletDeleteCmd) runE(cmd *cobra.Command, args []string) error {
	address, err := c.e.Delete(c.address, c.passphrase)
	if err != nil {
		return err
	}
	fmt.Printf("%s Wallet %s deleted\n", pretty.SuccessSign, address)
	return nil
}
