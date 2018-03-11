package cmdService

import (
	"fmt"
	"time"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// Deploy run the deploy command for a service
var Deploy = &cobra.Command{
	Use:     "deploy",
	Short:   "Deploy a new service",
	Args:    cobra.MinimumNArgs(1),
	Example: "mesg-cli service deploy service.yml",
	Run:     deployHandler,
}

func deployHandler(cmd *cobra.Command, args []string) {
	if !cmdUtils.Confirm(cmd, "Are you sure ?") {
		return
	}
	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Deployment in progress..."})
	time.Sleep(2 * time.Second)
	s.Stop()
	// TODO deploy the service
	fmt.Println("service deployed", args)
}

func init() {
	Deploy.Flags().BoolP("confirm", "c", false, "Confirm")
}
