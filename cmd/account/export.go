package cmdAccount

import (
	"errors"
	"os"

	"github.com/mesg-foundation/application/account"
	"github.com/mesg-foundation/application/cmd/utils"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// Export an account into a json file
var Export = &cobra.Command{
	Use:   "export",
	Short: "Export an account",
	Long: `This method creates a file containing the information about your account.
The private key of your account is encrypted with the password you choose.

**Warning:** This method does **NOT** export your password. You have to manage your password yourself.

You can import the backup file on any other MESG Application with the [import method](mesg-cli_account_import.md).`,
	Example: `mesg-cli account export
mesg-cli account export --account ACCOUNT --password PASSWORD --new-password PASSWORD --path ./PATH_TO_BACKUP_FILE`,
	Run:               exportHandler,
	DisableAutoGenTag: true,
}

func exportHandler(cmd *cobra.Command, args []string) {
	path := cmd.Flag("path").Value.String()
	acc := cmdUtils.AccountFromFlagOrAsk(cmd, "Choose the account to export:")
	password := cmd.Flag("password").Value.String()
	if password == "" && survey.AskOne(&survey.Password{Message: "Type current password:"}, &password, nil) != nil {
		os.Exit(0)
	}
	newPassword := cmd.Flag("new-password").Value.String()
	if newPassword == "" {
		var passwordConfirmation string
		if survey.AskOne(&survey.Password{Message: "Type new password for exportation:"}, &newPassword, nil) != nil {
			os.Exit(0)
		}
		if survey.AskOne(&survey.Password{Message: "Repeat password for exportation:"}, &passwordConfirmation, nil) != nil {
			os.Exit(0)
		}
		if newPassword != passwordConfirmation {
			panic(errors.New("Passwords are different"))
		}
	}

	err := account.Export(acc, password, newPassword, path)
	if err != nil {
		panic(err)
	}

	// TODO: show confirmation with path
}

func init() {
	cmdUtils.Accountable(Export)
	Export.Flags().StringP("password", "", "", "Current password of the account to export")
	Export.Flags().StringP("new-password", "", "", "New password of the exported account")
	Export.Flags().StringP("path", "p", "./export", "Path of the file your account will be exported in")
}
