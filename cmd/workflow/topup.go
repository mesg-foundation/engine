package cmdWorkflow

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/spf13/cobra"
)

// Topup a workflow
var Topup = &cobra.Command{
	Use:               "topup ID",
	Short:             "Top-up a workflow",
	Example:           "mesg-cli workflow topup xxx --amount XX",
	Run:               topupHandler,
	DisableAutoGenTag: true,
}

func topupHandler(cmd *cobra.Command, args []string) {
	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Select an account")
	amount := cmdUtils.GetOrAskAmount(cmd, "How much do you want to deposit in your workflow ?")
	var workflow = ""
	if len(args) > 0 {
		workflow = args[0]
	}
	if workflow == "" {
		// TODO add real list
		workflows := []string{"Workflow #1", "Workflow #2"}
		survey.AskOne(&survey.Select{
			Message: "Choose the workflow you want to pause",
			Default: workflows[0],
			Options: workflows,
		}, &workflow, nil)
	}
	if !cmdUtils.Confirm(cmd, "Are you sure ?") {
		return
	}
	// TODO topup the workflow onchain
	fmt.Println("workflow topup", args, account, amount)
}

func init() {
	cmdUtils.Confirmable(Topup)
	cmdUtils.Accountable(Topup)
	cmdUtils.Payable(Topup)
}
