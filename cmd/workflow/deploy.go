package cmdWorkflow

import (
	"fmt"
	"time"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// Deploy run the deploy command for a workflow
var Deploy = &cobra.Command{
	Use:               "deploy FILE",
	Short:             "Deploy a new workflow",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli workflow deploy workflow.yml",
	Run:               deployHandler,
	DisableAutoGenTag: true,
}

func deployHandler(cmd *cobra.Command, args []string) {
	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Select an account for deployment")
	amount, err := cmdUtils.GetOrAskAmount(cmd, "How much do you want to deposit in your workflow ?")
	if err != nil {
		panic(err)
	}
	if !cmdUtils.Confirm(cmd, "Are you sure ?") {
		return
	}
	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Deployment in progress..."})
	time.Sleep(2 * time.Second)
	s.Stop()
	// TODO deploy the workflow
	fmt.Println("workflow deployed", args, account, amount)
}

func init() {
	cmdUtils.Confirmable(Deploy)
	cmdUtils.Accountable(Deploy)
	cmdUtils.Payable(Deploy)
}
