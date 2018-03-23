package cmdAccount

import (
	"errors"
	"fmt"

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
	account := &account.Account{
		Password: cmd.Flag("password").Value.String(),
		Name:     cmd.Flag("name").Value.String(),
	}
	if err := checkPassword(account); err != nil {
		panic(err)
	}

	if err := checkName(account); err != nil {
		panic(err)
	}

	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Generating secure key..."})
	if err := account.Generate(); err != nil {
		panic(err)
	}
	s.Stop()

	displaySummary(account)
}

func checkPassword(account *account.Account) (err error) {
	if account.Password != "" {
		return
	}
	var passwordConfirmation string
	survey.AskOne(&survey.Password{Message: "Please set a password ?"}, &account.Password, nil)
	survey.AskOne(&survey.Password{Message: "Repeat your password ?"}, &passwordConfirmation, nil)
	if account.Password != passwordConfirmation {
		err = errors.New("Password confirmation invalid")
		return
	}
	return
}

func checkName(account *account.Account) (err error) {
	if account.Name != "" {
		return
	}
	survey.AskOne(&survey.Input{Message: "Choose a name for this account"}, &account.Name, nil)
	return
}

func displaySummary(account *account.Account) {
	fmt.Println("Here is all the details of your account:")
	fmt.Println()
	fmt.Printf("Account name: %s\n", aurora.Green(account.Name).Bold())
	fmt.Printf("Account address: %s\n", aurora.Green(account.Address.String()).Bold())
	fmt.Println()
	fmt.Printf("%s", aurora.Brown("Please make sure that you save those information").Bold())
	fmt.Println()
}

func init() {
	Create.Flags().StringP("name", "n", "", "Name of the account")
	Create.Flags().StringP("password", "p", "", "Password for the account")
}
