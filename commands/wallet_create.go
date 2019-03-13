package commands

import (
	"fmt"

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
)

type walletCreateCmd struct {
	baseWalletCmd

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
	c.setupFlags()
	return c
}

func (c *walletCreateCmd) preRunE(cmd *cobra.Command, args []string) error {
	return c.askPassphrase()
}

func (c *walletCreateCmd) runE(cmd *cobra.Command, args []string) error {
	address, err := c.e.Create(c.passphrase)
	if err != nil {
		return err
	}

	fmt.Printf("NOTE: remember to save passphrase\n\n")
	fmt.Printf("%s Account created with address %s\n", pretty.SuccessSign, pretty.Success(address))
	return nil
}
