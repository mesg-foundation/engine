package cmdAccount

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/application/account"
	"github.com/mesg-foundation/application/cmd/utils"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/spf13/cobra"
)

// Create run the create command for an account
var Create = &cobra.Command{
	Use:               "create",
	Short:             "Create a new account",
	Long:              "Create a new account composed of a name and a generated address",
	Example:           "mesg-cli account create",
	Run:               createHandler,
	DisableAutoGenTag: true,
}

func createHandler(cmd *cobra.Command, args []string) {
	password := cmd.Flag("password").Value.String()
	if password == "" {
		var passwordConfirmation string
		survey.AskOne(&survey.Password{Message: "Please set a password ?"}, &password, nil)
		survey.AskOne(&survey.Password{Message: "Repeat your password ?"}, &passwordConfirmation, nil)
		if password != passwordConfirmation {
			panic(errors.New("Password confirmation invalid"))
		}
	}

	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Generating secure key..."})
	acc, err := account.Generate(password)
	if err != nil {
		panic(err)
	}
	s.Stop()

	displaySummary(acc)
}

func displaySummary(acc accounts.Account) {
	fmt.Println("Here is all the details of your account:")
	fmt.Println()
	fmt.Printf("Account address: %s\n", aurora.Green(acc.Address.String()).Bold())
	fmt.Println()
	fmt.Printf("%s", aurora.Brown("Please make sure that you save those information").Bold())
	fmt.Println()
}

func init() {
	Create.Flags().StringP("password", "p", "", "Password for the account")
}
