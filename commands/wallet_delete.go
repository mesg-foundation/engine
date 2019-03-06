package commands

import (
	"fmt"

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
)

type walletDeleteCmd struct {
	baseWalletCmd

	account string

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
	if err := c.askPassphrase(); err != nil {
		return err
	}
	// TODO: if no account provided, the cli should ask to select one.
	c.account = args[0]
	return nil
}

func (c *walletDeleteCmd) runE(cmd *cobra.Command, args []string) error {
	account, err := c.e.Delete(c.account, c.passphrase)
	if err != nil {
		return err
	}
	fmt.Printf("%s Account %q deleted\n", pretty.SuccessSign, account)
	return nil
}
