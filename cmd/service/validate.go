package cmdService

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/application/service"

	"github.com/spf13/cobra"
)

// Validate a service
var Validate = &cobra.Command{
	Use:               "validate SERVICE_PATH",
	Short:             "Validate a service. Check the mesg.yml file for format and rules and do some additional tests about the directory",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli service validate service.yml",
	Run:               validateHandler,
	DisableAutoGenTag: true,
}

func validateHandler(cmd *cobra.Command, args []string) {

	valid, warnings, err := service.ValidService(args[0])
	if err != nil {
		panic(err)
	}
	if valid {
		fmt.Println(aurora.Green("Service is valid"))
	} else {
		fmt.Println(aurora.Red("Service contains errors"))
		for _, err := range warnings {
			fmt.Println(err)
		}
	}
}
