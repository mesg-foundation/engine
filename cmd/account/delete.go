package cmdAccount

import (
	"fmt"

	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// Delete a specific accounts
var Delete = &cobra.Command{
	Use:     "delete",
	Short:   "Delete an account",
	Example: "mesg-cli service delete 0x0000000000000000000000000000000000000000",
	Run:     deleteHandler,
}

func deleteHandler(cmd *cobra.Command, args []string) {
	var confirm bool
	var account = ""
	if len(args) > 0 {
		account = args[0]
	}
	if account == "" {
		// TODO add real list
		accounts := []string{"0x0000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000001"}
		survey.AskOne(&survey.Select{
			Message: "Choose the account you want to delete",
			Default: accounts[0],
			Options: accounts,
		}, &account, nil)
	}
	survey.AskOne(&survey.Confirm{Message: "Are you sure ? You can always re-import an account with your private seed"}, &confirm, nil)
	if confirm {
		// TODO add real deletion
		fmt.Println("account deleted", account)
	}
}
