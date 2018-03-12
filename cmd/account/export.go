package cmdAccount

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/mesg-foundation/application/account"

	"github.com/spf13/cobra"
)

// Export an account into a json file
var Export = &cobra.Command{
	Use:               "export ACCOUNT",
	Short:             "Export account details in order to be able to re-import it with the import command",
	Example:           "mesg-cli account export AccountX",
	Run:               exportHandler,
	DisableAutoGenTag: true,
}

func exportHandler(cmd *cobra.Command, args []string) {
	var account *account.Account
	if len(args) > 0 {
		account = cmdUtils.FindAccount(args[0])
	}
	if account == nil {
		account = cmdUtils.AskAccount("Choose the account you want to export")
	}
	path, err := account.Export()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Account exported: %s\n", cmdUtils.SuccessColor(path))
}

func init() {
	Export.Flags().StringP("name", "n", "", "Name of the account")
}
