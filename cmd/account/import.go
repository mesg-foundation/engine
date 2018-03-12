package cmdAccount

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/mesg-foundation/application/account"

	"github.com/spf13/cobra"
)

// Import an account from an exported file
var Import = &cobra.Command{
	Use:               "import FILE",
	Short:             "Import an account based on a file exported with the export command",
	Example:           "mesg-cli account import file.json",
	Run:               importHandler,
	DisableAutoGenTag: true,
}

func importHandler(cmd *cobra.Command, args []string) {
	account, err := account.Import(args[0], cmd.Flag("name").Value.String())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Account imported: %s\n", cmdUtils.SuccessColor(account))
}

func init() {
	Import.Flags().StringP("name", "n", "", "Name of the account")
}
