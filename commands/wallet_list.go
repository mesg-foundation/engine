package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

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
	accounts, err := c.e.List()
	if err != nil {
		return err
	}
	if len(accounts) == 0 {
		fmt.Println("No account")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	fmt.Fprintf(w, "ID\tACCOUNT\t\n")
	for i, account := range accounts {
		fmt.Fprintf(w, "%d\t%s\t\n", i+1, account)
	}
	return w.Flush()
}
