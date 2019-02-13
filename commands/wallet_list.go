package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type walletListCmd struct {
	baseCmd

	e WalletExecutor
}

func newWalletListCmd(e WalletExecutor) *walletListCmd {
	c := &walletListCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "list",
		Short:   "list accounts",
		Long:    "list the addresses of existing accounts",
		Example: "mesg-core wallet list",
		RunE:    c.runE,
	})

	return c
}

func (c *walletListCmd) runE(cmd *cobra.Command, args []string) error {
	addresses, err := c.e.List()
	if err != nil {
		return err
	}
	for _, address := range addresses {
		fmt.Println(address)
	}
	return nil
}
