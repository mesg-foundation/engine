package service

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
	"github.com/spf13/cobra"
)

// Status command returns started services
var Status = &cobra.Command{
	Use:               "status",
	Short:             "List started services",
	Example:           "mesg-core service status",
	Run:               statusHandler,
	DisableAutoGenTag: true,
}

func statusHandler(cmd *cobra.Command, args []string) {
	hashes, err := service.ListRunning() // TODO: should use the API
	utils.HandleError(err)
	fmt.Println("Running services:")
	for _, hash := range hashes {
		service, err := services.Get(hash)
		utils.HandleError(err)
		fmt.Println(aurora.Bold(" - " + hash + " - " + service.Name))
	}
	if len(hashes) == 0 {
		fmt.Println("No service are running")
	}
}
