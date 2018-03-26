package cmdService

import (
	"fmt"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// Publish a service to the marketplace
var Publish = &cobra.Command{
	Use:               "publish SERVICE_PATH",
	Short:             "Publish a new service",
	Long:              "Deploy a service to the Network from a given service. Validate it first. The user will need to provide an account and the password of the account.",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli service publish /path/to/the/service/folder",
	Run:               deployHandler,
	DisableAutoGenTag: true,
}

func deployHandler(cmd *cobra.Command, args []string) {
	if !validateServicePath(args[0]) {
		return
	}
	service, err := importService(args[0])
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Select an account")
	if !cmdUtils.Confirm(cmd, "Are you sure ?") {
		return
	}
	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Deployment of " + service.Name + " in progress..."})
	time.Sleep(2 * time.Second)
	s.Stop()
	// TODO deploy the service
	fmt.Println("service deployed with", account)
}

func init() {
	cmdUtils.Confirmable(Publish)
	cmdUtils.Accountable(Publish)
}
