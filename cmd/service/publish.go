package cmdService

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/application/cmd/utils"
	"github.com/mesg-foundation/application/service"

	"github.com/spf13/cobra"
)

// Publish a service to the marketplace
var Publish = &cobra.Command{
	Use:               "publish SERVICE_FILE",
	Short:             "Publish a new service",
	Long:              "Deploy a service to the Network from a given service file. Validate it first. The user will need to provide an account and the password of the account.",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli service publish service.yml",
	Run:               deployHandler,
	DisableAutoGenTag: true,
}

func deployHandler(cmd *cobra.Command, args []string) {
	valid, warnings, err := service.ValidService(args[0])
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	if !valid {
		for _, warning := range warnings {
			fmt.Println(aurora.Brown(warning))
		}
		return
	}
	serviceFile := filepath.Join(args[0], "mesg.yml")
	service, err := service.ImportFromFile(serviceFile)
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
