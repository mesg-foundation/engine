package cmdAccount

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/mesg-foundation/application/cmd/utils"
	"github.com/mesg-foundation/application/types"
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
	account := &types.Account{
		Password: cmd.Flag("password").Value.String(),
		Name:     cmd.Flag("name").Value.String(),
	}
	if err := checkPassword(account); err != nil {
		fmt.Println(err)
	}

	if err := checkName(account); err != nil {
		fmt.Println(err)
	}

	if err := generateAccount(account); err != nil {
		fmt.Println(err)
	}

	displayResume(account)
}

func checkPassword(account *types.Account) error {
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

func checkName(account *types.Account) error {
	if account.Name != "" {
		return nil
	}
	survey.AskOne(&survey.Input{Message: "Choose a name for this account"}, &account.Name, nil)
	return nil
}

func generateAccount(account *types.Account) error {
	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Generating secure key..."})
	time.Sleep(time.Second)
	s.Stop()

	// TODO add real account creation
	account.Address = "0x0000000000000000000000000000000000000000"
	account.Seed = "this is my long secure seed that help me regenerate my account keys"
	return nil
}

func displayResume(account *types.Account) {
	success := color.New(color.FgGreen, color.Bold).SprintFunc()
	warning := color.New(color.FgYellow, color.Bold).SprintFunc()
	fmt.Println("Here is all the details of your account:")
	fmt.Println()
	fmt.Printf("Account name: %s\n", success(account.Name))
	fmt.Printf("Account address: %s\n", success(account.Address))
	fmt.Printf("Seed: %s\n", success(account.Seed))
	fmt.Println()
	fmt.Printf("%s", warning("Please make sure that you save those informations and especially the following seed that might be needed to regenerate this address"))
	fmt.Println()
}

func init() {
	Create.Flags().StringP("name", "n", "", "Name of the account")
	Create.Flags().StringP("password", "p", "", "Password for the account")
}
