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
	service, err := service.ImportFromFile(args[0])
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	valid, errs := service.IsValid()
	if valid == true {
		fmt.Println(aurora.Green("Service " + service.Name + " is valid"))
	} else {
		fmt.Println(aurora.Red("Service " + service.Name + " contains errors"))
		if errs != nil {
			fmt.Println(errs)
		}
	}
}
