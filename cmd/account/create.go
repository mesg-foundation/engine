package cmdAccount

import (
	"errors"
	"fmt"
	"time"

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

	if err := generateAccount(account); err != nil {
		panic(err)
	}

	displaySummary(account)
}

func checkPassword(account *account.Account) error {
	if account.Password != "" {
		return nil
	}
	var passwordConfirmation string
	survey.AskOne(&survey.Password{Message: "Please set a password ?"}, &account.Password, nil)
	survey.AskOne(&survey.Password{Message: "Repeat your password ?"}, &passwordConfirmation, nil)
	if account.Password != passwordConfirmation {
		return errors.New("Password confirmation invalid")
	}
	return nil
}

func checkName(account *account.Account) error {
	if account.Name != "" {
		return nil
	}
	survey.AskOne(&survey.Input{Message: "Choose a name for this account"}, &account.Name, nil)
	return nil
}

func generateAccount(account *account.Account) error {
	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Generating secure key..."})
	time.Sleep(time.Second)
	s.Stop()

	return account.Generate()
}

func displaySummary(account *account.Account) {
	fmt.Println("Here is all the details of your account:")
	fmt.Println()
	fmt.Printf("Account name: %s\n", aurora.Green(account.Name).Bold())
	fmt.Printf("Account address: %s\n", aurora.Green(account.Address).Bold())
	fmt.Printf("Seed: %s\n", aurora.Green(account.Seed).Bold())
	fmt.Println()
	fmt.Printf("%s", aurora.Brown("Please make sure that you save those information and especially the following seed that might be needed to regenerate this address").Bold())
	fmt.Println()
}

func init() {
	Create.Flags().StringP("name", "n", "", "Name of the account")
	Create.Flags().StringP("password", "p", "", "Password for the account")
}
