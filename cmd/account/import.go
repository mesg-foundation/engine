package cmdAccount

import (
	"errors"
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/application/account"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/spf13/cobra"
)

// Import an account from an exported file
var Import = &cobra.Command{
	Use:   "import ./PATH_TO_BACKUP_FILE",
	Short: "Import an account from a backup file",
	Long:  `This method imports a previously exported backup file of your account created with the [export method](mesg-cli_account_export.md).`,
	Example: `mesg-cli account import ./PATH_TO_BACKUP_FILE
mesg-cli account import ./PATH_TO_BACKUP_FILE --password PASSWORD --new-password PASSWORD`,
	Args:              cobra.MinimumNArgs(1),
	Run:               importHandler,
	DisableAutoGenTag: true,
}

func importHandler(cmd *cobra.Command, args []string) {
	password := cmd.Flag("password").Value.String()
	if password == "" && survey.AskOne(&survey.Password{Message: "Type current password:"}, &password, nil) != nil {
		os.Exit(0)
	}
	newPassword := cmd.Flag("new-password").Value.String()
	if newPassword == "" {
		var passwordConfirmation string
		if survey.AskOne(&survey.Password{Message: "Type new password:"}, &newPassword, nil) != nil {
			os.Exit(0)
		}
		if survey.AskOne(&survey.Password{Message: "Repeat password:"}, &passwordConfirmation, nil) != nil {
			os.Exit(0)
		}
		if newPassword != passwordConfirmation {
			panic(errors.New("Passwords are different"))
		}
	}

	account, err := account.Import(args[0], password, newPassword)
	if err != nil {
		panic(err)
	}

	// TODO: can facto with a displaySummary function like in create.go
	fmt.Printf("Account imported: %s\n", aurora.Green(account).Bold())
}

func init() {
	Import.Flags().StringP("password", "", "", "Current password of the account to import")
	Import.Flags().StringP("new-password", "", "", "New password of the imported account")
}
