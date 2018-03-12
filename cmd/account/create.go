package cmdAccount

import (
	"fmt"
	"time"

	"github.com/mesg-foundation/application/cmd/utils"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/spf13/cobra"
)

// Create run the create command for an account
var Create = &cobra.Command{
	Use:               "create",
	Short:             "Create a new account",
	Example:           "mesg-cli account create",
	Run:               createHandler,
	DisableAutoGenTag: true,
}

func createHandler(cmd *cobra.Command, args []string) {
	var password, passwordConfirmation, privateSeed string
	survey.AskOne(&survey.Password{Message: "Please set a password ?"}, &password, nil)
	survey.AskOne(&survey.Password{Message: "Repeat your password ?"}, &passwordConfirmation, nil)
	if password != passwordConfirmation {
		fmt.Println("Password confirmation invalid")
		return
	}
	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Generating secure key..."})
	time.Sleep(time.Second)
	s.Stop()
	survey.AskOne(&survey.Input{Message: "Repeat your private seed"}, &privateSeed, nil)
	// TODO add real account creation
	s = cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "processing..."})
	time.Sleep(2 * time.Second)
	s.Stop()
}
