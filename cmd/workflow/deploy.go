package workflow

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// Deploy run the deploy command for a workflow
var Deploy = &cobra.Command{
	Use:   "deploy ./PATH_TO_WORKFLOW_FILE",
	Short: "Deploy a workflow",
	Long: `Deploy a workflow on the Network.

To get more information, see the [deploy page from the documentation](https://docs.mesg.tech/workflow/deploy.html)`,
	Args: cobra.MinimumNArgs(1),
	Example: `mesg-core workflow deploy ./PATH_TO_WORKFLOW_FILE.yml
mesg-core workflow deploy ./PATH_TO_WORKFLOW_FILE.yml --account ACCOUNT --amount AMOUNT --confirm`,
	Run:               deployHandler,
	DisableAutoGenTag: true,
}

func deployHandler(cmd *cobra.Command, args []string) {
	account := utils.AccountFromFlagOrAsk(cmd, "Select an account:")
	amount, err := utils.GetOrAskAmount(cmd, "How much do you want to deposit in your workflow?")
	if err != nil {
		panic(err)
	}
	if !utils.Confirm(cmd, "Are you sure?") {
		return
	}
	s := utils.StartSpinner(utils.SpinnerOptions{Text: "Deployment in progress..."})
	time.Sleep(2 * time.Second)
	s.Stop()
	// TODO deploy the workflow
	fmt.Println("Workflow deployed with success", args, account, amount)
}

func init() {
	utils.Confirmable(Deploy)
	utils.Accountable(Deploy)
	utils.Payable(Deploy)
}
