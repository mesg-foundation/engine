package account

import (
	"fmt"

	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/spf13/cobra"
)

// Detail one account
var Detail = &cobra.Command{
	Use:   "detail",
	Short: "Show detailed information of an account",
	Long:  `To show the balance, previous transactions and some other information of an account`,
	Example: `mesg-core account detail
mesg-core account detail --account ACCOUNT`,
	Run:               detailHandler,
	DisableAutoGenTag: true,
}

func detailHandler(cmd *cobra.Command, args []string) {
	acc := utils.AccountFromFlagOrAsk(cmd, "Choose account:")

	// TODO: can facto with a displaySummary function like in create.go
	fmt.Println(acc.Address.String())
}

func init() {
	utils.Accountable(Detail)
}
