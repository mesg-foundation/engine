package cmdAccount

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// Detail one account
var Detail = &cobra.Command{
	Use:               "detail",
	Short:             "Show detailed information of an account",
	Long:              `To show the balance, previous transactions and some other information of an account`,
	Example:           "mesg-cli detail",
	Run:               detailHandler,
	DisableAutoGenTag: true,
}

func detailHandler(cmd *cobra.Command, args []string) {
	acc := cmdUtils.AccountFromFlagOrAsk(cmd, "Choose account:")

	// TODO: can facto with a displaySummary function like in create.go
	fmt.Println(acc.Address.String())
}

func init() {
	cmdUtils.Accountable(Detail)
}
