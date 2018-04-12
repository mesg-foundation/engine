package cmdService

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/application/service"

	"github.com/spf13/cobra"
)

// Validate a service
var Validate = &cobra.Command{
	Use:   "validate",
	Short: "Validate a service file",
	Long: `Validate a service file. Check the yml format and rules.

All the definitions of the service file can be found in the page [Service File from the documentation](https://docs.mesg.tech/service/service-file.html).`,
	Example: `mesg-cli service validate
mesg-cli service validate ./SERVICE_FOLDER`,
	Run:               validateHandler,
	DisableAutoGenTag: true,
}

func validateHandler(cmd *cobra.Command, args []string) {
	path := "./"
	if len(args) > 0 {
		path = args[0]
	}

	warnings, err := service.ValidService(path)
	if err != nil {
		fmt.Println(aurora.Red("Service error").Bold())
		fmt.Println(err)

		for _, warning := range warnings {
			fmt.Println(aurora.Red("The service file contains errors:").Bold())
			fmt.Println(warning)
		}
	}
	if err == nil {
		fmt.Println(aurora.Green("Service is valid"))
	}
}
