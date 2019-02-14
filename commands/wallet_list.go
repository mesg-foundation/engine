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
	if len(addresses) == 0 {
		fmt.Println("No account")
		return nil
	}
	for _, address := range addresses {
		fmt.Printf("Address: %s\n", address)
	}
	return nil
}
