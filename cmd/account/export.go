package cmdAccount

import (
	"errors"
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/mesg-foundation/application/account"

	"github.com/spf13/cobra"
)

// Export an account into a json file
var Export = &cobra.Command{
	Use:               "export",
	Short:             "Export accounts in order to be able to re-import it with the import command",
	Example:           "mesg-cli account export",
	Run:               exportHandler,
	DisableAutoGenTag: true,
}

func exportHandler(cmd *cobra.Command, args []string) {
	var account *account.Account
	if name := cmd.Flag("name").Value.String(); name != "" {
		account = cmdUtils.FindAccount(name)
		if account == nil {
			panic(errors.New("Account '" + name + "' does not exist"))
		}
	}
	if account == nil {
		account = cmdUtils.AskAccount("Choose the account you want to export")
	}
	path, err := account.Export()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Account exported: %s\n", aurora.Green(path).Bold())
}

func init() {
	Export.Flags().StringP("name", "n", "", "Name of the account")
}
