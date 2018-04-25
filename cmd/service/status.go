package cmdService

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/service"

	"github.com/spf13/cobra"
)

// Status command returns started services
var Status = &cobra.Command{
	Use:               "status",
	Short:             "List started services",
	Example:           "mesg-cli service status",
	Run:               statusHandler,
	DisableAutoGenTag: true,
}

func statusHandler(cmd *cobra.Command, args []string) {
	services, err := service.List()
	if err != nil {
		fmt.Println(aurora.Red(err))
		os.Exit(0)
	}
	fmt.Println("Services running :")
	for _, service := range services {
		fmt.Println(aurora.Bold(" - " + service))
	}
	if len(services) == 0 {
		fmt.Println("No services running")
	}
}
