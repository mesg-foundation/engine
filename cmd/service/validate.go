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

All the definitions of the service file can be found in the page [Service File from the documentation](https://docs.mesg.tech/service/develop/service-file.html).`,
	Args: cobra.MinimumNArgs(1),
	Example: `mesg-cli service validate
mesg-cli service validate ./SERVICE_FOLDER`,
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
