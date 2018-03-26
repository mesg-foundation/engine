package cmdAccount

import (
	"errors"

	"github.com/mesg-foundation/application/account"

	"github.com/mesg-foundation/application/cmd/utils"
	survey "gopkg.in/AlecAivazis/survey.v1"

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
	path := cmd.Flag("path").Value.String()
	acc := cmdUtils.AccountFromFlagOrAsk(cmd, "Choose the account you want to export")
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

	err := account.Export(acc, password, newPassword, path)
	if err != nil {
		panic(err)
	}
}

func init() {
	cmdUtils.Accountable(Export)
	Export.Flags().StringP("password", "", "", "Current password for the account you export")
	Export.Flags().StringP("new-password", "", "", "New password for the account you export")
	Export.Flags().StringP("path", "p", "./export", "Path of the file where your account will be exported")
}
