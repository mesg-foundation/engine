package cmdWorkflow

import (
	"fmt"
	"time"

	"github.com/mesg-foundation/core/cmd/utils"

	"github.com/spf13/cobra"
)

// Deploy run the deploy command for a workflow
var Deploy = &cobra.Command{
	Use:   "deploy ./PATH_TO_WORKFLOW_FILE",
	Short: "Deploy a workflow",
	Long: `Deploy a workflow on the Network.

To get more information, see the [deploy page from the documentation](https://docs.mesg.tech/workflow/deploy.html)`,
	Args: cobra.MinimumNArgs(1),
	Example: `mesg-cli workflow deploy ./PATH_TO_WORKFLOW_FILE.yml
mesg-cli workflow deploy ./PATH_TO_WORKFLOW_FILE.yml --account ACCOUNT --amount AMOUNT --confirm`,
	Run:               deployHandler,
	DisableAutoGenTag: true,
}

func deployHandler(cmd *cobra.Command, args []string) {
	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Select an account:")
	amount, err := cmdUtils.GetOrAskAmount(cmd, "How much do you want to deposit in your workflow?")
	if err != nil {
		panic(err)
	}
	if !cmdUtils.Confirm(cmd, "Are you sure?") {
		return
	}
	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Deployment in progress..."})
	time.Sleep(2 * time.Second)
	s.Stop()
	// TODO deploy the workflow
	fmt.Println("Workflow deployed with success", args, account, amount)
}

func init() {
	cmdUtils.Confirmable(Deploy)
	cmdUtils.Accountable(Deploy)
	cmdUtils.Payable(Deploy)
}
