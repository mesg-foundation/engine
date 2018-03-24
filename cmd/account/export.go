package cmdAccount

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// Export an account into a json file
var Export = &cobra.Command{
	Use:               "export",
	Short:             "Export account details in order to be able to re-import it with the import command",
	Example:           "mesg-cli account export --account AccountX",
	Run:               exportHandler,
	DisableAutoGenTag: true,
}

func exportHandler(cmd *cobra.Command, args []string) {
	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Choose the account you want to export")
	fmt.Println("account exported : " + account.Address.String())
	// TODO
	// path, err := account.Export()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Account exported: %s\n", aurora.Green(path).Bold())
}

func init() {
	cmdUtils.Accountable(Export)
}
