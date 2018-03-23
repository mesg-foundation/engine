package cmdService

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/application/service"

	"github.com/spf13/cobra"
)

// Validate a service
var Validate = &cobra.Command{
	Use:               "validate SERVICE_FILE",
	Short:             "Validate a service file. Check the yml format and rules.",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli service validate service.yml",
	Run:               validateHandler,
	DisableAutoGenTag: true,
}

func validateHandler(cmd *cobra.Command, args []string) {
	res, err := service.ValidServiceFile(args[0])
	if err != nil {
		panic(err)
	}
	if res.Valid() {
		fmt.Println(aurora.Green("Service is valid"))
	} else {
		fmt.Println(aurora.Red("Service contains errors"))
		for _, err := range res.Errors() {
			fmt.Println(err)
		}
	}
}
