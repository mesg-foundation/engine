package cmdAccount

import (
	"errors"
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/application/account"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/spf13/cobra"
)

// Import an account from an exported file
var Import = &cobra.Command{
	Use:               "import FILE",
	Short:             "Import an account based on a file exported with the export command",
	Example:           "mesg-cli account import file.json",
	Args:              cobra.MinimumNArgs(1),
	Run:               importHandler,
	DisableAutoGenTag: true,
}

func importHandler(cmd *cobra.Command, args []string) {
	password := cmd.Flag("password").Value.String()
	if password == "" {
		survey.AskOne(&survey.Password{Message: "Type the current password ?"}, &password, nil)
	}
	newPassword := cmd.Flag("new-password").Value.String()
	if newPassword == "" {
		var passwordConfirmation string
		survey.AskOne(&survey.Password{Message: "Type the new password for your account ?"}, &newPassword, nil)
		survey.AskOne(&survey.Password{Message: "Repeat your password ?"}, &passwordConfirmation, nil)
		if newPassword != passwordConfirmation {
			panic(errors.New("Password confirmation invalid"))
		}
	}

	account, err := account.Import(args[0], password, newPassword)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Account imported: %s\n", aurora.Green(account).Bold())
}

func init() {
	Import.Flags().StringP("password", "", "", "Current password for the account you import")
	Import.Flags().StringP("new-password", "", "", "New password for the account you import")
}
