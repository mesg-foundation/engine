package workflow

import (
	"fmt"
	"os"

	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// Topup a workflow
var Topup = &cobra.Command{
	Use:   "topup WORKFLOW_ID",
	Short: "Top-up a workflow",
	Long: `Top-up a workflow.
Add more token to a existing workflow.`,
	Example: `mesg-core workflow topup WORKFLOW_ID
mesg-core workflow topup WORKFLOW_ID --amount AMOUNT --account ACCOUNT_ID --confirm`,
	Run:               topupHandler,
	DisableAutoGenTag: true,
}

func topupHandler(cmd *cobra.Command, args []string) {
	account := utils.AccountFromFlagOrAsk(cmd, "Select an account:")
	amount, err := utils.GetOrAskAmount(cmd, "How much do you want to deposit in your workflow?")
	if err != nil {
		panic(err)
	}
	var workflow = ""
	if len(args) > 0 {
		workflow = args[0]
	}
	if workflow == "" {
		// TODO add real list
		workflows := []string{"Workflow #1", "Workflow #2"}
		if survey.AskOne(&survey.Select{
			Message: "Choose the workflow you want to pause",
			Default: workflows[0],
			Options: workflows,
		}, &workflow, nil) != nil {
			os.Exit(0)
		}
	}
	if !utils.Confirm(cmd, "Are you sure ?") {
		return
	}
	// TODO topup the workflow onchain
	fmt.Println("workflow topup", args, account, amount)
}

func init() {
	utils.Confirmable(Topup)
	utils.Accountable(Topup)
	utils.Payable(Topup)
}
